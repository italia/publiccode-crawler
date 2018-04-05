package crawler

import (
	"encoding/json"

	"github.com/italia/developers-italia-backend/httpclient"
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

	for {
		body, err := httpclient.GetURL(nextPage) // TODO: from 1 to 4 seconds to retrieve this data. Bottleneck.
		if err != nil {
			return err
		}

		var result response
		json.Unmarshal(body, &result)

		for _, v := range result.Values {
			repositories <- Repository{
				Name: v.FullName,
				//URL:  v.Links.Clone[0].Href + "/raw/default/publiccode.yml",
				URL: v.Links.Clone[0].Href + "/raw/default/.gitignore",
			}
		}

		if len(result.Next) == 0 {
			return nil
		}

		nextPage = result.Next
	}
}
