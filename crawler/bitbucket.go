package crawler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-redis/redis"
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
	var nextPage = host.URL

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	for {
		// AAA: aggiungi a redis K: nextPage e V: false
		err := redisClient.Set(nextPage, false, 0).Err()
		if err != nil {
			panic(err)
		}
		log.Debug("GetRepositories -> NextPage at false: " + nextPage)
		///
		body, status, err := httpclient.GetURL(nextPage) // TODO: from 1 to 4 seconds to retrieve this data. Bottleneck.
		if err != nil && status.StatusCode != http.StatusOK {
			return err
		}

		var result response
		json.Unmarshal(body, &result)

		for _, v := range result.Values {
			repositories <- Repository{
				Name: v.FullName,
				//URL:  v.Links.Clone[0].Href + "/raw/default/publiccode.yml",
				URL:    v.Links.Clone[0].Href + "/raw/default/.gitignore",
				Source: nextPage,
			}
		}

		// AAA: cambia a redis K: nextPage e V: after parsed.

		u, _ := url.Parse(nextPage)
		q := u.Query()
		value := q.Get("after")

		err = redisClient.Set(nextPage, value, 0).Err()
		if err != nil {
			panic(err)
		}
		log.Debug("GetRepositories -> Set at " + value + ": " + nextPage)
		///
		if len(result.Next) == 0 {
			if len(repositories) == 0 {
				// if i want to restart:
				//nextPage = "https://api.bitbucket.org/2.0/repositories?pagelen=100"
				//and comment "close(repositories)"
				log.Debug("End of repos")
				close(repositories)
			}
		} else {
			nextPage = result.Next
		}

	}
}
