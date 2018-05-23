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
// with delay of 60 seconds.
func handlerOneRepoListWithDelay(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(60 * time.Second)
	fmt.Fprint(w, "{\"data\": \"example data\"}")
}

// handlerEmpty print an empty response
func handlerEmpty(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "")
}

// handlerHeaderInResponse print something and return a specific Header.
func handlerHeaderInResponse(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("X-PowOfTwo", "4")
	w.Header().Write(w)

	fmt.Fprint(w, "{\"data\": \"example data\"}")
}

// TestGetUrl should test if a getUrl to a valid http resource will end without errors.
func TestGetUrl(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerOneRepoList))
	defer ts.Close()

	resp, _ := GetURL(ts.URL, nil)
	r := "{\"data\": \"example data\"}"
	if string(resp.Body) != r && resp.Status.Code == 200 {
		t.Errorf("Call was incorrect, got: %s, want: %v.", resp.Body, r)
	}
}

// TestIncorrectProtocolUrl should test if a getUrl to incorrect protocol url will fail.
func TestIncorrectProtocolUrl(t *testing.T) {
	resp, err := GetURL("hktp://incorrectprotocol.url", nil)

	if err == nil && resp.Status.Code != -1 {
		t.Errorf("TestInexistentUrlCall was incorrect, got error: %v", err)
	}
}

// TestInexistentUrl should test if a getUrl to inexistent url will fail.
func TestInexistentUrl(t *testing.T) {
	resp, err := GetURL("http://inexistent.url", nil)
	if err == nil && resp.Status.Code == -1 {
		t.Errorf("TestInexistentUrlCall was incorrect, got error: %v", err)
	}
}

// TestEmptyResponse should test if a getUrl to a valid http resource with empty response will
// end without errors. Not really useful for getUrl(), use it as example.
func TestEmptyResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerEmpty))
	defer ts.Close()

	resp, _ := GetURL(ts.URL, nil)
	r := ""
	if string(resp.Body) != r {
		t.Errorf("TestEmptyCall was incorrect, got: %s, want: %v.", string(resp.Body), r)
	}
}

// TestGetUrlWithDelayResponse getUrl test with 10 seconds response time.
func TestGetUrlWithDelayResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerOneRepoListWithDelay))
	defer ts.Close()

	resp, _ := GetURL(ts.URL, nil)
	if resp.Body != nil {
		t.Errorf("TestCallWithDelay was incorrect, got: %s, want: %v.", string(resp.Body), nil)
	}
}

// TestGetUrlWithHeadersResponse getUrl test with headers in return.
func TestGetUrlWithHeadersResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerHeaderInResponse))
	defer ts.Close()

	resp, _ := GetURL(ts.URL, nil)
	if resp.Headers.Get("X-PowOfTwo") != "4" {
		t.Errorf("TestCallWithDelay was incorrect, got: %s, want: %s.", resp.Headers.Get("X-PowOfTwo"), "4")
	}
}
