package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/crywolf/goexperiments/protobuf/todo"
	"github.com/golang/protobuf/proto"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing subcommand: list or add")
		os.Exit(1)
	}

	var err error
	switch cmd := flag.Arg(0); cmd {
	case "list":
		err = list()
	case "add":
		err = add(strings.Join(flag.Args()[1:], " "))
	default:
		err = fmt.Errorf("unknown subcommand %s", cmd)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const (
	sizeOfLength = 8
	dbPath       = "myDb.pb"
)

func add(text string) error {
	task := &todo.Task{
		Text: text,
		Done: false,
	}

	b, err := proto.Marshal(task)
	if err != nil {
		return fmt.Errorf("coud not encode task %v", err)
	}

	buf := proto.NewBuffer(nil)
	if err = buf.EncodeFixed64(uint64(len(b))); err != nil {
		return fmt.Errorf("could not encode message length:, %v", err)
	}

	messageLength := buf.Bytes()
	b = append(messageLength, b...)

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("could not open file %s: %v", dbPath, err)
	}

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("could not write task to file %v", err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("could not close file %s: %v", dbPath, err)
	}

	return nil
}

func list() error {
	b, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return fmt.Errorf("could not read file %s: %v", dbPath, err)
	}

	for {
		if len(b) == 0 {
			return nil
		}
		if len(b) < sizeOfLength {
			return fmt.Errorf("remaining wrong number of bytes: %d", len(b))
		}

		var length uint64
		buf := proto.NewBuffer(b)
		length, err := buf.DecodeFixed64()
		if err != nil {
			return fmt.Errorf("could not decode message length:, %v", err)
		}
		b = b[sizeOfLength:]

		var task todo.Task
		if err = proto.Unmarshal(b[:length], &task); err != nil {
			return fmt.Errorf("could not read task: %v", err)
		}
		b = b[length:]

		if task.Done {
			fmt.Printf("[x]")
		} else {
			fmt.Printf("[-]")
		}
		fmt.Printf(" %s\n", task.Text)
	}
}
