package httpclient

import (
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"

	version "github.com/italia/developers-italia-backend/version"
	log "github.com/sirupsen/logrus"
	"github.com/tomnomnom/linkheader"
)

// HttpResponse wraps body, Status and Headers from the http.Response.
type HttpResponse struct {
	Body    []byte
	Status  ResponseStatus
	Headers http.Header
}

// ResponseStatus contains the status and statusCode of a response.
type ResponseStatus struct {
	Text string // e.g. "200 OK"
	Code int    // e.g. 200
}

const (
	userAgent = "Golang_italia_backend_bot"

	headerRetryAfter    = "Retry-After"
	headerRateReset     = "X-RateLimit-Reset"
	headerRateRemaining = "X-RateLimit-Remaining"
)

// GetURL retrieves data, status and response headers from an URL.
// It uses some technique to slow down the requests if it get a 429 (Too Many Requests) response.
func GetURL(URL string, headers map[string]string) (HttpResponse, error) {
	expBackoffAttempts := 0

	var sleep time.Duration
	const timeout = time.Duration(60 * time.Second)

	client := http.Client{
		// Request Timeout.
		Timeout: timeout,
	}

	for {
		req, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			return HttpResponse{
				Body:    nil,
				Status:  ResponseStatus{Text: err.Error(), Code: -1},
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
			return HttpResponse{
				Body:    nil,
				Status:  ResponseStatus{Text: err.Error(), Code: -1},
				Headers: nil,
			}, err
		}

		// Check if the request results in http notFound.
		if resp.StatusCode == http.StatusNotFound {
			return HttpResponse{
				Body:    nil,
				Status:  ResponseStatus{Text: resp.Status, Code: resp.StatusCode},
				Headers: resp.Header,
			}, nil

			// Check if the request results in http OK.
		}
		if resp.StatusCode == http.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorf(err.Error())
				return HttpResponse{
					Body:    nil,
					Status:  ResponseStatus{Text: resp.Status, Code: resp.StatusCode},
					Headers: resp.Header,
				}, err
			}

			err = resp.Body.Close()
			if err != nil {
				log.Errorf(err.Error())
			}

			return HttpResponse{
				Body:    body,
				Status:  ResponseStatus{Text: resp.Status, Code: resp.StatusCode},
				Headers: resp.Header,
			}, nil

			// Check if the request results in http RateLimit error.
		}
		if resp.StatusCode == http.StatusTooManyRequests {

			if retryAfter := resp.Header.Get(headerRetryAfter); retryAfter != "" {
				// If Retry-after is set, use that value.
				log.Infof("Waiting: %s seconds for %s. (The value of %s)", retryAfter, URL, headerRetryAfter)
				secondsAfterRetry, _ := strconv.Atoi(retryAfter)
				time.Sleep(time.Second * time.Duration(secondsAfterRetry))
			} else {
				// Calculate ExpBackoff
				expBackoffWait := expBackoffCalc(expBackoffAttempts)
				// Perform a backoff sleep time.
				sleep = time.Duration(expBackoffWait) * time.Second
				expBackoffAttempts = expBackoffAttempts + 1
				log.Info("Rate limit reached, sleep %v \n", sleep)
				time.Sleep(sleep)
			}

			// Check if the request result in http Forbidden status
		} else if resp.StatusCode == http.StatusForbidden {

			if retryAfter := resp.Header.Get(headerRetryAfter); retryAfter != "" {
				// If Retry-after is set, use that value.
				log.Infof("Waiting: %s seconds for %s. (The value of %s)", retryAfter, URL, headerRetryAfter)
				secondsAfterRetry, _ := strconv.Atoi(retryAfter)
				time.Sleep(time.Second * time.Duration(secondsAfterRetry))

			} else if reset := resp.Header.Get(headerRateReset); reset != "" {
				// If X-rateLimit-remaining
				if remaining := resp.Header.Get(headerRateRemaining); reset != "" {
					rateRemaining, _ := strconv.Atoi(remaining)
					if rateRemaining != 0 {
						// In this case there is another StatusForbidden and i should skip.
						log.Errorf("Forbidden error on %s.", URL)
						return HttpResponse{
							Body:    nil,
							Status:  ResponseStatus{Text: resp.Status, Code: resp.StatusCode},
							Headers: resp.Header,
						}, err
					} else {
						retryEpoch, _ := strconv.Atoi(reset)
						secondsAfterRetry := int64(retryEpoch) - time.Now().Unix()
						log.Infof("Waiting %s seconds. (The difference between header %s and time.Now())", strconv.FormatInt(secondsAfterRetry, 10), headerRateReset)
						time.Sleep(time.Second * time.Duration(secondsAfterRetry))
					}
				}
			} else {
				// Generic forbidden.
				log.Errorf("Forbidden error on %s.", URL)
				return HttpResponse{
					Body:    nil,
					Status:  ResponseStatus{Text: resp.Status, Code: resp.StatusCode},
					Headers: resp.Header,
				}, err
			}

		} else {
			// Generic invalid status code.
			// Calculate ExpBackoff
			expBackoffWait := expBackoffCalc(expBackoffAttempts)
			// Perform a backoff sleep time.
			sleep = time.Duration(expBackoffWait) * time.Second
			expBackoffAttempts = expBackoffAttempts + 1
			log.Infof("Invalid status code on %s : sleep %v \n", URL, sleep)
			time.Sleep(sleep)
		}
	}
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
