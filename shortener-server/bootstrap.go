package main

import (
	"log"
	"os"
)

const (
	DEFAULT_PORT        = 3333
	ID_QUERY_STRING_KEY = "id"
	TARGET_URL_KEY      = "target_url"
	MAIN_PATH           = "/url"
)

func main() {
	log.Println("[INFO] - Injecting dependencies")
	server, err := Inject()
	if err != nil {
		log.Println("[ERROR] - Error while injecting dependencies: ", err)
		os.Exit(1)
	}
	log.Println("[INFO] - Starting server")
	err = server.Start()
	if err != nil {
		log.Println("[ERROR] - Error while starting server: ", err)
		os.Exit(1)
	}

}
