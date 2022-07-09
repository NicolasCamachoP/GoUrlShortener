package main

type ShortenerInterface interface {
	GetUrl(id string) (string, error)
	SaveUrl(targetUrl string) (string, error)
}
