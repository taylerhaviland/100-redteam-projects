package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	SERVER_CONN_HOST = "localhost"
	SERVER_CONN_PORT = "3333"
	SERVER_CONN_TYPE = "tcp"
)

type userInfo struct {
	username   string
	connection net.Conn
	channel    chan string
}

func main() {
	// c1 := make(chan string)
	var userSlice []userInfo

	l, err := net.Listen(SERVER_CONN_TYPE, SERVER_CONN_HOST+":"+SERVER_CONN_PORT)
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		os.Exit(1)
	}

	defer l.Close()
	fmt.Println("Listening on: " + SERVER_CONN_HOST + ":" + SERVER_CONN_PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: " + err.Error())
			os.Exit(1)
		}
		fmt.Println("Remote Network Connected: " + conn.RemoteAddr().String())

		newUser := addUser(conn)
		fmt.Printf("New user added: %s\n", newUser.username)

		userSlice = append(userSlice, newUser)
		go newUserConnection(newUser)
		go func() {
			for i := 0; i < len(userSlice); i++ {
				username := userSlice[i].username
				msg := userSlice[i].channel
				select {
				case <-msg:
					for j := 0; j < len(userSlice); j++ {
						if username != userSlice[j].username {
							userSlice[j].connection.Write([]byte(username + ": " + <-msg + "\n"))
						}
					}
				default:
					// pass
				}
			}
		}()
	}
}

func welcomeMessage(conn net.Conn) {
	conn.Write([]byte("Welcome to Moose's chat server!\n"))
}

func addUser(conn net.Conn) userInfo {
	welcomeMessage(conn)
	buf := make([]byte, 1024)
	conn.Write([]byte("Enter a username: "))
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection: " + err.Error())
		conn.Close()
	}
	username := bytes.NewBuffer(buf).String()
	username = strings.TrimRight(username, "\x00")
	username = strings.TrimSpace(username)
	newUser := userInfo{
		username:   username,
		connection: conn,
		channel:    make(chan string),
	}
	newUser.connection.Write([]byte("Welcome " + username + "\n"))

	return newUser
}

func newUserConnection(user userInfo) {
	for {
		user.connection.Write([]byte(user.username + ": "))
		buf := make([]byte, 1024)
		_, err := user.connection.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println(user.username + " has left the server.")
				user.connection.Close()
				return
			}
		}
		msg := bytes.NewBuffer(buf).String()
		msg = strings.TrimRight(msg, "\x00")
		msg = strings.TrimSpace(msg)
		fmt.Println(user.username + ": " + string(msg))
		user.channel <- msg
	}
}
