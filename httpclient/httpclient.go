package httpclient

import (
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// GetURL retrieves data from an URL.
// It uses some technique to slow down the requests if it get a 429 (Too Many Requests) response.
func GetURL(URL string) ([]byte, error) {
	var sleep time.Duration
	const timeout = time.Duration(10 * time.Second)

	client := http.Client{
		Timeout: timeout,
	}

	for {
		resp, err := client.Get(URL)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == 429 {
			sleep = sleep + (5 * time.Minute)
			log.Infof("Rate limit reached, sleep %v minutes\n", sleep)
			time.Sleep(sleep)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body.Close()

			return body, nil
		}
	}
}
