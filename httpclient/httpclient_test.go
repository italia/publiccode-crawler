package httpclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// handlerOneRepoList print on ResponseWriter one element of response list.
func handlerOneRepoList(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "{\"data\": \"example data\"}")
}

// handlerOneRepoListWithDelay print on ResponseWriter one element of response.
// with delay of 15 seconds
func handlerOneRepoListWithDelay(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(30 * time.Second)
	fmt.Fprint(w, "{\"data\": \"example data\"}")
}

// handlerEmpty print an empty response
func handlerEmpty(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "")
}

// TestGetUrl should test if a getUrl to a valid http resource will end without errors.
func TestGetUrl(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerOneRepoList))
	defer ts.Close()

	data, status, _ := GetURL(ts.URL, nil)
	r := "{\"data\": \"example data\"}"
	if string(data) != r && status.StatusCode == 200 {
		t.Errorf("Call was incorrect, got: %s, want: %v.", data, r)
	}
}

// TestIncorrectProtocolUrl should test if a getUrl to incorrect protocol url will fail.
func TestIncorrectProtocolUrl(t *testing.T) {
	_, status, err := GetURL("hktp://incorrectprotocol.url", nil)

	if err == nil && status.StatusCode != -1 {
		t.Errorf("TestInexistentUrlCall was incorrect, got error: %v", err)
	}
}

// TestInexistentUrl should test if a getUrl to inexistent url will fail.
func TestInexistentUrl(t *testing.T) {
	_, status, err := GetURL("http://inexistent.url", nil)
	fmt.Printf("inexistent %+v", status)
	if err == nil {
		t.Errorf("TestInexistentUrlCall was incorrect, got error: %v", err)
	}
}

// TestEmptyResponse should test if a getUrl to a valid http resource with empty response will
// end without errors. Not really useful for getUrl(), use it as example.
func TestEmptyResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerEmpty))
	defer ts.Close()

	data, _, _ := GetURL(ts.URL, nil)
	r := ""
	if string(data) != r {
		t.Errorf("TestEmptyCall was incorrect, got: %s, want: %v.", data, r)
	}
}

// TestGetUrlWithDelayResponse getUrl test with 10 seconds response time.
func TestGetUrlWithDelayResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerOneRepoListWithDelay))
	defer ts.Close()

	data, _, _ := GetURL(ts.URL, nil)
	if data != nil {
		t.Errorf("TestCallWithDelay was incorrect, got: %s, want: %v.", data, nil)
	}
}
