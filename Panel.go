package main

import (
	"github.com/yangmls/vcron/rest"
	"log"
	"net/http"
)

type Panel struct {
	Port string
}

func (panel *Panel) Run() {

	api := rest.NewApi()

	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(JobsRoutes()...)
	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("."))))

	log.Fatal(http.ListenAndServe(":"+panel.Port, nil))
}
