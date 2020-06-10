package main

import (
	"context"
	"net"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"cronie/messages"
)

const (
	server_address = "localhost:50050"
	client_port    = ":50051"
)

func SendTaskStatus(status int32) {
	conn, err := grpc.Dial(server_address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := messages.NewCronyTaskSlaveClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.TaskStatus(ctx, &messages.TaskStatusMsg{Status: status})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Task Status: %s", r.Ok)
}

func runTask(task *messages.TaskDef) {
	cmd := exec.Command(task.Executable, task.Args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	log.Infof("combined out:\n%s\n", string(out))
	SendTaskStatus(2)
}

type server struct {
}

func (s *server) RunTask(ctx context.Context, task *messages.TaskDef) (*messages.Ack, error) {
	log.Printf("Received: %v", task.Executable)
	go runTask(task)
	return &messages.Ack{Ok: true}, nil
}

func (s *server) Ping(ctx context.Context, emp *messages.Empty) (*messages.Ack, error) {
	return &messages.Ack{}, nil
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	lis, err := net.Listen("tcp", client_port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	messages.RegisterCronyTaskMasterServer(s, &server{})
	log.Infoln("Cronie Slave Starting on port %s/%s\n", client_port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
