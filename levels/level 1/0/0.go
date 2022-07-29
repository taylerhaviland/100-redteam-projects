package main

import (
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		os.Exit(1)
	}

	defer l.Close()
	fmt.Println("Listening on: " + CONN_HOST + ":" + CONN_PORT)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: " + err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}

}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection: " + err.Error())
	}
	fmt.Println("Message: " + string(buf))
	conn.Write([]byte("Message received.\n"))

	conn.Close()
}
