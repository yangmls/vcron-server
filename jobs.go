package main

import (
	"encoding/json"
	"fmt"
	"github.com/yangmls/vcron"
	"io/ioutil"
	"os/user"
	"time"
)

var (
	JobId = 0
	Jobs  = make(map[int]*Job)
)

type Job struct {
	Name       string
	Expression string
	Command    string
	Timer      *time.Timer
}

type JobStore struct {
	Name       string
	Expression string
	Command    string
}

func AddJob(n string, e string, c string) int {
	job := &Job{
		Name:       n,
		Expression: e,
		Command:    c,
	}

	JobId = JobId + 1
	Jobs[JobId] = job
	go StartJob(JobId)
	go StoreJobs()
	return JobId
}

func RemoveJob(id int) {
	StopJob(id)
	delete(Jobs, id)
	go StoreJobs()
}

func UpdateJob(id int, n string, e string, c string) {
	StopJob(id)
	job := Jobs[id]
	job.Name = n
	job.Expression = e
	job.Command = c
	go StartJob(id)
	go StoreJobs()
}

func StartJob(id int) {
	job := Jobs[id]
	cron := *vcron.NewCron(job.Expression)

	for {
		job.Timer = cron.GetNextTimer()
		<-job.Timer.C
		go DispatchCommandByName(job.Name, job.Command)
	}
}

func StopJob(id int) {
	job := Jobs[id]
	job.Timer.Stop()
}

func JobsPath() string {
	var (
		current *user.User
		err     error
	)

	current, err = user.Current()

	if err != nil {
		fmt.Println("can not get home dir")
	}

	dir := current.HomeDir
	return dir + "/vcron.json"
}

func LoadJobs() {
	path := JobsPath()
	b, err := ioutil.ReadFile(path)

	if err != nil {
		return
	}

	var data []*JobStore
	json.Unmarshal(b, &data)

	for _, row := range data {
		job := &Job{
			Name:       row.Name,
			Expression: row.Expression,
			Command:    row.Command,
		}

		JobId = JobId + 1
		Jobs[JobId] = job
		go StartJob(JobId)
	}
}

func StoreJobs() {
	path := JobsPath()

	data := make([]*JobStore, len(Jobs))

	i := 0
	for _, job := range Jobs {
		data[i] = &JobStore{
			Name:       job.Name,
			Expression: job.Expression,
			Command:    job.Command,
		}
		i++
	}

	b, _ := json.Marshal(&data)
	ioutil.WriteFile(path, b, 0644)
}
