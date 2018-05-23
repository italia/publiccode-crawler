package crawler

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"sync"

	"github.com/italia/developers-italia-backend/httpclient"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// GithubOrgs is the complete result from the Github API respose for /orgs/<orgNamee>/repos.
type GithubOrgs []struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	FullName         string    `json:"full_name"`
	Owner            Owner     `json:"owner"`
	Private          bool      `json:"private"`
	HTMLURL          string    `json:"html_url"`
	Description      string    `json:"description"`
	Fork             bool      `json:"fork"`
	URL              string    `json:"url"`
	ForksURL         string    `json:"forks_url"`
	KeysURL          string    `json:"keys_url"`
	CollaboratorsURL string    `json:"collaborators_url"`
	TeamsURL         string    `json:"teams_url"`
	HooksURL         string    `json:"hooks_url"`
	IssueEventsURL   string    `json:"issue_events_url"`
	EventsURL        string    `json:"events_url"`
	AssigneesURL     string    `json:"assignees_url"`
	BranchesURL      string    `json:"branches_url"`
	TagsURL          string    `json:"tags_url"`
	BlobsURL         string    `json:"blobs_url"`
	GitTagsURL       string    `json:"git_tags_url"`
	GitRefsURL       string    `json:"git_refs_url"`
	TreesURL         string    `json:"trees_url"`
	StatusesURL      string    `json:"statuses_url"`
	LanguagesURL     string    `json:"languages_url"`
	StargazersURL    string    `json:"stargazers_url"`
	ContributorsURL  string    `json:"contributors_url"`
	SubscribersURL   string    `json:"subscribers_url"`
	SubscriptionURL  string    `json:"subscription_url"`
	CommitsURL       string    `json:"commits_url"`
	GitCommitsURL    string    `json:"git_commits_url"`
	CommentsURL      string    `json:"comments_url"`
	IssueCommentURL  string    `json:"issue_comment_url"`
	ContentsURL      string    `json:"contents_url"`
	CompareURL       string    `json:"compare_url"`
	MergesURL        string    `json:"merges_url"`
	ArchiveURL       string    `json:"archive_url"`
	DownloadsURL     string    `json:"downloads_url"`
	IssuesURL        string    `json:"issues_url"`
	PullsURL         string    `json:"pulls_url"`
	MilestonesURL    string    `json:"milestones_url"`
	NotificationsURL string    `json:"notifications_url"`
	LabelsURL        string    `json:"labels_url"`
	ReleasesURL      string    `json:"releases_url"`
	DeploymentsURL   string    `json:"deployments_url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	PushedAt         time.Time `json:"pushed_at"`
	GitURL           string    `json:"git_url"`
	SSHURL           string    `json:"ssh_url"`
	CloneURL         string    `json:"clone_url"`
	SvnURL           string    `json:"svn_url"`
	Homepage         string    `json:"homepage"`
	Size             int       `json:"size"`
	StargazersCount  int       `json:"stargazers_count"`
	WatchersCount    int       `json:"watchers_count"`
	Language         string    `json:"language"`
	HasIssues        bool      `json:"has_issues"`
	HasProjects      bool      `json:"has_projects"`
	HasDownloads     bool      `json:"has_downloads"`
	HasWiki          bool      `json:"has_wiki"`
	HasPages         bool      `json:"has_pages"`
	ForksCount       int       `json:"forks_count"`
	MirrorURL        string    `json:"mirror_url"`
	Archived         bool      `json:"archived"`
	OpenIssuesCount  int       `json:"open_issues_count"`
	License          struct {
		Key    string `json:"key"`
		Name   string `json:"name"`
		SpdxID string `json:"spdx_id"`
		URL    string `json:"url"`
	} `json:"license"`
	Forks         int    `json:"forks"`
	OpenIssues    int    `json:"open_issues"`
	Watchers      int    `json:"watchers"`
	DefaultBranch string `json:"default_branch"`
	Permissions   struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permissions"`
}

