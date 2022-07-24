package main

import (
	"fmt"
	"os"
	"os/signal"
	"server/router"
	"syscall"
)

var manager = router.Manager{
	Clients:         make(map[*router.Client]router.ClientDetail),
	InComingMessage: make(chan router.MessagePacket),
	Connect:         make(chan *router.Client),
	Disconnect:      make(chan *router.Client),
}

func main() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./virtualself-gseo-b664ca37d0c7.json")
	// ctx := context.Background() // Use context to cancel running child go routines after deadline of 450 miliseconds
	// ctx, cancel := context.WithCancel(ctx)

	port := getPort()
	fmt.Println("Connecting in port: " + port)
	// Create Server
	server, router := NewServer(
		"",
		port,
	)

	// go server.ListenAndServe(ctx)
	go server.ListenAndServe()
	go manager.Start(router)
	go manager.MonitorDetails()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	<-c
	close(c)
	fmt.Println("Closing...")

	// cancel() // Close all child routines
}

func getPort() string {
	port, ok := os.LookupEnv("PORT")
	// Set a default port if there is nothing in the environment
	if !ok {
		port = "9090"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return port
}
