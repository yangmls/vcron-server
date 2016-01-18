package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yangmls/vcron"
	"net"
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
	var (
		message *vcron.Message
		err     error
	)

	defer conn.Close()

	message, err = read(conn)

	if message == nil {
		return
	}

	if message.GetType() != "register" {
		return
	}

	id := AddConn(message.GetName(), conn)

	for {
		message, err = read(conn)

		if err != nil {
			RemoveCon(id)
			break
		}
	}
}

func read(conn net.Conn) (*vcron.Message, error) {
	data := make([]byte, 4096)
	len, readErr := conn.Read(data)

	if readErr != nil {
		return nil, readErr
	}

	message := &vcron.Message{}
	uncodeErr := proto.Unmarshal(data[0:len], message)

	if uncodeErr != nil {
		fmt.Println(uncodeErr)
		return nil, nil
	}

	return message, nil
}