// GithubRepo is a complete result from the Github API respose for a single repository.
type GithubRepo struct {
	ID               int         `json:"id"`
	Name             string      `json:"name"`
	FullName         string      `json:"full_name"`
	Owner            Owner       `json:"owner"`
	Private          bool        `json:"private"`
	HTMLURL          string      `json:"html_url"`
	Description      string      `json:"description"`
	Fork             bool        `json:"fork"`
	URL              string      `json:"url"`
	ForksURL         string      `json:"forks_url"`
	KeysURL          string      `json:"keys_url"`
	CollaboratorsURL string      `json:"collaborators_url"`
	TeamsURL         string      `json:"teams_url"`
	HooksURL         string      `json:"hooks_url"`
	IssueEventsURL   string      `json:"issue_events_url"`
	EventsURL        string      `json:"events_url"`
	AssigneesURL     string      `json:"assignees_url"`
	BranchesURL      string      `json:"branches_url"`
	TagsURL          string      `json:"tags_url"`
	BlobsURL         string      `json:"blobs_url"`
	GitTagsURL       string      `json:"git_tags_url"`
	GitRefsURL       string      `json:"git_refs_url"`
	TreesURL         string      `json:"trees_url"`
	StatusesURL      string      `json:"statuses_url"`
	LanguagesURL     string      `json:"languages_url"`
	StargazersURL    string      `json:"stargazers_url"`
	ContributorsURL  string      `json:"contributors_url"`
	SubscribersURL   string      `json:"subscribers_url"`
	SubscriptionURL  string      `json:"subscription_url"`
	CommitsURL       string      `json:"commits_url"`
	GitCommitsURL    string      `json:"git_commits_url"`
	CommentsURL      string      `json:"comments_url"`
	IssueCommentURL  string      `json:"issue_comment_url"`
	ContentsURL      string      `json:"contents_url"`
	CompareURL       string      `json:"compare_url"`
	MergesURL        string      `json:"merges_url"`
	ArchiveURL       string      `json:"archive_url"`
	DownloadsURL     string      `json:"downloads_url"`
	IssuesURL        string      `json:"issues_url"`
	PullsURL         string      `json:"pulls_url"`
	MilestonesURL    string      `json:"milestones_url"`
	NotificationsURL string      `json:"notifications_url"`
	LabelsURL        string      `json:"labels_url"`
	ReleasesURL      string      `json:"releases_url"`
	DeploymentsURL   string      `json:"deployments_url"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	PushedAt         time.Time   `json:"pushed_at"`
	GitURL           string      `json:"git_url"`
	SSHURL           string      `json:"ssh_url"`
	CloneURL         string      `json:"clone_url"`
	SvnURL           string      `json:"svn_url"`
	Homepage         string      `json:"homepage"`
	Size             int         `json:"size"`
	StargazersCount  int         `json:"stargazers_count"`
	WatchersCount    int         `json:"watchers_count"`
	Language         string      `json:"language"`
	HasIssues        bool        `json:"has_issues"`
	HasProjects      bool        `json:"has_projects"`
	HasDownloads     bool        `json:"has_downloads"`
	HasWiki          bool        `json:"has_wiki"`
	HasPages         bool        `json:"has_pages"`
	ForksCount       int         `json:"forks_count"`
	MirrorURL        interface{} `json:"mirror_url"`
	Archived         bool        `json:"archived"`
	OpenIssuesCount  int         `json:"open_issues_count"`
	License          interface{} `json:"license"`
	Forks            int         `json:"forks"`
	OpenIssues       int         `json:"open_issues"`
	Watchers         int         `json:"watchers"`
	DefaultBranch    string      `json:"default_branch"`
	NetworkCount     int         `json:"network_count"`
	SubscribersCount int         `json:"subscribers_count"`
}

// Owner of the repository.
type Owner struct {
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
}

// RegisterGithubAPI register the crawler function for Github API.
// It get the list of repositories on "link" url.
// If a next page is available return its url.
// Otherwise returns an empty ("") string.
func RegisterGithubAPI() Handler {
	return func(domain Domain, link string, repositories chan Repository, wg *sync.WaitGroup) (string, error) {
		// Set BasicAuth header
		headers := make(map[string]string)
		if domain.BasicAuth != nil {
			rand.Seed(time.Now().Unix())
			n := rand.Int() % len(domain.BasicAuth)
			headers["Authorization"] = "Basic " + domain.BasicAuth[n]
		}

		// Get List of repositories.
		resp, err := httpclient.GetURL(link, headers)
		if err != nil {
			return link, err
		}
		if resp.Status.Code != http.StatusOK {
			log.Warnf("Request returned: %s", string(resp.Body))
			return link, errors.New("request returned an incorrect http.Status: " + resp.Status.Text)
		}

		// Fill response as list of values (repositories data).
		var results GithubOrgs
		err = json.Unmarshal(resp.Body, &results)
		if err != nil {
			return link, err
		}

		// Add repositories to the channel that will perform the check on everyone.
		for _, v := range results {
			// Join file raw URL.
			u, err := url.Parse(domain.RawBaseUrl)
			if err != nil {
				return link, err
			}
			u.Path = path.Join(u.Path, v.FullName, v.DefaultBranch, viper.GetString("CRAWLED_FILENAME"))
			// Add repository to channel.
			repositories <- Repository{
				Name:       v.FullName,
				FileRawURL: u.String(),
				Domain:     domain,
				Headers:    headers,
			}
		}

		// Return next url.
		nextLink := httpclient.NextHeaderLink(resp.Headers.Get("Link"))

		// if last page for this organization, the nextLink is empty.
		if nextLink == "" {
			return "", nil
		}

		return nextLink, nil
	}
}

// RegisterSingleBitbucketAPI register the crawler function for single repository Github API.
// Return nil if the repository was successfully added to repositories channel.
// Otherwise return the generated error.
func RegisterSingleGithubAPI() SingleHandler {
	return func(domain Domain, link string, repositories chan Repository) error {
		// Set BasicAuth header.
		headers := make(map[string]string)
		if domain.BasicAuth != nil {
			rand.Seed(time.Now().Unix())
			n := rand.Int() % len(domain.BasicAuth)
			headers["Authorization"] = "Basic " + domain.BasicAuth[n]
		}

		u, err := url.Parse(link)
		if err != nil {
			log.Error(err)
		}

		// Clear the url.
		fullName := strings.Trim(u.Path, "/")

		// Generate fullURL using go templates. It will replace {{.Name}} with fullName.
		fullURL := domain.ApiRepoURL
		data := struct{ Name string }{Name: url.QueryEscape(fullName)}
		// Create a new template and parse the Url into it.
		t := template.Must(template.New("url").Parse(fullURL))
		buf := new(bytes.Buffer)
		// Execute the template: add "data" data in "url".
		t.Execute(buf, data)
		fullURL = buf.String()

		// Fill response as list of values (repositories data).
		result, _, err := getGithubRepoInfos(fullURL, headers)
		if err != nil {
			return err
		}

		// Join file raw URL.
		u, err = url.Parse(domain.RawBaseUrl)
		if err != nil {
			return err
		}
		u.Path = path.Join(u.Path, result.FullName, result.DefaultBranch, viper.GetString("CRAWLED_FILENAME"))

		// If the repository was never used, the Mainbranch is empty ("").
		if result.DefaultBranch != "" {
			repositories <- Repository{
				Name:       result.FullName,
				FileRawURL: u.String(),
				Domain:     domain,
				Headers:    headers,
			}
		} else {
			return errors.New("repository is empty")
		}

		return nil
	}
}

// getGithubRepoInfos retrieve single repository info GithubRepo.
func getGithubRepoInfos(URL string, headers map[string]string) (GithubRepo, string, error) {
	var results GithubRepo

	// Get List of repositories.
	resp, err := httpclient.GetURL(URL, headers)
	if err != nil {
		return results, URL, err
	}
	if resp.Status.Code != http.StatusOK {
		log.Warnf("Request for single Github repository returned: %s", string(resp.Body))
		return results, URL, errors.New("request returned an incorrect http.Status: " + resp.Status.Text)
	}

	// Fill response as list of values (repositories data).
	err = json.Unmarshal(resp.Body, &results)
	if err != nil {
		return results, URL, err
	}

	return results, URL, nil

}
