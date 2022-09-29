package main

import "fmt"

func Inject() (*Server, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, fmt.Errorf("error while getting configuration: %w", err)
	}
	dbOptions := &DbOptions{
		Database:      config.DB_DATABASE,
		UrlCollection: config.DB_COLLECTION,
		UserName:      config.DB_USERNAME,
		Password:      config.DB_PASSWORD,
		Host:          config.DB_HOST,
		Port:          config.DB_PORT_NUMBER,
	}
	shortenerOptions := &ShortenerOptions{}
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
	mongoSvc, err := NewMongoService(&DbOptions{})
	if err != nil {
		return nil, fmt.Errorf("error creating repository handler: %w", err)
	}
	shortener, err := NewShortener(shortenerOpts, mongoSvc)
	if err != nil {
		return nil, fmt.Errorf("error creating shortener: %w", err)
	}
	return NewServer(serverOpts, shortener), nil
}
