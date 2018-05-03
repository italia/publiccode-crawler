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

// Bitbucket is a Crawler for the Bitbucket API.
type Bitbucket struct {
	Pagelen int `json:"pagelen"`
	Values  []struct {
		Scm     string `json:"scm"`
		Website string `json:"website"`
		HasWiki bool   `json:"has_wiki"`
		Name    string `json:"name"`
		Links   struct {
			Watchers struct {
				Href string `json:"href"`
			} `json:"watchers"`
			Branches struct {
				Href string `json:"href"`
			} `json:"branches"`
			Tags struct {
				Href string `json:"href"`
			} `json:"tags"`
			Commits struct {
				Href string `json:"href"`
			} `json:"commits"`
			Clone []struct {
				Href string `json:"href"`
				Name string `json:"name"`
			} `json:"clone"`
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			Source struct {
				Href string `json:"href"`
			} `json:"source"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
			Hooks struct {
				Href string `json:"href"`
			} `json:"hooks"`
			Forks struct {
				Href string `json:"href"`
			} `json:"forks"`
			Downloads struct {
				Href string `json:"href"`
			} `json:"downloads"`
			Pullrequests struct {
				Href string `json:"href"`
			} `json:"pullrequests"`
		} `json:"links"`
		ForkPolicy string `json:"fork_policy"`
		UUID       string `json:"uuid"`
		Language   string `json:"language"`
		CreatedOn  string `json:"created_on"`
		Mainbranch struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"mainbranch"`
		FullName  string `json:"full_name"`
		HasIssues bool   `json:"has_issues"`
		Owner     struct {
			Username    string `json:"username"`
			DisplayName string `json:"display_name"`
			Type        string `json:"type"`
			UUID        string `json:"uuid"`
			Links       struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Avatar struct {
					Href string `json:"href"`
				} `json:"avatar"`
			} `json:"links"`
		} `json:"owner"`
		UpdatedOn   string `json:"updated_on"`
		Size        int    `json:"size"`
		Type        string `json:"type"`
		Slug        string `json:"slug"`
		IsPrivate   bool   `json:"is_private"`
		Description string `json:"description"`
		Project     struct {
			Key   string `json:"key"`
			Type  string `json:"type"`
			UUID  string `json:"uuid"`
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Avatar struct {
					Href string `json:"href"`
				} `json:"avatar"`
			} `json:"links"`
			Name string `json:"name"`
		} `json:"project,omitempty"`
		Parent struct {
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Avatar struct {
					Href string `json:"href"`
				} `json:"avatar"`
			} `json:"links"`
			Type     string `json:"type"`
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			UUID     string `json:"uuid"`
		} `json:"parent,omitempty"`
	} `json:"values"`
	Next string `json:"next"`
}

// RegisterBitbucketAPI register the crawler function for Bitbucket API.
func RegisterBitbucketAPI() func(domain Domain, url string, repositories chan Repository) (string, error) {
	return func(domain Domain, url string, repositories chan Repository) (string, error) {
		// Set BasicAuth header
		headers := make(map[string]string)
		if domain.BasicAuth != "" {
			headers["Authorization"] = "Basic " + domain.BasicAuth
		}

		// Get List of repositories
		resp, err := httpclient.GetURL(url, headers)
		if err != nil {
			return url, err
		}
		if resp.Status.Code != http.StatusOK {
			log.Warnf("Request returned: %s", string(resp.Body))
			return url, errors.New("request returned an incorrect http.Status: " + resp.Status.Text)
		}

		// Fill response as list of values (repositories data).
		var result Bitbucket
		err = json.Unmarshal(resp.Body, &result)
		if err != nil {
			return url, err
		}

		// Add repositories to the channel that will perform the check on everyone.
		for _, v := range result.Values {
			// If the repository was never used, the Mainbranch is empty ("")
			if v.Mainbranch.Name != "" {
				repositories <- Repository{
					Name:       v.FullName,
					FileRawURL: v.Links.HTML.Href + "/raw/" + v.Mainbranch.Name + "/" + os.Getenv("CRAWLED_FILENAME"),
					Domain:     domain.Id,
					Headers:    headers,
				}
			}
		}

		// Bitbucket end reached.
		if len(result.Next) == 0 {
			for len(repositories) != 0 {
				time.Sleep(time.Second)
			}
			log.Info("Bitbucket repositories status: end reached. Restart from domain value:" + domain.URL)

			// Restart.
			return domain.URL, nil
		}

		// Return next url
		return result.Next, nil
	}
}
