package main

import (
	"fmt"
	mockhttp "github.com/karupanerura/go-mock-http-response"
	"net/http"
	"testing"
)

const validAccessToken = "valid"
const invalidAccessToken = "invalid"
const baseURL = "http://..."

func mockResponse(statusCode int, headers map[string]string, body []byte) {
	client = mockhttp.NewResponseMock(statusCode, headers, body).MakeClient()
}

func TestInvalidAccesToken(t *testing.T) {
	client := &SearchClient{URL: baseURL, AccessToken: invalidAccessToken}
	mockResponse(http.StatusUnauthorized, map[string]string{"AccessToken": validAccessToken}, []byte("Invalid access token"))
	request := SearchRequest{}
	response, err := client.FindUsers(request)
	fmt.Printf("response %v, error %v\n", response, err)
	if response != nil && err != fmt.Errorf("Bad AccessToken") {
		t.Errorf("Access with invalid AccessToken %v should be prohibited\n", invalidAccessToken)
	}
}
