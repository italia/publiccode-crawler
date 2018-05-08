package crawler

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/italia/developers-italia-backend/httpclient"
	"github.com/prometheus/common/log"
)

// Gitlab is a Crawler for the Gitlab API.
type Gitlab []struct {
	ID                int           `json:"id"`
	Description       string        `json:"description"`
	Name              string        `json:"name"`
	NameWithNamespace string        `json:"name_with_namespace"`
	Path              string        `json:"path"`
	PathWithNamespace string        `json:"path_with_namespace"`
	CreatedAt         string        `json:"created_at"`
	DefaultBranch     string        `json:"default_branch"`
	TagList           []interface{} `json:"tag_list"`
	SSHURLToRepo      string        `json:"ssh_url_to_repo"`
	HTTPURLToRepo     string        `json:"http_url_to_repo"`
	WebURL            string        `json:"web_url"`
	AvatarURL         interface{}   `json:"avatar_url"`
	StarCount         int           `json:"star_count"`
	ForksCount        int           `json:"forks_count"`
	LastActivityAt    string        `json:"last_activity_at"`
}

// RegisterGitlabAPI register the crawler function for Gitlab API.
func RegisterGitlabAPI() func(domain Domain, url string, wg *sync.WaitGroup, repositories chan Repository) (string, error) {
	return func(domain Domain, url string, wg *sync.WaitGroup, repositories chan Repository) (string, error) {
		log.Debugf("RegisterGitlabAPI: %s ")

		// Set BasicAuth header
		headers := make(map[string]string)
		if domain.BasicAuth != nil {
			rand.Seed(time.Now().Unix())
			n := rand.Int() % len(domain.BasicAuth)
			headers["Authorization"] = "Basic " + domain.BasicAuth[n]
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
		var results Gitlab
		err = json.Unmarshal(resp.Body, &results)
		if err != nil {
			return url, err
		}

		// Add repositories to the channel that will perform the check on everyone.
		for _, v := range results {
			if v.DefaultBranch != "" {
				repositories <- Repository{
					Name:       v.PathWithNamespace,
					FileRawURL: "https://gitlab.com/" + v.PathWithNamespace + "/raw/" + v.DefaultBranch + "/" + os.Getenv("CRAWLED_FILENAME"),
					Domain:     domain.Id,
					Headers:    headers,
				}
			}
		}

		if len(resp.Headers.Get("Link")) == 0 {
			for len(repositories) != 0 {
				time.Sleep(time.Second)
			}
			log.Info("Gitlab repositories status: end reached.")

			// If Restart: uncomment next line.
			// return "domain.URL", nil
			return "", nil

		}

		// Return next url
		parsedLink := httpclient.NextHeaderLink(resp.Headers.Get("Link"))
		if parsedLink == "" {
			log.Info("Gitlab repositories status: end reached (no more ref=Next header).")
			return domain.URL, nil
		}

		return parsedLink, nil
	}
}
