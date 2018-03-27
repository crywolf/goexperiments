package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"log"

	"github.com/crywolf/goexperiments/grpc/todo"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing subcommand: list or add")
		os.Exit(1)
	}

	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to backend: %v", err)
	}
	client := todo.NewTasksClient(conn)

	switch cmd := flag.Arg(0); cmd {
	case "list":
		err = list(context.Background(), client)
	case "add":
		err = add(context.Background(), client, strings.Join(flag.Args()[1:], " "))
	default:
		err = fmt.Errorf("unknown subcommand %s", cmd)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func add(ctx context.Context, client todo.TasksClient, text string) error {
	_, err := client.Add(ctx, &todo.Text{text})
	if err != nil {
		return fmt.Errorf("could not add task: %v", err)
	}
	return nil
}

func list(ctx context.Context, client todo.TasksClient) error {
	l, err := client.List(ctx, &todo.Void{})
	if err != nil {
		return fmt.Errorf("could not fetch task list: %v", err)
	}

	for _, task := range l.Tasks {
		if task.Done {
			fmt.Printf("[x]")
		} else {
			fmt.Printf("[-]")
		}
		fmt.Printf(" %s\n", task.Text)
	}
	return nil
}
