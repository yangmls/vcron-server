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
	go func() {
		name := "Clives-Air.local"
		command := "date"
		expression := "0 * * * * * *"

		cron := *vcron.NewCron(expression)

		for {
			timer := cron.GetNextTimer()
			<-timer.C
			go DispatchCommandByName(name, command)
		}
	}()
}
