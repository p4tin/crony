package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/robfig/cron"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"cronie/messages"
)

/*
TODO:
	- Add a job
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

type AJob struct {
	JobID     string    `json:"job_id"`
	JobName   string    `json:"job_name"`
	Cmd       string    `json:"cmd"`
	Args      []string  `json:"args"`
	Schedule  string    `json:"schedule"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

var AJobs map[string]AJob
var crn *cron.Cron

func init() {
	crn = cron.New()
}

type TaskRunner struct {
	JobId string
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
	log.Infof("Task Status: %s", r.Ok)
	return true
}

func AddJob() {
	jid := xid.New().String()
	AJobs[jid] = AJob{
		JobID:     jid,
		Cmd:       "sh",
		Args:      []string{"test.sh"},
		Schedule:  "00 * * * * *",
		Status:    "New",
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
	err := crn.AddFunc(AJobs[jid].Schedule, func() { TaskRunner{JobId: jid}.Run() })
	if err != nil {
		log.Errorln(err.Error())
		os.Exit(0)
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
	AJobs = make(map[string]AJob)
	AddJob()
	go runRESTServer()
	go runGRPCServer()
	log.Infoln("Cronie Server Started on ports %s/%s\n", grpc_server_port, rest_server_port)

	<-sigChan

}

type client struct{}

func (c *client) TaskStatus(ctx context.Context, task *messages.TaskStatusMsg) (*messages.Ack, error) {
	log.Printf("Received: %d\n", task.Status)
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
	jid := ctx.Param("jid")
	return ctx.JSON(http.StatusOK, AJobs[jid])
}

func AddAJob(ctx echo.Context) error {
	job := AJob{}
	if err := ctx.Bind(&job); err != nil {
		return err
	}
	job.JobID = xid.New().String()
	job.Status = "New"
	job.UpdatedAt = time.Now()
	job.CreatedAt = time.Now()

	AJobs[job.JobID] = job
	return ctx.JSON(http.StatusOK, job)
}

func runRESTServer() {
	c := echo.New()
	c.HideBanner = true
	c.HidePort = true
	c.GET("/jobs", GetJobs)     // List all jobs
	c.GET("/job/:jid", GetAJob) // List all jobs
	c.POST("/jobs", AddAJob)

	c.Start(rest_server_port)
}
