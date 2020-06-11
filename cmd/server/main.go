package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"cronie/messages"
)

/*
TODO:
	- Pause a job
	- Update a job
	- Delete a job
	- Save jobs to leveldb?
	- More than 1 slave (environments/slave selection)
	- Job composition
		-- on success
			--- Chaining tasks on success
		-- on failure
			--- Alert on failure is just another job to run
	- A rudimentary UI (native/web ???)
*/
const (
	grpc_server_port = ":50050"
	rest_server_port = ":8080"
	client_address   = "localhost:50051"
)

var StatusNames = map[int32]string{
	0: "None",
	1: "Started",
	2: "Completed",
	4: "Error",
}

type AJob struct {
	JobID     int32    `json:"job_id"`
	JobName   string    `json:"job_name"`
	Cmd       string    `json:"cmd"`
	Args      []string  `json:"args"`
	Schedule  string    `json:"schedule"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

var AJobs map[int32]AJob
var crn *cron.Cron

func init() {
	crn = cron.New(cron.WithSeconds())
}

type TaskRunner struct {
	JobId int32
}

func (tr TaskRunner) Run() bool {
	conn, err := grpc.Dial(client_address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := messages.NewCronyTaskMasterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.RunTask(ctx, &messages.TaskDef{
		RequestID:  tr.JobId,
		Executable: AJobs[tr.JobId].Cmd,
		Args:       AJobs[tr.JobId].Args,
		Env:        nil,
		Dir:        "",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Infof("Task received by slave: %t", r.Ok)
	return true
}

func AddJob() {
	tr := TaskRunner{}
	eid, err := crn.AddFunc("00 * * * * *", func() { tr.Run() })
	if err != nil {
		log.Errorln(err.Error())
		os.Exit(0)
	}
	tr.JobId = int32(eid)

	AJobs[int32(eid)] = AJob{
		JobID:     int32(eid),
		JobName: "Test 0",
		Cmd:       "sh",
		Args:      []string{"test.sh"},
		Schedule:  "00 * * * * *",
		Status:    "New",
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	crn.Start()
}

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		crn.Stop()
		log.Infoln("Clean exit")
	}()
	log.SetFormatter(&log.JSONFormatter{})
	AJobs = make(map[int32]AJob)
	AddJob()
	go runRESTServer()
	go runGRPCServer()
	log.Infoln("Cronie Server Started on ports %s/%s\n", grpc_server_port, rest_server_port)

	<-sigChan

}

type client struct{}

func (c *client) TaskStatus(ctx context.Context, task *messages.TaskStatusMsg) (*messages.Ack, error) {
	log.Printf("Received: %s/%s" +
		"\n", AJobs[task.RequestID].JobName, StatusNames[task.Status])
	return &messages.Ack{Ok: true}, nil
}

func (c *client) Register(ctx context.Context, task *messages.Registration) (*messages.Empty, error) {
	return &messages.Empty{}, nil
}

func runGRPCServer() {
	lis, err := net.Listen("tcp", grpc_server_port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	messages.RegisterCronyTaskSlaveServer(s, &client{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func GetJobs(ctx echo.Context) error {
	jobs := make([]AJob, 0)
	for _, job := range AJobs {
		jobs = append(jobs, job)
	}
	return ctx.JSON(http.StatusOK, jobs)
}

func GetAJob(ctx echo.Context) error {
	jidStr := ctx.Param("jid")
	jid, err := strconv.Atoi(jidStr)
	if err != nil {
		return err
	}
	if job, ok := AJobs[int32(jid)]; ok {
		return ctx.JSON(http.StatusOK, job)
	}
	return ctx.JSON(http.StatusNotFound, map[string]string{"status": "not found"})
}

func AddAJob(ctx echo.Context) error {
	job := AJob{}
	if err := ctx.Bind(&job); err != nil {
		return err
	}

	tr := TaskRunner{}
	eid, err := crn.AddFunc(job.Schedule, func() { tr.Run() })
	if err != nil {
		log.Errorln(err.Error())
		os.Exit(0)
	}
	tr.JobId = int32(eid)

	job.JobID = int32(eid)
	job.Status = "New"
	job.UpdatedAt = time.Now()
	job.CreatedAt = time.Now()

	AJobs[job.JobID] = job

	return ctx.JSON(http.StatusOK, job)
}

func DeleteAJob(ctx echo.Context) error {
	jidStr := ctx.Param("jid")
	jid, err := strconv.Atoi(jidStr)
	if err != nil {
		return err
	}
	if _, ok := AJobs[int32(jid)]; ok {
		crn.Remove(cron.EntryID(jid))
		delete(AJobs, int32(jid))

		return ctx.JSON(http.StatusOK, map[string]string{"deleted": fmt.Sprintf("job %d=ok", jid)})
	}
	return ctx.JSON(http.StatusNotFound, map[string]string{"status": "not found"})
}

func runRESTServer() {
	c := echo.New()
	c.HideBanner = true
	c.HidePort = true
	c.GET("/jobs", GetJobs)     	// List all jobs
	c.GET("/jobs/:jid", GetAJob) 	// List all jobs
	c.POST("/jobs", AddAJob) 		// Add a job to the job list
	c.DELETE("/jobs/:jid", DeleteAJob)

	c.Start(rest_server_port)
}
