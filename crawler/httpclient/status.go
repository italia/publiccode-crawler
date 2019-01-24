package httpclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// ResponseStatus contains the status and statusCode of a response.
type ResponseStatus struct {
	Text string // e.g. "200 OK"
	Code int    // e.g. 200
}

const (
	headerRetryAfter    = "Retry-After"
	headerRateReset     = "X-RateLimit-Reset"
	headerRateRemaining = "X-RateLimit-Remaining"
)

// statusOK returns an HTTPResponse with the data from response.
func statusOK(resp *http.Response) (HTTPResponse, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf(err.Error())
		return HTTPResponse{
			Body:    nil,
			Status:  ResponseStatus{Text: resp.Status, Code: resp.StatusCode},
			Headers: resp.Header,
		}, err
	}

	err = resp.Body.Close()
	if err != nil {
		log.Errorf(err.Error())
	}

	return HTTPResponse{
		Body:    body,
		Status:  ResponseStatus{Text: resp.Status, Code: resp.StatusCode},
		Headers: resp.Header,
	}, nil
}

// statusNotFound returns an HTTPResponse with the data from response.
func statusNotFound(resp *http.Response) (HTTPResponse, error) {
	return HTTPResponse{
		Body:    nil,
		Status:  ResponseStatus{Text: resp.Status, Code: resp.StatusCode},
		Headers: resp.Header,
	}, fmt.Errorf("not found")
}

// statusTooManyRequests returns an HTTPResponse with the data from response.
func statusTooManyRequests(resp *http.Response, expBackoffAttempts int) (int, error) {
	// If Retry-after Header is set, use the header value.
	if retryAfter := resp.Header.Get(headerRetryAfter); retryAfter != "" {
		log.Infof("Waiting: %s seconds. (The value of %s)", retryAfter, headerRetryAfter)
		secondsAfterRetry, err := strconv.Atoi(retryAfter)
		if err != nil {
			log.Warn(err)
		}
		time.Sleep(time.Second * time.Duration(secondsAfterRetry))
		return expBackoffAttempts, nil
	}
	// Calculate ExpBackoff
	expBackoffWait := expBackoffCalc(expBackoffAttempts)
	// Perform a backoff sleep time.
	sleep := time.Duration(expBackoffWait) * time.Second
	log.Infof("Rate limit reached, sleep %v \n", sleep)
	time.Sleep(sleep)

	return expBackoffAttempts + 1, nil
}

// statusForbidden returns an HTTPResponse with the data from response.
func statusForbidden(resp *http.Response, expBackoffAttempts int) (int, error) {
	// If Retry-after is set, use that value.
	if retryAfter := resp.Header.Get(headerRetryAfter); retryAfter != "" {
		log.Infof("Waiting: %s seconds. (The value of %s)", retryAfter, headerRetryAfter)
		secondsAfterRetry, err := strconv.Atoi(retryAfter)
		if err != nil {
			log.Warn(err)
		}
		time.Sleep(time.Second * time.Duration(secondsAfterRetry))
		return expBackoffAttempts, nil
	}

	// If X-rateLimit-remaining
	if reset := resp.Header.Get(headerRateReset); reset != "" {
		// If X-RateLimit-Remaining is set
		if remaining := resp.Header.Get(headerRateRemaining); reset != "" {
			rateRemaining, err := strconv.Atoi(remaining)
			if err != nil {
				log.Warn(err)
			}
			if rateRemaining != 0 {
				// In this case there is another StatusForbidden and i should skip.
				return expBackoffAttempts, fmt.Errorf("forbidden resource")
			}

			retryEpoch, err := strconv.Atoi(reset)
			if err != nil {
				log.Warn(err)
			}
			secondsAfterRetry := int64(retryEpoch) - time.Now().Unix()
			log.Infof("Waiting %s seconds for %s. (The difference between header %s and time.Now())", strconv.FormatInt(secondsAfterRetry, 10), headerRateReset, reset)
			time.Sleep(time.Second * time.Duration(secondsAfterRetry))
			return expBackoffAttempts, nil
		}
	}

	// Generic forbidden.
	return expBackoffAttempts, fmt.Errorf("forbidden resource")

}
