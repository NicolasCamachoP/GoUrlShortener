package main

import (
	"log"
	"os"
	"os/signal"
)

var server *Server

func main() {
	//Channel to receive notifications
	channel := make(chan os.Signal, 1)
	//Notify when interrupt signal is received
	signal.Notify(channel, os.Interrupt)
	go func() {
		for sig := range channel {
			log.Printf("[INFO] - %v captured, cleaning up...", sig)
			if server != nil {
				server.ShutDown()
			}
			os.Exit(1)
		}
	}()
	log.Println("[INFO] - Injecting dependencies")
	var err error
	server, err = Inject()
	if err != nil {
		log.Println("[ERROR] - Error while injecting dependencies: ", err)
		os.Exit(1)
	}
	log.Println("[INFO] - Starting server")
	if err = server.Start(); err != nil {
		log.Println("[ERROR] - Error while starting server: ", err)
		log.Println("[INFO] - Shutting down services")
		if err = server.ShutDown(); err != nil {
			log.Println("[ERROR] - Unable to shut down services: ", err)
		}
		os.Exit(1)
	}
}
