package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "Hello from client!\n")
	output, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(output)
}
