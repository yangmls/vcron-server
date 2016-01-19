package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yangmls/vcron"
	"net"
	"sync"
)

var (
	ConnId = 0
	Conns  = make(map[int]*Conn)
)

type Conn struct {
	Id    int
	Name  string
	I     net.Conn
	Mutex *sync.Mutex
}

func AddConn(name string, conn net.Conn) int {
	ConnId = ConnId + 1
	c := &Conn{
		Id:    ConnId,
		Name:  name,
		I:     conn,
		Mutex: new(sync.Mutex),
	}
	Conns[ConnId] = c

	fmt.Println("Add conn", ConnId)

	return ConnId
}

func RemoveCon(id int) {
	conn := Conns[id]
	conn.I.Close()
	delete(Conns, id)

	fmt.Println("Remove conn", id)
}

func DispatchCommandByName(name string, command string) {

	for _, value := range Conns {
		if name == value.Name {
			go DispatchCommand(value, command)
		}
	}

}

func DispatchCommand(conn *Conn, command string) {
	conn.Mutex.Lock()
	message := &vcron.Message{
		Type:    proto.String("run"),
		Command: proto.String(command),
	}
	data, _ := proto.Marshal(message)
	conn.I.Write(data)
	conn.Mutex.Unlock()
}
