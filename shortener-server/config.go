package main

import "fmt"

type Config struct {
	DB_DOMAIN               string
	DB_PORT_NUMBER          int
	DB_PASSWORD             string
	DB_KEY_EXPIRATION_HOURS int
	SERVER_PORT_NUMBER      int
	SERVER_BASEPATH         string
}

func GetConfig() (*Config, error) {
	//TODO read and validate all config from EnvVars
	return nil, fmt.Errorf("not implemented")
}
