package main

import (
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/yangmls/vcron"
)

var ConnId int
var Conns map[int]net.Conn
var ConnNames map[int]string
var ConnChans map[int]chan int

func init() {
	ConnId = 0
	Conns = make(map[int]net.Conn)
	ConnNames = make(map[int]string)
	ConnChans = make(map[int]chan int)
}

func GetConnId() int {
	ConnId = ConnId + 1

	return ConnId
}

func AddConn(name string, conn net.Conn) {
	id := GetConnId()

	ConnNames[id] = name
	Conns[id] = conn
	ConnChans[id] = make(chan int)
}

func RemoveCon(conn net.Conn) {

}

func DispatchCommandByName(name string, command string) {

	for key, value := range ConnNames {
		if name == value {
			go DispatchCommand(Conns[key], command)
		}
	}

}

func DispatchCommand(conn net.Conn, command string) {
	message := &vcron.Message{
		Type:    proto.String("run"),
		Command: proto.String(command),
	}
	data, _ := proto.Marshal(message)
	conn.Write(data)
}
