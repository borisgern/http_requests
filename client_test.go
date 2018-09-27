package main

import (
	//"fmt"
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
		t.Errorf("Should be server internal error\n")
	}
}

func TestJSONUnmarshalError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad json"))
	}))
	client := &SearchClient{URL: server.URL}
	response, err := client.FindUsers(SearchRequest{})
	//fmt.Println(response, err)
	if response != nil || err.Error() != "cant unpack error json: invalid character 'b' looking for beginning of value" {
		t.Errorf("Shouldn't be able to parse bad json\n")
	}
}
