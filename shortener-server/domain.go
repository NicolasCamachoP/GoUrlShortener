package main

import "net/url"

type UrlRequest struct {
	TargetUrl string `json:"target_url"`
}

func (ur *UrlRequest) IsValid() bool {
	_, err := url.ParseRequestURI(ur.TargetUrl)
	return len(ur.TargetUrl) > 0 && err == nil
}

type UrlResponse struct {
	ShortenedUrl string `json:"shortened_url"`
}
