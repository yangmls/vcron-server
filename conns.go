package main

import (
	"encoding/binary"
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
	ID      int
	Name    string
	Timer   *time.Timer
	C       net.Conn
	Mutex   *sync.Mutex
	Running bool
}

func (conn *Conn) Run() {
	fmt.Println("running ", conn.ID)

	conn.Running = true

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
		msg := <-conn.Timer.C

		if !conn.Running {
			return
		}

		fmt.Println(msg)

		go conn.RunJobs(jobs)
	}
}

func (conn *Conn) Stop() {
	conn.Running = false
	conn.Timer.Stop()
}

func (conn *Conn) Remove() {
	fmt.Println("removing conn", conn.ID)

	conn.Stop()
	conn.C.Close()
	delete(Conns, conn.ID)
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

	err := conn.SendRequest(request)

	if err == nil {
		_, err := conn.WaitResponse()

		if err != nil {
			conn.Remove()
		}
	} else {
		conn.Remove()
	}

	conn.Mutex.Unlock()
}

func (conn *Conn) SendRequest(request *vcron.Request) error {
	var (
		data []byte
		err  error
	)

	fmt.Println("sending request")

	data, err = proto.Marshal(request)
	prefix := make([]byte, 4, 4)
	size := uint64(len(data))
	binary.PutUvarint(prefix, size)

	_, err = conn.C.Write(prefix)

	if err != nil {
		return err
	}

	_, err = conn.C.Write(data)

	if err != nil {
		return err
	}

	fmt.Println("sent request")

	return nil
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
		ID:    ConnID,
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

	go conn.Run()
}
