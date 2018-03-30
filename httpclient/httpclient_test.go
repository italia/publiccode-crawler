package httpclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// handlerOneRepoList print on ResponseWriter one element of response list
func handlerOneRepoList(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "{\"data\": \"SUPERDATA\"}")
}

// handlerOneRepoListWithDelay print on ResponseWriter one element of response
// with delay of 15 seconds
func handlerOneRepoListWithDelay(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(15 * time.Second)
	fmt.Fprint(w, "{\"data\": \"SUPERDATA\"}")
}

// handlerEmpty print an empty response
func handlerEmpty(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "")
}

// TestGetUrl should test if a getUrl to a valid http resource will end without errors
func TestGetUrl(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerOneRepoList))
	defer ts.Close()

	data, _ := GetURL(ts.URL)
	r := "{\"data\": \"SUPERDATA\"}"

	if string(data) != r {
		t.Errorf("Call was incorrect, got: %s, want: %v.", data, r)
	}
}

// TestInexistentUrl should test if a getUrl to inexistent url will fail
func TestInexistentUrl(t *testing.T) {
	_, err := GetURL("inexistent.url")

	if err == nil {
		t.Errorf("TestInexistentUrlCall was incorrect, got: %v, want: %v.", nil, err)
	}
}

// TestEmptyResponse should test if a getUrl to a valid http resource with empty response will
// end without errors. Not really useful for getUrl(), use it as example.
func TestEmptyResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerEmpty))
	defer ts.Close()

	data, _ := GetURL(ts.URL)
	r := ""

	if string(data) != r {
		t.Errorf("TestEmptyCall was incorrect, got: %s, want: %v.", data, r)
	}
}

// TestGetUrlWithDelayResponse getUrl test with 10 seconds response time
func TestGetUrlWithDelayResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerOneRepoListWithDelay))
	defer ts.Close()

	data, _ := GetURL(ts.URL)

	if data != nil {
		t.Errorf("TestCallWithDelay was incorrect, got: %s, want: %v.", data, nil)
	}
}
