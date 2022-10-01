package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-http-utils/headers"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestHandler_ServeHTTP(t *testing.T) {
	handlerFn := (&Handler{shortener: &MockShotener{}}).ServeHTTP
	httpHandler := http.HandlerFunc(handlerFn)
	muxHandler := mux.NewRouter()
	muxHandler.HandleFunc(fmt.Sprintf("/test/{%s}", KEY_ID), handlerFn)

	//--- Unsuported PATCH method
	patchReq, err := http.NewRequest(http.MethodPatch, "/", nil)
	require.NoError(t, err, "unexpected error")
	require.NotNil(t, patchReq, "request expected")
	rr := httptest.NewRecorder()
	httpHandler.ServeHTTP(rr, patchReq)
	require.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	//--- Supported POST method (postUrl test)
	jsonBody := []byte(`{"target_url": "https://google.com"}`)
	bodyReader := bytes.NewReader(jsonBody)
	postReq, err := http.NewRequest(http.MethodPost, "/test", bodyReader)
	require.NoError(t, err, "unexpected error")
	require.NotNil(t, postReq, "request expected")
	rr = httptest.NewRecorder()
	httpHandler.ServeHTTP(rr, postReq)
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, `{"shortened_url":"/test/test"}`, rr.Body.String())

	//--- Supported GET method (getUrl test)
	getReq, err := http.NewRequest(http.MethodGet, "/test/test", nil)
	require.NoError(t, err, "unexpected error")
	require.NotNil(t, getReq, "request expected")
	rr = httptest.NewRecorder()
	muxHandler.ServeHTTP(rr, getReq)
	require.Equal(t, http.StatusTemporaryRedirect, rr.Code)
	require.Equal(t, "https://google.com", rr.Header().Get(headers.Location))
	require.Equal(t, "absolute-URI", rr.Header().Get(headers.ContentLocation))
}

func TestServer_ShutDown(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		handler *Handler
	}{
		{
			name:    "TestServer_ShutDown - Empty handler",
			handler: nil,
			wantErr: true,
		},
		{
			name: "TestServer_ShutDown - Empty handler",
			handler: &Handler{
				shortener: &MockShotener{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{}
			handler = tt.handler
			if err := s.ShutDown(); (err != nil) != tt.wantErr {
				t.Errorf("Server.ShutDown() error = %v, wantErr %v", err, tt.wantErr)
			}
			handler = nil
		})
	}
}

//--------------------- Mocks ---------------------//

type MockShotener struct{ ShortenerInterface }

func (ms *MockShotener) ShutDown() error { return nil }

func (ms *MockShotener) GetUrl(key string) (string, error) {
	if key == "test" {
		return "https://google.com", nil
	}
	return "", fmt.Errorf("unknown key")
}

func (ms *MockShotener) SaveUrl(targetUrl string) (string, error) {
	return "/test", nil
}
