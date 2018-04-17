package crawler

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/italia/developers-italia-backend/httpclient"
	log "github.com/sirupsen/logrus"
)

// Bitbucket is a Crawler for the Bitbucket hosting.
type Bitbucket struct {
	URL       string
	RateLimit struct {
		ReqH int `yaml:"req/h"`
		ReqM int `yaml:"req/m"`
	} `yaml:"rate-limit"`
	BasicAuth string `yaml:"basic-auth"`
}

type response struct {
	Values []struct {
		Name  string `json:"name"`
		Links struct {
			Clone []struct {
				Href string `json:"href"`
				Name string `json:"name"`
			} `json:"clone"`
		} `json:"links"`
		FullName    string `json:"full_name"`
		Description string `json:"description"`
	} `json:"values"`
	Next string `json:"next"`
}

// GetRepositories retrieves the list of all repository from an hosting.
func (host Bitbucket) GetRepositories(repositories chan Repository) error {
	var sourceURL = host.URL

	// Redis connection.
	redisClient, err := redisClientFactory("redis:6379")
	if err != nil {
		log.Error(err)
		return err
	}

	for {
		// Save the current processed page as "not already processed".
		// Set a redis K: sourceURL e V: false
		err := redisClient.Set(sourceURL, false, 0).Err()
		if err != nil {
			log.Error(err)
			return err
		}
		log.Debugf("Redis: %s saved as 'false'.", sourceURL)

		// Set BasicAuth header
		headers := make(map[string]string)
		if host.BasicAuth != "" {
			headers["Authorization"] = "Basic " + host.BasicAuth
		}

		// Get List of repositories
		body, status, err := httpclient.GetURL(sourceURL, headers)
		if err != nil {
			return err
		}
		if status.StatusCode != http.StatusOK {
			log.Warnf("Request returned: %s", string(body))
			return errors.New("requets returned an incorrect http.Status: " + status.Status)
		}

		// Fill response as list of values (repositories data).
		var result response
		err = json.Unmarshal(body, &result)
		if err != nil {
			return err
		}

		// Add repositories to the channel that will perform the check on everyone.
		for _, v := range result.Values {
			repositories <- Repository{
				Name:    v.FullName,
				URL:     v.Links.Clone[0].Href + "/raw/default/" + os.Getenv("CRAWLED_FILENAME"),
				Source:  sourceURL,
				Headers: headers,
			}
		}

		// If reached, the page was correctly retrieved.
		// Set the value of sourceURL on redis to actual timestamp.
		timestamp := time.Now().String()
		err = redisClient.Set(sourceURL, timestamp, 0).Err()
		if err != nil {
			log.Error(err)
			return err
		}
		log.Debugf("Redis: set %s value as '%s'.", sourceURL, timestamp)

		// Check if the end of bitbucket repositories is reached and the repositories are all processed.
		if len(result.Next) == 0 {
			// If I want to restart when it ends:
			// sourceURL = "https://api.bitbucket.org/2.0/repositories?pagelen=100&after=2008-08-13"
			// and comment the line "close(repositories)"
			log.Info("Bitbucket repositories status: end reached.")
			close(repositories)

		} else {
			// Set the new URL to retrieve and continue.
			sourceURL = result.Next
		}

	}
}
