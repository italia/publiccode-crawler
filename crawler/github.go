package crawler

import (
	"encoding/json"
	"errors"
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

// GithubOrgs is the complete result from the Github API respose for /orgs/<Name>/repos.
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

// GithubFiles is a list of files in repository
type GithubFiles []struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Sha         string `json:"sha"`
	Size        int    `json:"size"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	GitURL      string `json:"git_url"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"`
	Links       struct {
		Self string `json:"self"`
		Git  string `json:"git"`
		HTML string `json:"html"`
	} `json:"_links"`
}

// RegisterGithubAPI register the crawler function for Github API.
// It get the list of repositories on "link" url.
// If a next page is available return its url.
// Otherwise returns an empty ("") string.
func RegisterGithubAPI() OrganizationHandler {
	return func(domain Domain, link string, repositories chan Repository, wg *sync.WaitGroup) (string, error) {
		// Set BasicAuth header
		headers := make(map[string]string)
		if domain.BasicAuth != nil {
			n, err := generateRandomInt(len(domain.BasicAuth))
			if err != nil {
				return link, err
			}
			headers["Authorization"] = domain.BasicAuth[n]
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
			// Marshal all the repository metadata.
			metadata, err := json.Marshal(v)
			if err != nil {
				log.Errorf("github metadata: %v", err)
			}
			contents := strings.Replace(v.ContentsURL, "{+path}", "", -1)
			// Get List of files.
			resp, err := httpclient.GetURL(contents, headers)
			if err != nil {
				return link, err
			}
			if resp.Status.Code != http.StatusOK {
				log.Infof("Request returned an invalid status code: %s", string(resp.Body))
			}
			// Fill response as list of values (repositories data).
			var files GithubFiles
			err = json.Unmarshal(resp.Body, &files)
			if err != nil {
				log.Infof("Repository is empty: %s", string(resp.Body))
			}

			// Search a file with a valid name and a downloadURL.
			for _, f := range files {
				if f.Name == viper.GetString("CRAWLED_FILENAME") && f.DownloadURL != "" {
					// Add repository to channel.
					repositories <- Repository{
						Name:       v.FullName,
						FileRawURL: f.DownloadURL,
						Domain:     domain,
						Headers:    headers,
						Metadata:   metadata,
					}
				}
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

// RegisterSingleGithubAPI register the crawler function for single repository Github API.
// Return nil if the repository was successfully added to repositories channel.
// Otherwise return the generated error.
func RegisterSingleGithubAPI() SingleRepoHandler {
	return func(domain Domain, link string, repositories chan Repository) error {
		// Set BasicAuth header.
		headers := make(map[string]string)
		if domain.BasicAuth != nil {
			n, err := generateRandomInt(len(domain.BasicAuth))
			if err != nil {
				return err
			}
			headers["Authorization"] = domain.BasicAuth[n]
		}

		u, err := url.Parse(link)
		if err != nil {
			return err
		}

		u.Path = path.Join("repos", u.Path)
		u.Path = strings.Trim(u.Path, "/")
		u.Host = "api." + u.Host

		// Get List of repositories.
		resp, err := httpclient.GetURL(u.String(), headers)
		if err != nil {
			return err
		}
		if resp.Status.Code != http.StatusOK {
			log.Warnf("Request returned: %s", string(resp.Body))
			return errors.New("request returned an incorrect http.Status: " + resp.Status.Text)
		}

		var v GithubRepo
		err = json.Unmarshal(resp.Body, &v)
		if err != nil {
			return err
		}

		// Marshal all the repository metadata.
		metadata, err := json.Marshal(v)
		if err != nil {
			log.Errorf("github metadata: %v", err)
		}
		contents := strings.Replace(v.ContentsURL, "{+path}", "", -1)
		// Get List of files.
		resp, err = httpclient.GetURL(contents, headers)
		if err != nil {
			return err
		}
		if resp.Status.Code != http.StatusOK {
			log.Infof("Request returned an invalid status code: %s", string(resp.Body))
		}
		// Fill response as list of values (repositories data).
		var files GithubFiles
		err = json.Unmarshal(resp.Body, &files)
		if err != nil {
			log.Infof("Repository is empty: %s", string(resp.Body))
		}

		// Search a file with a valid name and a downloadURL.
		for _, f := range files {
			if f.Name == viper.GetString("CRAWLED_FILENAME") && f.DownloadURL != "" {
				// Add repository to channel.
				repositories <- Repository{
					Name:       v.FullName,
					FileRawURL: f.DownloadURL,
					Domain:     domain,
					Headers:    headers,
					Metadata:   metadata,
				}
			} else {
				return errors.New("Repository does not contain " + viper.GetString("CRAWLED_FILENAME"))
			}
		}

		return nil
	}
}

// GenerateGithubAPIURL returns the api url of given Gitlab organization link.
// IN: https://github.com/italia
// OUT:https://api.github.com/orgs/italia/repos
func GenerateGithubAPIURL() GeneratorAPIURL {
	return func(in string) (string, error) {
		u, err := url.Parse(in)
		if err != nil {
			return in, err
		}
		u.Path = path.Join("orgs", u.Path, "repos")
		u.Path = strings.Trim(u.Path, "/")
		u.Host = "api." + u.Host

		return u.String(), nil
	}
}

// IsGithub returns "true" if the url can use Github API.
func IsGithub(link string) bool {
	if len(link) == 0 {
		log.Errorf("IsGithub: empty link %s.", link)
		return false
	}

	u, err := url.Parse(link)
	if err != nil {
		log.Errorf("IsGithub: impossible to parse %s.", link)
		return false
	}
	u.Path = "rate_limit"
	u.Host = "api." + u.Host

	resp, err := httpclient.GetURL(u.String(), nil)
	if err != nil {
		log.Debugf("can %s use Github API? No.", link)
		return false
	}
	if resp.Status.Code != http.StatusOK {
		log.Debugf("can %s use Github API? No.", link)
		return false
	}

	log.Debugf("can %s use Github API? Yes.", link)
	return true
}
