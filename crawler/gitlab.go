package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/italia/developers-italia-backend/httpclient"
	log "github.com/sirupsen/logrus"
	linkheader "github.com/tomnomnom/linkheader"
)

// Gitlab is a Crawler for the Gitlab hosting.
type Gitlab struct {
	URL       string
	RateLimit struct {
		ReqH int `yaml:"req/h"`
		ReqM int `yaml:"req/m"`
	} `yaml:"rate-limit"`
	BasicAuth string `yaml:"basic-auth"`
}

type gitlabResponse []struct {
	ID                int           `json:"id"`
	Description       string        `json:"description"`
	Name              string        `json:"name"`
	NameWithNamespace string        `json:"name_with_namespace"`
	Path              string        `json:"path"`
	PathWithNamespace string        `json:"path_with_namespace"`
	CreatedAt         time.Time     `json:"created_at"`
	DefaultBranch     string        `json:"default_branch"`
	TagList           []interface{} `json:"tag_list"`
	SSHURLToRepo      string        `json:"ssh_url_to_repo"`
	HTTPURLToRepo     string        `json:"http_url_to_repo"`
	WebURL            string        `json:"web_url"`
	AvatarURL         interface{}   `json:"avatar_url"`
	StarCount         int           `json:"star_count"`
	ForksCount        int           `json:"forks_count"`
	LastActivityAt    time.Time     `json:"last_activity_at"`
}

// GetRepositories retrieves the list of all repository from an hosting.
// Return the URL from where it should restart (Next or actual if fails) and error.
func (host Gitlab) GetRepositories(url string, repositories chan Repository) (string, error) {
	fmt.Println("get repo:" + url)
	// Set BasicAuth header
	headers := make(map[string]string)
	if host.BasicAuth != "" {
		headers["Authorization"] = "Basic " + host.BasicAuth
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
	var results gitlabResponse
	err = json.Unmarshal(body, &results)
	if err != nil {
		return url, err
	}

	// Add repositories to the channel that will perform the check on everyone.
	for _, v := range results {
		repositories <- Repository{
			Name:    v.PathWithNamespace,
			URL:     "https://gitlab.com/" + v.PathWithNamespace + "/raw/" + v.DefaultBranch + "/" + os.Getenv("CRAWLED_FILENAME"),
			Source:  "gitlab.com",
			Headers: headers,
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
		return host.URL, nil
	}

	// Return next url
	parsedLink := parseHeaderLinkGitlab(respHeaders.Get("Link"))

	return parsedLink, nil
}

// parseHeaderLink parse the Gitlab Header Link to nect link of repositories.
// original Link: <https://gitlab.com/api/v4/projects?archived=false&membership=false&order_by=created_at&owned=false&page=2&per_page=20&simple=false&sort=desc&starred=false&statistics=false&with_custom_attributes=false&with_issues_enabled=false&with_merge_requests_enabled=false>; rel="next", <https://gitlab.com/api/v4/projects?archived=false&membership=false&order_by=created_at&owned=false&page=1&per_page=20&simple=false&sort=desc&starred=false&statistics=false&with_custom_attributes=false&with_issues_enabled=false&with_merge_requests_enabled=false>; rel="first", <https://gitlab.com/api/v4/projects?archived=false&membership=false&order_by=created_at&owned=false&page=21994&per_page=20&simple=false&sort=desc&starred=false&statistics=false&with_custom_attributes=false&with_issues_enabled=false&with_merge_requests_enabled=false>; rel="last"
// parsedLink: https://gitlab.com/api/v4/projects?archived=false&membership=false&order_by=created_at&owned=false&page=2&per_page=20&simple=false&sort=desc&starred=false&statistics=false&with_custom_attributes=false&with_issues_enabled=false&with_merge_requests_enabled=false
func parseHeaderLinkGitlab(link string) string {
	parsedLinks := linkheader.Parse(link)

	for _, link := range parsedLinks {
		if link.Rel == "next" {
			return link.URL
		}
	}
	return link
}
