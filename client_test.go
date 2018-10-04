package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

type root struct {
	XMLName xml.Name `xml:"root"`
	Rows    []row    `xml:"row"`
}

type row struct {
	XMLName   xml.Name `xml:"row"`
	ID        int      `xml:"id"`
	Age       int      `xml:"age"`
	Gender    string   `xml:"gender"`
	LastName  string   `xml:"last_name"`
	FirstName string   `xml:"first_name"`
	About     string   `xml:"about"`
}

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
	if response != nil || err.Error() != "offset must be > 0" {
		t.Errorf("Shouldn't be able to proceed request with offset < 0\n")
	}
}

func queryXML(query string, limit string) []User {
	xmlFile, err := os.Open("dataset.xml")
	if err != nil {
		fmt.Printf("Can't open file: %v\n", err)
	}
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	var usersXML root
	err = xml.Unmarshal(byteValue, &usersXML)
	if err != nil {
		fmt.Printf("Can't parse XML: %v\n", err)
	}
	users := make([]User, 0)
	for _, user := range usersXML.Rows {
		name := user.FirstName + " " + user.LastName
		if strings.Contains(name, query) || strings.Contains(user.About, query) {
			var u User
			u.Id = user.ID
			u.About = user.About
			u.Age = user.Age
			u.Gender = user.Gender
			u.Name = name
			users = append(users, u)
		}
	}
	return users
}

func TestResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		query := r.URL.Query().Get("query")
		limit := r.URL.Query().Get("limit")
		users := queryXML(query, limit)
		b, err := json.Marshal(users)
		if err != nil {
			fmt.Printf("Can't marshal users: %v\n", err)
		}
		w.Write(b)
	}))
	client := &SearchClient{URL: server.URL}
	response, err := client.FindUsers(SearchRequest{Query: "Wolf", Limit: 17})
	if response == nil || err != nil {
		t.Errorf("Shouldn't be able to receive result\n")
	}
	response, err = client.FindUsers(SearchRequest{Query: "enim", Limit: 17})
	if response.NextPage != true || err != nil {
		t.Errorf("Should show next page true\n")
	}
}

func TestResponseUnknownError(t *testing.T) {
	httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	client := &SearchClient{}
	response, err := client.FindUsers(SearchRequest{})
	if response != nil || err.Error() != `unknown error Get ?limit=1&offset=0&order_by=0&order_field=&query=: unsupported protocol scheme ""` {
		t.Errorf("Shouldn't be able to receive result\n")
	}
}

func TestResponseError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
	}))
	client := &SearchClient{URL: server.URL}
	response, err := client.FindUsers(SearchRequest{})
	if response != nil || err.Error() != "timeout for limit=1&offset=0&order_by=0&order_field=&query=" {
		t.Errorf("Should be timeout\n")
	}
}
