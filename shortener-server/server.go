package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-http-utils/headers"
	"github.com/gorilla/mux"
)

const (
	KEY_ID = "KEY"
)

//Point to the handler to be able to shut down the service completely
var handler *Handler

type Server struct {
	options     *ServerOptions
	httpHandler http.Handler
}

type ServerOptions struct {
	PortNumber int
	BasePath   string //Will be used in both GET and POST requests
}

type Handler struct {
	shortener ShortenerInterface
	basePath  string
}

func NewServer(serverOpts *ServerOptions, shortener ShortenerInterface) *Server {
	router := mux.NewRouter()
	handler = &Handler{shortener, serverOpts.BasePath}
	router.Handle(fmt.Sprintf("%v/", serverOpts.BasePath), handler)
	router.Handle(fmt.Sprintf("%v/{%v}", serverOpts.BasePath, KEY_ID), handler)
	return &Server{serverOpts, router}
}

func (s *Server) Start() error {
	log.Printf("[INFO] [Start] - Listening and serving for port %v and basepath %s", s.options.PortNumber, s.options.BasePath)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", s.options.PortNumber), s.httpHandler); err != nil {
		return fmt.Errorf("error while doing ListenAndServe operation: %w", err)
	}
	return nil
}

func (s *Server) ShutDown() error {
	if handler != nil {
		return handler.shortener.ShutDown()
	}
	return fmt.Errorf("handler is nil")
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getUrl(w, r)
	case http.MethodPost:
		h.postUrl(w, r)
	default:
		w.Header().Set("Allow", fmt.Sprintf("%v, %v", http.MethodGet, http.MethodPost))
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getUrl(w http.ResponseWriter, r *http.Request) {
	key, ok := mux.Vars(r)[KEY_ID]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Missing key parameter...")
		return
	}
	associatedUrl, err := h.shortener.GetUrl(key)
	if err != nil {
		log.Println("[ERROR] [GET] - Error searching for url: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if associatedUrl == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set(headers.Location, associatedUrl)
	w.Header().Set(headers.ContentLocation, "absolute-URI")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) postUrl(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("[ERROR] - Error reading body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	request := &UrlRequest{}

	if err = json.Unmarshal(body, request); err != nil {
		log.Println("[ERROR] - Error unmarshalling body: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !request.IsValid() {
		log.Println("[ERROR] - Incomplete body or wrong url format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := h.shortener.SaveUrl(request.TargetUrl)
	if err != nil {
		log.Println("[ERROR] - Error saving url: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := &UrlResponse{fmt.Sprintf("%s%s%s", r.Host, r.URL.Path, id)}
	rawResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("[ERROR] - Error marshalling response: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(headers.ContentType, "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(rawResponse))
}
