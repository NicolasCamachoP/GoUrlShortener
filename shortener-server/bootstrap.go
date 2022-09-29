package main

import (
	"log"
	"os"
)

func main() {
	log.Println("[INFO] - Injecting dependencies")
	server, err := Inject()
	if err != nil {
		log.Println("[ERROR] - Error while injecting dependencies: ", err)
		os.Exit(1)
	}
	log.Println("[INFO] - Starting server")
	if err = server.Start(); err != nil {
		log.Println("[ERROR] - Error while starting server: ", err)
		os.Exit(1)
	}

}
