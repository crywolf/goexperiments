package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var data map[string]string

func main() {
	data = make(map[string]string)

	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	li, err := net.Listen("tcp", ":"+port)
	log.Println("Listening on localhost:" + port)
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

	io.WriteString(conn, "\r\nIN-MEMORY DATABASE\r\n\r\n"+
		"USE:\r\n"+
		"\tSET key value \r\n"+
		"\tGET key \r\n"+
		"\tDEL key \r\n\r\n"+
		"EXAMPLE:\r\n"+
		"\tSET fav chocolate \r\n"+
		"\tGET fav \r\n\r\n\r\n")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		log.Println(ln)
		resp, err := processCommand(ln)
		if err != nil {
			log.Println(err)
		}
		io.WriteString(conn, resp)
	}
}

func processCommand(ln string) (string, error) {
	ret := ""

	fs := strings.Fields(ln)
	if len(fs) == 0 {
		return "Missing command!\n", nil
	}

	cmd := fs[0]
	args := fs[1:]

	switch cmd {
	case "GET":
		ok, message := checkCommand(cmd, args)
		if !ok {
			return message, nil
		}
		k := fs[1]
		ret = data[k] + "\n"
	case "SET":
		ok, message := checkCommand(cmd, args)
		if !ok {
			return message, nil
		}
		k := fs[1]
		v := fs[2]
		data[k] = v
	case "DEL":
		ok, message := checkCommand(cmd, args)
		if !ok {
			return message, nil
		}
		k := fs[1]
		delete(data, k)
	default:
		ret = "Invalid command!\n"
	}

	return ret, nil
}

func checkCommand(cmd string, args []string) (ok bool, message string) {
	switch cmd {
	case "GET", "DEL":
		if len(args) == 0 {
			return false, "Missing key!\n"
		}
	case "SET":
		if len(args) == 0 {
			return false, "Missing key and value!\n"
		}
		if len(args) == 1 {
			return false, "Missing value!\n"
		}
	}
	return true, ""
}
