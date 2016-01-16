package main

import (
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/yangmls/vcron"
)

type Server struct {
	Port  string
	Conns []net.Conn
}

func (server *Server) Run() {
	Listener, err := net.Listen("tcp", ":"+server.Port)

	if err != nil {
		return
	}

	defer Listener.Close()

	fmt.Println("Server is running on " + server.Port)

	for {
		conn, err := Listener.Accept()

		if err != nil {
			break
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	message := read(conn)

	if message == nil {
		return
	}

	if message.GetType() != "register" {
		return
	}

	AddConn(message.GetName(), conn)

	for {
		message := read(conn)

		if message == nil {
			conn.Close()
			return
		}

		fmt.Println(message)
	}
}

func read(conn net.Conn) (m *vcron.Message) {
	data := make([]byte, 4096)
	len, readErr := conn.Read(data)

	if readErr != nil {
		return
	}

	message := &vcron.Message{}
	uncodeErr := proto.Unmarshal(data[0:len], message)

	if uncodeErr != nil {
		fmt.Println(uncodeErr)
		return
	}

	return message
}
