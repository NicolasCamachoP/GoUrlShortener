package main

import "fmt"

func Inject() (*Server, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, fmt.Errorf("error while getting configuration: %w", err)
	}
	dbOptions := &DbOptions{
		Domain:     config.DB_DOMAIN,
		PortNumber: config.DB_PORT_NUMBER,
		Password:   config.DB_PASSWORD,
	}
	shortenerOptions := &ShortenerOptions{
		UrlExpirationHours: config.DB_KEY_EXPIRATION_HOURS,
	}
	serverOptions := &ServerOptions{
		PortNumber: config.SERVER_PORT_NUMBER,
		BasePath:   config.SERVER_BASEPATH,
	}
	server, err := CreateService(serverOptions, shortenerOptions, dbOptions)
	if err != nil {
		return nil, fmt.Errorf("error while initializing handler: %w", err)
	}
	return server, nil
}

func CreateService(serverOpts *ServerOptions, shortenerOpts *ShortenerOptions, dbOpts *DbOptions) (*Server, error) {
	redisHandler, err := NewRedisHandler(dbOpts)
	if err != nil {
		return nil, fmt.Errorf("error creating redis handler: %w", err)
	}

	shortener, err := NewShortener(shortenerOpts, redisHandler)
	if err != nil {
		return nil, fmt.Errorf("error creating shortener: %w", err)
	}
	server := NewServer(serverOpts, shortener)
	if err != nil {
		return nil, fmt.Errorf("error starting server: %w", err)
	}
	return server, nil
}
