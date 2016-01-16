package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Panel struct {
	Port string
}

func (panel *Panel) Run() {
	RegisterRoutes()

	fmt.Println("Panel is running on " + panel.Port)

	err := http.ListenAndServe(":"+panel.Port, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func RegisterRoutes() {
	http.HandleFunc("/hello", HelloServer)
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}
