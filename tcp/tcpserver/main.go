package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panic(err)
	}

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
		}
		log.Println("Connection accepted")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	conn.SetDeadline(time.Now().Add(time.Second * 5))
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		io.WriteString(conn, "Hello from server. Thanks for calling me!\n")
	}

	log.Println("Closing connection")
}
