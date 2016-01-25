package main

var exit = make(chan bool)

func main() {
	go StartPanel()
	go StartCron()
	go StartServer()

	<-exit
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
	LoadJobs()
}
