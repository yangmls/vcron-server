package main

import (
	"github.com/yangmls/vcron/rest"
	"net/http"
)

type JobModel struct {
	Name       string
	Expression string
	Command    string
}

func JobsRoutes() []*rest.Route {
	return []*rest.Route{
		rest.Get("/jobs", getJobs),
		rest.Post("/jobs", postJobs),
		rest.Put("/jobs/:name", putJobs),
		rest.Delete("/jobs/:name", deleteJobs),
	}
}

func getJobs(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(Jobs)
}

func postJobs(w rest.ResponseWriter, r *rest.Request) {
	model := &JobModel{}
	err := r.DecodeJsonPayload(model)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if model.Name == "" {
		rest.Error(w, "job name required", 400)
		return
	}
	if model.Command == "" {
		rest.Error(w, "job command required", 400)
		return
	}
	if model.Expression == "" {
		rest.Error(w, "job expression required", 400)
		return
	}
	job := &Job{
		Expression: model.Expression,
		Command:    model.Command,
	}
	AddJob(model.Name, job)
	w.WriteJson(model)
}

func putJobs(w rest.ResponseWriter, r *rest.Request) {

}

func deleteJobs(w rest.ResponseWriter, r *rest.Request) {

}
