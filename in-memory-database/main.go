package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/crywolf/goexperiments/in-memory-database/storage"
)

var (
	database storage.Db
	port     = flag.String("port", "8080", "Listen on port number")
)

func main() {
	flag.Parse()
	database = storage.GetDatabase()

	li, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	log.Printf("Listening on localhost:%s\n", *port)
	if err != nil {
		log.Panic(err)
	}

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	provideIntroductionHelp(conn)

	handleCommands(conn)
}

func provideIntroductionHelp(conn net.Conn) (int, error) {
	return io.WriteString(conn, "\r\nIN-MEMORY DATABASE\r\n\r\n"+
		"USE:\r\n"+
		"\tSET key value \r\n"+
		"\tGET key \r\n"+
		"\tDEL key \r\n"+
		"\tQUIT\r\n\r\n"+
		"EXAMPLE:\r\n"+
		"\tSET fav chocolate \r\n"+
		"\tGET fav \r\n\r\n\r\n")
}

func handleCommands(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()

		fs := strings.Fields(ln)
		if len(fs) == 0 {
			io.WriteString(conn, "-> Missing command!\n")
			continue
		}

		cmd := fs[0]
		args := fs[1:]

		if cmd == "QUIT" {
			break
		}

		resp, err := processCommand(cmd, args)
		if err != nil {
			log.Println(ln, ":", err)
			continue
		}

		log.Println(ln)

		if resp != "" {
			resp = fmt.Sprintf("-> %s\n", resp)
		}

		io.WriteString(conn, resp)
	}
}

func processCommand(cmd string, args []string) (string, error) {
	switch cmd {
	case "GET":
		ok, message := checkCommand(cmd, args)
		if !ok {
			return message, nil
		}
		key := args[0]
		return database.Get(key)
	case "SET":
		ok, message := checkCommand(cmd, args)
		if !ok {
			return message, nil
		}
		key := args[0]
		val := args[1]
		return database.Set(key, val)
	case "DEL":
		ok, message := checkCommand(cmd, args)
		if !ok {
			return message, nil
		}
		key := args[0]
		return database.Del(key)
	default:
		return "Invalid command!", nil
	}
}

func checkCommand(cmd string, args []string) (ok bool, message string) {
	switch cmd {
	case "GET", "DEL":
		if len(args) == 0 {
			return false, "Missing key!"
		}
	case "SET":
		if len(args) == 0 {
			return false, "Missing key and value!"
		}
		if len(args) == 1 {
			return false, "Missing value!"
		}
	}

	return true, ""
}
