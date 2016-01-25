package main

import (
	"fmt"
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

		go AddConn(conn)
	}
}
