package crawler

import (
	"encoding/json"

	"github.com/italia/developers-italia-backend/httpclient"
)

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

// Bitbucket is a Crawler for the Bitbucket hosting.
type Bitbucket struct {
	URL string
}

// GetRepositories retrieves the list of all repository from an hosting.
func (c Bitbucket) GetRepositories(repositories chan Repository) error {
	var nextPage = c.URL

	for {
		body, err := httpclient.GetURL(nextPage)
		if err != nil {
			return err
		}

		var result response
		json.Unmarshal(body, &result)

		for _, v := range result.Values {
			repositories <- Repository{
				Name: v.Name,
				URL:  v.Links.Clone[0].Href + "/raw/default/publiccode.yml",
			}
		}

		if len(result.Next) == 0 {
			return nil
		}

		nextPage = result.Next
	}
}
