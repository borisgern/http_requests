package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInvalidAccesToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	client := &SearchClient{URL: server.URL}
	response, err := client.FindUsers(SearchRequest{})
	if response != nil || err.Error() != "Bad AccessToken" {
		t.Errorf("Access without valid AccessToken should be prohibited\n")
	}
}

func TestInternalServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	client := &SearchClient{URL: server.URL}
	response, err := client.FindUsers(SearchRequest{})
	if response != nil || err.Error() != "SearchServer fatal error" {
		t.Errorf("Server internal error\n")
	}
}
