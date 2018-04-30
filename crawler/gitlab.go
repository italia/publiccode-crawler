package crawler

import (
	"github.com/italia/developers-italia-backend/httpclient"
	"net/http"
	"github.com/prometheus/common/log"
	"errors"
	"encoding/json"
	"os"
	"time"
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

// GetRepositories retrieves the list of all repository from a domain.
// Return the URL from where it should restart (Next or actual if fails) and error.
func RegisterGitlabAPI() func(domain Domain, url string, repositories chan Repository) (string, error) {
	return func(domain Domain, url string, repositories chan Repository) (string, error) {
		// Set BasicAuth header
		headers := make(map[string]string)
		if domain.BasicAuth != "" {
			headers["Authorization"] = "Basic " + domain.BasicAuth
		}

		// Get List of repositories
		body, status, respHeaders, err := httpclient.GetURL(url, headers)
		if err != nil {
			return url, err
		}
		if status.StatusCode != http.StatusOK {
			log.Warnf("Request returned: %s", string(body))
			return url, errors.New("requets returned an incorrect http.Status: " + status.Status)
		}

		// Fill response as list of values (repositories data).
		var results Gitlab
		err = json.Unmarshal(body, &results)
		if err != nil {
			return url, err
		}

		// Add repositories to the channel that will perform the check on everyone.
		for _, v := range results {
			repositories <- Repository{
				Name:       v.PathWithNamespace,
				FileRawURL: "https://gitlab.com/" + v.PathWithNamespace + "/raw/" + v.DefaultBranch + "/" + os.Getenv("CRAWLED_FILENAME"),
				Domain:     domain.Id,
				Headers:    headers,
			}
		}

		if len(respHeaders.Get("Link")) == 0 {
			for len(repositories) != 0 {
				time.Sleep(time.Second)
			}
			// if wants to end the program when repo list ends (last page) decomment
			// close(repositories)
			// return url, nil
			log.Info("Gitlab repositories status: end reached.")

			// Restart.
			return domain.URL, nil
		}

		// Return next url
		parsedLink := httpclient.NextHeaderLink(respHeaders.Get("Link"))
		if parsedLink == "" {
			log.Info("Gitlab repositories status: end reached (no more ref=Next header). Restart from: " + domain.URL)
			return domain.URL, nil
		}

		return parsedLink, nil
	}
}
