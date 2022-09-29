package main

type ShortenerInterface interface {
	//GetUrl: retrieves an URL by Id
	//Input: URL id
	//Output: URL
	GetUrl(id string) (string, error)
	//SaveUrl: saver and URL and creates associated ID
	//Input: URL to save
	//Output: created ID and error if any
	SaveUrl(targetUrl string) (string, error)
	//ShutDown: close all possible connections
	//Output: error if any
	ShutDown() error
}
