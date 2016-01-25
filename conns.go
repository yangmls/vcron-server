package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yangmls/vcron"
	"github.com/yangmls/vcron/cronexpr"
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
	Mutex *sync.Mutex
}

func (conn *Conn) Run() {
	for {
		now := time.Now()
		jobs := make([]*vcron.Job, 0)
		d := time.Duration(0)

		for _, job := range Jobs {
			if job.Name != conn.Name {
				continue
			}

			expression := cronexpr.MustParse(job.Expression)
			next := expression.Next(now)

			if d == time.Duration(0) && next.Sub(now) > d {
				d = next.Sub(now)
				j := &vcron.Job{
					Command: job.Command,
				}
				jobs = append(jobs, j)
			} else if d != time.Duration(0) && d == next.Sub(now) {
				j := &vcron.Job{
					Command: job.Command,
				}
				jobs = append(jobs, j)
			}
		}

		if d == time.Duration(0) {
			continue
		}

		conn.Timer = time.NewTimer(d)
		<-conn.Timer.C

		go conn.RunJobs(jobs)
	}
}

func (conn *Conn) Remove() {

}

func (conn *Conn) GetOrder() {

}

func (conn *Conn) IsFirst() {

}

func (conn *Conn) RunJobs(jobs []*vcron.Job) {
	conn.Mutex.Lock()

	request := &vcron.Request{
		Type: "run",
		Jobs: jobs,
	}

	conn.SendRequest(request)
	conn.WaitResponse()

	conn.Mutex.Unlock()
}

func (conn *Conn) SendRequest(request *vcron.Request) {
	fmt.Println("sending request")
	data, _ := proto.Marshal(request)
	conn.C.Write(data)
	fmt.Println("sent request")
}

func (conn *Conn) WaitResponse() (*vcron.Response, error) {
	buf := make([]byte, 2048)
	len, err := conn.C.Read(buf)

	if err != nil {
		return nil, err
	}

	data := buf[0:len]

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
		C:     c,
		Mutex: new(sync.Mutex),
	}
	Conns[ConnID] = conn
	ConnMutex.Unlock()

	request := &vcron.Request{
		Type: "register",
	}

	conn.SendRequest(request)
	response, err := conn.WaitResponse()

	if err != nil {
		return
	}

	conn.Name = response.Message

	fmt.Println("running ", ConnID)

	go conn.Run()
}
