package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

type ShortenerOptions struct {
	UrlExpirationHours int
}
type Shortener struct {
	dbHandler DbHandlerInterface
	opts      *ShortenerOptions
}

func NewShortener(shortenerOpts *ShortenerOptions, db DbHandlerInterface) (*Shortener, error) {

	return nil, nil
}

func (s *Shortener) SaveUrl(targetUrl string) (string, error) {
	id := genId(targetUrl)
	if !s.dbHandler.Exists(id) {
		err := s.dbHandler.SetItem(id, targetUrl, time.Duration(s.opts.UrlExpirationHours)*time.Hour)
		if err != nil {
			return "", fmt.Errorf("error while saving url: %w", err)
		}
	}
	return id, nil
}

func genId(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func (s *Shortener) GetUrl(id string) (string, error) {
	urlFound, err := s.dbHandler.GetItem(id)
	if err != nil {
		return "", fmt.Errorf("error while retrieving url: %w", err)
	} else if urlFound == "" {
		return "", nil
	}
	return fmt.Sprintf("%v", urlFound), nil
}
