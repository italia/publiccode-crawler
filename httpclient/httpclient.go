package httpclient

import (
	"math"
	"net/http"
	"time"

	version "github.com/italia/developers-italia-backend/version"
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

	}

	// Generic invalid status code.
	return HTTPResponse{
		Body:    nil,
		Status:  ResponseStatus{Text: "Invalid Status Code: " + URL, Code: -1},
		Headers: nil,
	}, err
}

// NextHeaderLink parse the Github Header Link to next link of repositories.
// original Link: <https://api.github.com/repositories?since=1592>; rel="next", <https://api.github.com/repositories{?since}>; rel="first"
// parsedLink: https://api.github.com/repositories?since=1592
func NextHeaderLink(linkHeader string) string {
	parsedLinks := linkheader.Parse(linkHeader)

	for _, link := range parsedLinks {
		if link.Rel == "next" {
			return link.URL
		}
	}

	return ""
}

// HeaderLink parse the Github Header Link to "next"/"last"/"first"/"prev" link of repositories.
// HeaderLink("next", link) or HeaderLink("prev", link) or HeaderLink("last", link).
func HeaderLink(command, linkHeader string) string {
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
