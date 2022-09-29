package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

type ShortenerOptions struct{}
type Shortener struct {
	dbHandler IRepository
}

func NewShortener(shortenerOpts *ShortenerOptions, db IRepository) *Shortener {
	return &Shortener{db}
}

func (s *Shortener) ShutDown() error {
	if s.dbHandler != nil {
		return s.dbHandler.ShutDown()
	}
	return fmt.Errorf("unable to shutdown dbHandler. nil value")
}

func (s *Shortener) SaveUrl(targetUrl string) (string, error) {
	id := genId(targetUrl)
	if !s.dbHandler.Exists(id) {
		if err := s.dbHandler.SetItem(id, targetUrl); err != nil {
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
	} else if urlFound == nil {
		return "", nil
	}
	return fmt.Sprint(urlFound), nil
}
