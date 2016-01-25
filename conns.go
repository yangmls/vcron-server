package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yangmls/vcron"
	"io/ioutil"
	"net"
	"sync"
	"time"
)

var (
	ConnID    = 0
	ConnMutex = new(sync.Mutex)
	Conns     = make(map[int]*Conn)
)

type Conn struct {
	Name  string
	Timer *time.Timer
	C     net.Conn
}

func (conn *Conn) Run() {
}

func (conn *Conn) Remove() {

}

func (conn *Conn) GetOrder() {

}

func (conn *Conn) IsFirst() {

}

func (conn *Conn) SendRequest(request *vcron.Request) {
	fmt.Println("sending request")
	data, _ := proto.Marshal(request)
	conn.C.Write(data)
	fmt.Println("sent request")
}

func (conn *Conn) WaitResponse() (*vcron.Response, error) {
	data, err := ioutil.ReadAll(conn.C)

	if err != nil {
		return nil, err
	}

	response := &vcron.Response{}

	err = proto.Unmarshal(data, response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func AddConn(c net.Conn) {
	ConnMutex.Lock()
	ConnID++
	conn := &Conn{
		C: c,
	}
	Conns[ConnID] = conn
	ConnMutex.Unlock()

	request := &vcron.Request{
		Type: "register",
	}

	conn.SendRequest(request)
	conn.WaitResponse()

	fmt.Println("running ", ConnID)

	go conn.Run()
}
