package main

import (
	"github.com/yangmls/vcron"
)

func main() {
	go StartPanel()
	go StartCron()

	StartServer()
}

func StartServer() {
	server := &Server{Port: "7023"}
	server.Run()
}

func StartPanel() {
	panel := &Panel{Port: "8080"}
	panel.Run()
}

func StartCron() {
	cron := *vcron.NewCron()

	go func() {
		name := "Clives-Air.local"
		command := "echo 1"

		cron.Add(name, "* * * * *", command)

		for {
			command := <-cron.Listeners[name]
			go DispatchCommandByName(name, command)
		}
	}()

	cron.Run()
}
