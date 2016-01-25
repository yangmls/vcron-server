package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
)

var (
	JobId = 0
	Jobs  = make(map[int]*Job)
)

type Job struct {
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
	go StoreJobs()
	return JobId
}

func RemoveJob(id int) {
	delete(Jobs, id)
	go StoreJobs()
}

func UpdateJob(id int, n string, e string, c string) {
	job := Jobs[id]
	job.Name = n
	job.Expression = e
	job.Command = c
	go StoreJobs()
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

	var data []*Job
	json.Unmarshal(b, &data)

	for _, row := range data {
		job := &Job{
			Name:       row.Name,
			Expression: row.Expression,
			Command:    row.Command,
		}

		JobId = JobId + 1
		Jobs[JobId] = job
	}
}

func StoreJobs() {
	path := JobsPath()

	data := make([]*Job, len(Jobs))

	i := 0
	for _, job := range Jobs {
		data[i] = job
		i++
	}

	b, _ := json.Marshal(&data)
	ioutil.WriteFile(path, b, 0644)
}
