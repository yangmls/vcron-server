package main

import (
	"github.com/yangmls/vcron"
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

func AddJob(job *Job) int {
	JobId = JobId + 1
	Jobs[JobId] = job
	go StartJob(JobId)

	return JobId
}

func RemoveJob(id int) {
	StopJob(id)
	delete(Jobs, id)
}

func UpdateJob(id int, job *Job) {
	StopJob(id)
	Jobs[id] = job
	go StartJob(id)
}

func StartJobs() {
	for id, _ := range Jobs {
		go StartJob(id)
	}
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
