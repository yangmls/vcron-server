package main

import (
	"github.com/yangmls/vcron"
	"time"
)

var (
	Jobs = make(map[string]*Job)
)

type Job struct {
	Expression string
	Command    string
	Timer      *time.Timer
}

func AddJob(name string, job *Job) {
	Jobs[name] = job
	go StartJob(name)
}

func RemoveJob(name string) {
	StopJob(name)
	delete(Jobs, name)
}

func UpdateJob(name string, job *Job) {
	StopJob(name)
	Jobs[name] = job
	go StartJob(name)
}

func StartJobs() {
	for name, _ := range Jobs {
		go StartJob(name)
	}
}

func StartJob(name string) {
	job := Jobs[name]
	cron := *vcron.NewCron(job.Expression)

	for {
		job.Timer = cron.GetNextTimer()
		<-job.Timer.C
		go DispatchCommandByName(name, job.Command)
	}
}

func StopJob(name string) {
	job := Jobs[name]
	job.Timer.Stop()
}
