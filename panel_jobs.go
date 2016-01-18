package main

import (
	"github.com/yangmls/vcron/rest"
	"net/http"
	"strconv"
)

type JobModel struct {
	Id         int
	Name       string
	Expression string
	Command    string
}

func JobsRoutes() []*rest.Route {
	return []*rest.Route{
		rest.Get("/jobs", getJobs),
		rest.Post("/jobs", postJobs),
		rest.Put("/jobs/:id", putJobs),
		rest.Delete("/jobs/:id", deleteJobs),
	}
}

func getJobs(w rest.ResponseWriter, r *rest.Request) {
	jobs := make([]*JobModel, len(Jobs))

	i := 0
	for id, job := range Jobs {
		jobs[i] = &JobModel{
			Id:         id,
			Name:       job.Name,
			Expression: job.Expression,
			Command:    job.Command,
		}
		i++
	}

	w.WriteJson(&jobs)
}

func postJobs(w rest.ResponseWriter, r *rest.Request) {
	model := payload(w, r)

	if model == nil {
		return
	}

	id := AddJob(model.Name, model.Expression, model.Command)
	model.Id = id
	w.WriteJson(model)
}

func putJobs(w rest.ResponseWriter, r *rest.Request) {
	s := r.PathParam("id")
	id, err := strconv.Atoi(s)
	if err != nil {
		return
	}

	model := payload(w, r)

	if model == nil {
		return
	}

	UpdateJob(id, model.Name, model.Expression, model.Command)
	model.Id = id
	w.WriteJson(model)
}

func deleteJobs(w rest.ResponseWriter, r *rest.Request) {
	s := r.PathParam("id")
	id, err := strconv.Atoi(s)
	if err != nil {
		return
	}
	RemoveJob(id)
	w.WriteHeader(http.StatusOK)
}

func payload(w rest.ResponseWriter, r *rest.Request) *JobModel {
	model := &JobModel{}
	err := r.DecodeJsonPayload(model)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	if model.Name == "" {
		rest.Error(w, "job name required", 400)
		return nil
	}
	if model.Command == "" {
		rest.Error(w, "job command required", 400)
		return nil
	}
	if model.Expression == "" {
		rest.Error(w, "job expression required", 400)
		return nil
	}

	return model
}
