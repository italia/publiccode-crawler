package httpclient

import (
	"math"
	"net/http"
	"time"

	version "github.com/italia/developers-italia-backend/crawler/version"
	log "github.com/sirupsen/logrus"
	"github.com/tomnomnom/linkheader"
)

// HTTPResponse wraps body, Status and Headers from the http.Response.
type HTTPResponse struct {
	Body    []byte
	Status  ResponseStatus
	Headers http.Header
}

const (
	userAgent = "Golang_italia_backend_bot"
)

// GetURL retrieves data, status and response headers from an URL.
// It uses some technique to slow down the requests if it get a 429 (Too Many Requests) response.
func GetURL(URL string, headers map[string]string) (HTTPResponse, error) {
	expBackoffAttempts := 0
	const timeout = 60 * time.Second
	const maxBackOffAttempts = 8 // 2 minutes.
	var err error

	client := http.Client{
		// Request Timeout.
		Timeout: timeout,
	}

	for expBackoffAttempts < maxBackOffAttempts {
		log.Debugf("attempt %d for %s", expBackoffAttempts, URL)

		req, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			return HTTPResponse{
				Body:    nil,
				Status:  ResponseStatus{Text: err.Error() + URL, Code: -1},
				Headers: nil,
			}, err
		}

		// Set headers.
		for k, v := range headers {
			req.Header.Add(k, v)
		}

		// Set special user agent for bot. Note: in github reqs the User-Agent must be set.
		req.Header.Add("User-Agent", userAgent+"/"+version.VERSION)

		// Perform the request.
		resp, err := client.Do(req)
		if err != nil {
			return HTTPResponse{
				Body:    nil,
				Status:  ResponseStatus{Text: err.Error() + URL, Code: -1},
				Headers: nil,
			}, err
		}

		// Check if the request results in http OK.
		if resp.StatusCode == http.StatusOK {
			return statusOK(resp)
		}

		// Check if the request results in http notFound.
		if resp.StatusCode == http.StatusNotFound {
			log.Debugf("Status: %s - Resource: %s", resp.Status, URL)
			return statusNotFound(resp)
		}

		// Check if the request results in http RateLimit error.
		if resp.StatusCode == http.StatusTooManyRequests {
			log.Debugf("Status: %s - Resource: %s", resp.Status, URL)
			expBackoffAttempts, err = statusTooManyRequests(resp, expBackoffAttempts)
			if err != nil {
				return HTTPResponse{
					Body:    nil,
					Status:  ResponseStatus{Text: err.Error() + URL, Code: -1},
					Headers: nil,
				}, err
			}

		}
		// Check if the request result in http Forbidden status.
		if resp.StatusCode == http.StatusForbidden {
			log.Debugf("Status: %s - Resource: %s", resp.Status, URL)
			expBackoffAttempts, err = statusForbidden(resp, expBackoffAttempts)
			if err != nil {
				return HTTPResponse{
					Body:    nil,
					Status:  ResponseStatus{Text: err.Error() + URL, Code: -1},
					Headers: nil,
				}, err
			}
		}

		defer resp.Body.Close()
		log.Errorf("Unknown http status code: %s - Resource: %s", resp.Status, URL)
		// increment attempts
		expBackoffAttempts++
	}

	// Generic invalid status code.
	return HTTPResponse{
		Body:    nil,
		Status:  ResponseStatus{Text: "Invalid Status Code: " + URL, Code: -1},
		Headers: nil,
	}, err
}

// HeaderLink parse the Github Header Link to "next"/"last"/"first"/"prev" link of repositories.
// Example: HeaderLink(link,"next") or HeaderLink(link, "prev") or HeaderLink(link,"last").
func HeaderLink(linkHeader, command string) string {
	parsedLinks := linkheader.Parse(linkHeader)

	for _, link := range parsedLinks {
		if link.Rel == command {
			return link.URL
		}
	}

	return ""
}

// expBackoffCalc calculate the exponential backoff given.
func expBackoffCalc(attempts int) float64 {
	return (math.Pow(2, float64(attempts)) - 1) / 2
}
