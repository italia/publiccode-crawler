package httpclient

import (
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// RespStatus contains the status and statusCode of a response.
type RespStatus struct {
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200
}

// GetURL retrieves data and statusCode from an URL.
// It uses some technique to slow down the requests if it get a 429 (Too Many Requests) response.
func GetURL(URL string, headers map[string]string) ([]byte, RespStatus, error) {
	var sleep time.Duration
	const timeout = time.Duration(20 * time.Second)

	client := http.Client{
		// Request Timeout.
		Timeout: timeout,
	}

	for {
		req, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			return nil, RespStatus{Status: err.Error(), StatusCode: -1}, err
		}

		// Set headers.
		for k, v := range headers {
			req.Header.Add(k, v)
		}

		// Set special user agent for bot.
		req.Header.Add("User-Agent", "Golang_talia_backend_bot/0.0.1")

		// Perform the request.
		resp, err := client.Do(req)
		if err != nil {
			return nil, RespStatus{Status: err.Error(), StatusCode: -1}, err
		}

		// Check if the request results in http RateLimit error.
		if resp.StatusCode == http.StatusTooManyRequests {
			sleep = sleep + (5 * time.Minute)
			log.Info("Rate limit reached, sleep %v minutes\n", sleep)
			time.Sleep(sleep)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, RespStatus{Status: resp.Status, StatusCode: resp.StatusCode}, err
			}

			err = resp.Body.Close()
			if err != nil {
				log.Info("Error closing Body in httpclient.go\n")
			}

			return body, RespStatus{Status: resp.Status, StatusCode: resp.StatusCode}, nil
		}
	}
}
