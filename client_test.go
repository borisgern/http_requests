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
	if response != nil || err.Error() != "cant unpack error json: invalid character 'b' looking for beginning of value" {
		t.Errorf("Shouldn't be able to parse bad json\n")
	}
}

func TestJSONUnmarshalResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("result"))
	}))
	client := &SearchClient{URL: server.URL}
	response, err := client.FindUsers(SearchRequest{})
	if response != nil || err.Error() != "cant unpack result json: invalid character 'r' looking for beginning of value" {
		t.Errorf("Shouldn't be able to parse bad json in result\n")
	}
}

func TestBadOrderField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"ErrorBadOrderField"}`))
	}))
	client := &SearchClient{URL: server.URL}
	response, err := client.FindUsers(SearchRequest{OrderField: "bad"})
	if response != nil || err.Error() != "OrderFeld bad invalid" {
		t.Errorf("Should return if OrderField is invalid\n")
	}
}

func TestUnknownBadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"unknown"}`))
	}))
	client := &SearchClient{URL: server.URL}
	response, err := client.FindUsers(SearchRequest{})
	if response != nil || err.Error() != "unknown bad request error: unknown" {
		t.Errorf("Should return if unknown error\n")
	}
}

func TestReqParams(t *testing.T) {
	client := &SearchClient{}
	response, err := client.FindUsers(SearchRequest{Limit: -1})
	if response != nil || err.Error() != "limit must be > 0" {
		t.Errorf("Shouldn't be able to proceed request with limit < 0\n")
	}
	response, err = client.FindUsers(SearchRequest{Limit: 26, Offset: -1})
	//fmt.Println(response, err)
	if response != nil || err.Error() != "offset must be > 0" {
		t.Errorf("Shouldn't be able to proceed request with offset < 0\n")
	}
}
