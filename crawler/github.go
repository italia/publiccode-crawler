package crawler

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/italia/developers-italia-backend/httpclient"
	log "github.com/sirupsen/logrus"
)

// Github is a Crawler for the Github hosting.
type Github struct {
	URL       string
	RateLimit struct {
		ReqH int `yaml:"req/h"`
		ReqM int `yaml:"req/m"`
	} `yaml:"rate-limit"`
	BasicAuth string `yaml:"basic-auth"`
}

type githubResponse []struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"owner"`
	Private          bool   `json:"private"`
	HTMLURL          string `json:"html_url"`
	Description      string `json:"description"`
	Fork             bool   `json:"fork"`
	URL              string `json:"url"`
	ForksURL         string `json:"forks_url"`
	KeysURL          string `json:"keys_url"`
	CollaboratorsURL string `json:"collaborators_url"`
	TeamsURL         string `json:"teams_url"`
	HooksURL         string `json:"hooks_url"`
	IssueEventsURL   string `json:"issue_events_url"`
	EventsURL        string `json:"events_url"`
	AssigneesURL     string `json:"assignees_url"`
	BranchesURL      string `json:"branches_url"`
	TagsURL          string `json:"tags_url"`
	BlobsURL         string `json:"blobs_url"`
	GitTagsURL       string `json:"git_tags_url"`
	GitRefsURL       string `json:"git_refs_url"`
	TreesURL         string `json:"trees_url"`
	StatusesURL      string `json:"statuses_url"`
	LanguagesURL     string `json:"languages_url"`
	StargazersURL    string `json:"stargazers_url"`
	ContributorsURL  string `json:"contributors_url"`
	SubscribersURL   string `json:"subscribers_url"`
	SubscriptionURL  string `json:"subscription_url"`
	CommitsURL       string `json:"commits_url"`
	GitCommitsURL    string `json:"git_commits_url"`
	CommentsURL      string `json:"comments_url"`
	IssueCommentURL  string `json:"issue_comment_url"`
	ContentsURL      string `json:"contents_url"`
	CompareURL       string `json:"compare_url"`
	MergesURL        string `json:"merges_url"`
	ArchiveURL       string `json:"archive_url"`
	DownloadsURL     string `json:"downloads_url"`
	IssuesURL        string `json:"issues_url"`
	PullsURL         string `json:"pulls_url"`
	MilestonesURL    string `json:"milestones_url"`
	NotificationsURL string `json:"notifications_url"`
	LabelsURL        string `json:"labels_url"`
	ReleasesURL      string `json:"releases_url"`
	DeploymentsURL   string `json:"deployments_url"`
}

// GetRepositories retrieves the list of all repository from an hosting.
// Return the URL from where it should restart (Next or actual if fails) and error.
func (host Github) GetRepositories(url string, repositories chan Repository) (string, error) {
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
		// Github abuse rate limit management: https://developer.github.com/v3/guides/best-practices-for-integrators/#dealing-with-abuse-rate-limits
		if len(respHeaders.Get("Retry-After")) > 0 {
			log.Infof("Waiting: %s seconds. (The value of Header Retry-After)", respHeaders.Get("Retry-After"))
			secondsAfterRetry, _ := strconv.Atoi(respHeaders.Get("Retry-After"))
			time.Sleep(time.Second * time.Duration(secondsAfterRetry))
			return url, nil
		} else {
			return url, errors.New("requets returned an incorrect http.Status: " + status.Status)
		}
	}

	// Fill response as list of values (repositories data).
	var results githubResponse
	err = json.Unmarshal(body, &results)
	if err != nil {
		return url, err
	}

	// Add repositories to the channel that will perform the check on everyone.
	for _, v := range results {
		repositories <- Repository{
			Name:    v.FullName,
			URL:     v.URL + "/contents/" + os.Getenv("CRAWLED_FILENAME"),
			Source:  url,
			Headers: headers,
		}
	}

	if len(respHeaders.Get("Link")) == 0 {
		// If I want to restart when it ends:
		// sourceURL = "https://api.bitbucket.org/2.0/repositories?pagelen=100&after=2008-08-13"
		// and comment the line "close(repositories)"
		log.Info("Github repositories status: end reached.")
		close(repositories)
		return url, nil
	}

	// Return next url
	parsedLink := parseHeaderLink(respHeaders.Get("Link"))

	return parsedLink, nil
}

// parseHeaderLink parse the Github Header Link to nect link of repositories.
// original Link: <https://api.github.com/repositories?since=1592>; rel="next", <https://api.github.com/repositories{?since}>; rel="first"
// parsedLink: https://api.github.com/repositories?since=1592
func parseHeaderLink(link string) string {
	parsedLink := strings.Split(link, ";")[0]
	parsedLink = parsedLink[1:][:len(parsedLink)-2]

	return parsedLink
}
