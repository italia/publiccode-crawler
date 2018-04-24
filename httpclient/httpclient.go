package httpclient

import (
	"io/ioutil"
	"net/http"
	"time"

	version "github.com/italia/developers-italia-backend/version"
	log "github.com/sirupsen/logrus"
)

// ResponseStatus contains the status and statusCode of a response.
type ResponseStatus struct {
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200
}

// GetURL retrieves data, status and response headers from an URL.
// It uses some technique to slow down the requests if it get a 429 (Too Many Requests) response.
func GetURL(URL string, headers map[string]string) ([]byte, ResponseStatus, http.Header, error) {
	var sleep time.Duration
	const timeout = time.Duration(20 * time.Second)

	client := http.Client{
		// Request Timeout.
		Timeout: timeout,
	}

	for {
		req, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			return nil, ResponseStatus{Status: err.Error(), StatusCode: -1}, nil, err
		}

		// Set headers.
		for k, v := range headers {
			req.Header.Add(k, v)
		}

		// Set special user agent for bot. Note: in github reqs the User-Agent must be set.
		req.Header.Add("User-Agent", "Golang_italia_backend_bot/"+version.VERSION)

		// Perform the request.
		resp, err := client.Do(req)
		if err != nil {
			return nil, ResponseStatus{Status: err.Error(), StatusCode: -1}, nil, err
		}

		// Check if the request results in http RateLimit error.
		if resp.StatusCode == http.StatusTooManyRequests {
			sleep = sleep + (5 * time.Minute)
			log.Info("Rate limit reached, sleep %v minutes\n", sleep)
			time.Sleep(sleep)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, ResponseStatus{Status: resp.Status, StatusCode: resp.StatusCode}, resp.Header, err
			}

			err = resp.Body.Close()
			if err != nil {
				log.Info("Error closing Body in httpclient.go\n")
			}

			return body, ResponseStatus{Status: resp.Status, StatusCode: resp.StatusCode}, resp.Header, nil
		}
	}
}
