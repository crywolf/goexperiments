package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"os"

	"github.com/crywolf/goexperiments/grpc/todo"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	port := flag.String("port", "8888", "defines port")
	flag.Parse()

	var taskSrv taskServer
	srv := grpc.NewServer()
	todo.RegisterTasksServer(srv, taskSrv)

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		log.Fatalf("could not start listening on port %s: %v", *port, err)
	}
	log.Fatal(srv.Serve(l))
}

const (
	sizeOfLength = 8
	dbPath       = "../../myDb.pb"
)

func (s taskServer) Add(ctx context.Context, text *todo.Text) (*todo.Task, error) {
	task := &todo.Task{
		Text: text.Text,
		Done: false,
	}

	b, err := proto.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("coud not encode task %v", err)
	}

	buf := proto.NewBuffer(nil)
	if err = buf.EncodeFixed64(uint64(len(b))); err != nil {
		return nil, fmt.Errorf("could not encode message length:, %v", err)
	}

	messageLength := buf.Bytes()
	b = append(messageLength, b...)

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %v", dbPath, err)
	}

	_, err = f.Write(b)
	if err != nil {
		return nil, fmt.Errorf("could not write task to file %v", err)
	}

	if err = f.Close(); err != nil {
		return nil, fmt.Errorf("could not close file %s: %v", dbPath, err)
	}

	return task, nil
}

type taskServer struct {
}

func (s taskServer) List(ctx context.Context, void *todo.Void) (*todo.TaskList, error) {
	b, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %v", dbPath, err)
	}

	var taskList todo.TaskList

	for {
		if len(b) == 0 {
			return &taskList, nil
		}
		if len(b) < sizeOfLength {
			return nil, fmt.Errorf("remaining wrong number of bytes: %d", len(b))
		}

		var length uint64
		buf := proto.NewBuffer(b)
		length, err := buf.DecodeFixed64()
		if err != nil {
			return nil, fmt.Errorf("could not decode message length:, %v", err)
		}
		b = b[sizeOfLength:]

		var task todo.Task
		if err = proto.Unmarshal(b[:length], &task); err != nil {
			return nil, fmt.Errorf("could not read task: %v", err)
		}
		b = b[length:]

		taskList.Tasks = append(taskList.Tasks, &task)
	}
}
