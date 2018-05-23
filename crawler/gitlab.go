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
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
)

// Gitlab represent a complete result for the Gitlab API respose from the groups call.
type GitlabGroups struct {
	ID                   int    `json:"id"`
	WebURL               string `json:"web_url"`
	Name                 string `json:"name"`
	Path                 string `json:"path"`
	Description          string `json:"description"`
	Visibility           string `json:"visibility"`
	LfsEnabled           bool   `json:"lfs_enabled"`
	AvatarURL            string `json:"avatar_url"`
	RequestAccessEnabled bool   `json:"request_access_enabled"`
	FullName             string `json:"full_name"`
	FullPath             string `json:"full_path"`
	ParentID             string `json:"parent_id"`
	Projects             []struct {
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
		AvatarURL         string        `json:"avatar_url"`
		StarCount         int           `json:"star_count"`
		ForksCount        int           `json:"forks_count"`
		LastActivityAt    time.Time     `json:"last_activity_at"`
		Links             struct {
			Self          string `json:"self"`
			Issues        string `json:"issues"`
			MergeRequests string `json:"merge_requests"`
			RepoBranches  string `json:"repo_branches"`
			Labels        string `json:"labels"`
			Events        string `json:"events"`
			Members       string `json:"members"`
		} `json:"_links"`
		Archived                       bool   `json:"archived"`
		Visibility                     string `json:"visibility"`
		ResolveOutdatedDiffDiscussions bool   `json:"resolve_outdated_diff_discussions"`
		ContainerRegistryEnabled       bool   `json:"container_registry_enabled"`
		IssuesEnabled                  bool   `json:"issues_enabled"`
		MergeRequestsEnabled           bool   `json:"merge_requests_enabled"`
		WikiEnabled                    bool   `json:"wiki_enabled"`
		JobsEnabled                    bool   `json:"jobs_enabled"`
		SnippetsEnabled                bool   `json:"snippets_enabled"`
		SharedRunnersEnabled           bool   `json:"shared_runners_enabled"`
		LfsEnabled                     bool   `json:"lfs_enabled"`
		CreatorID                      int    `json:"creator_id"`
		Namespace                      struct {
			ID       int         `json:"id"`
			Name     string      `json:"name"`
			Path     string      `json:"path"`
			Kind     string      `json:"kind"`
			FullPath string      `json:"full_path"`
			ParentID interface{} `json:"parent_id"`
		} `json:"namespace"`
		ImportStatus                              string        `json:"import_status"`
		OpenIssuesCount                           int           `json:"open_issues_count,omitempty"`
		PublicJobs                                bool          `json:"public_jobs"`
		CiConfigPath                              string        `json:"ci_config_path"`
		SharedWithGroups                          []interface{} `json:"shared_with_groups"`
		OnlyAllowMergeIfPipelineSucceeds          bool          `json:"only_allow_merge_if_pipeline_succeeds"`
		RequestAccessEnabled                      bool          `json:"request_access_enabled"`
		OnlyAllowMergeIfAllDiscussionsAreResolved bool          `json:"only_allow_merge_if_all_discussions_are_resolved"`
		PrintingMergeRequestLinkEnabled           bool          `json:"printing_merge_request_link_enabled"`
		MergeMethod                               string        `json:"merge_method"`
		ApprovalsBeforeMerge                      int           `json:"approvals_before_merge"`
	} `json:"projects"`
	SharedProjects []struct {
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
		AvatarURL         string        `json:"avatar_url"`
		StarCount         int           `json:"star_count"`
		ForksCount        int           `json:"forks_count"`
		LastActivityAt    time.Time     `json:"last_activity_at"`
		Links             struct {
			Self          string `json:"self"`
			Issues        string `json:"issues"`
			MergeRequests string `json:"merge_requests"`
			RepoBranches  string `json:"repo_branches"`
			Labels        string `json:"labels"`
			Events        string `json:"events"`
			Members       string `json:"members"`
		} `json:"_links"`
		Archived                       bool   `json:"archived"`
		Visibility                     string `json:"visibility"`
		ResolveOutdatedDiffDiscussions bool   `json:"resolve_outdated_diff_discussions"`
		ContainerRegistryEnabled       bool   `json:"container_registry_enabled"`
		IssuesEnabled                  bool   `json:"issues_enabled"`
		MergeRequestsEnabled           bool   `json:"merge_requests_enabled"`
		WikiEnabled                    bool   `json:"wiki_enabled"`
		JobsEnabled                    bool   `json:"jobs_enabled"`
		SnippetsEnabled                bool   `json:"snippets_enabled"`
		SharedRunnersEnabled           bool   `json:"shared_runners_enabled"`
		LfsEnabled                     bool   `json:"lfs_enabled"`
		CreatorID                      int    `json:"creator_id"`
		Namespace                      struct {
			ID       int         `json:"id"`
			Name     string      `json:"name"`
			Path     string      `json:"path"`
			Kind     string      `json:"kind"`
			FullPath string      `json:"full_path"`
			ParentID interface{} `json:"parent_id"`
		} `json:"namespace"`
		ForkedFromProject struct {
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
		} `json:"forked_from_project,omitempty"`
		ImportStatus     string      `json:"import_status"`
		OpenIssuesCount  int         `json:"open_issues_count,omitempty"`
		PublicJobs       bool        `json:"public_jobs"`
		CiConfigPath     interface{} `json:"ci_config_path"`
		SharedWithGroups []struct {
			GroupID          int    `json:"group_id"`
			GroupName        string `json:"group_name"`
			GroupAccessLevel int    `json:"group_access_level"`
		} `json:"shared_with_groups"`
		OnlyAllowMergeIfPipelineSucceeds          bool   `json:"only_allow_merge_if_pipeline_succeeds"`
		RequestAccessEnabled                      bool   `json:"request_access_enabled"`
		OnlyAllowMergeIfAllDiscussionsAreResolved bool   `json:"only_allow_merge_if_all_discussions_are_resolved"`
		PrintingMergeRequestLinkEnabled           bool   `json:"printing_merge_request_link_enabled"`
		MergeMethod                               string `json:"merge_method"`
		ApprovalsBeforeMerge                      int    `json:"approvals_before_merge"`
		Owner                                     struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Username  string `json:"username"`
			State     string `json:"state"`
			AvatarURL string `json:"avatar_url"`
			WebURL    string `json:"web_url"`
		} `json:"owner,omitempty"`
	} `json:"shared_projects"`
	LdapCn                    interface{} `json:"ldap_cn"`
	LdapAccess                interface{} `json:"ldap_access"`
	SharedRunnersMinutesLimit int         `json:"shared_runners_minutes_limit"`
}

// GitlabRepo represent a complete for the Gitlab API respose from a single repository.
type GitlabRepo struct {
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

// RegisterGitlabAPI register the crawler function for Gitlab API.
func RegisterGitlabAPI() Handler {
	return func(domain Domain, link string, repositories chan Repository, wg *sync.WaitGroup) (string, error) {
		log.Debugf("RegisterGitlabAPI: %s ")

		// Set BasicAuth header
		headers := make(map[string]string)
		if domain.BasicAuth != nil {
			rand.Seed(time.Now().Unix())
			n := rand.Int() % len(domain.BasicAuth)
			headers["Authorization"] = "Basic " + domain.BasicAuth[n]
		}

		// Get List of repositories
		resp, err := httpclient.GetURL(link, headers)
		if err != nil {
			return link, err
		}
		if resp.Status.Code != http.StatusOK {
			log.Warnf("Request returned: %s", string(resp.Body))
			return link, errors.New("request returned an incorrect http.Status: " + resp.Status.Text)
		}

		// Fill response as list of values (repositories data).
		var results GitlabGroups
		err = json.Unmarshal(resp.Body, &results)
		if err != nil {
			return link, err
		}

		// Add repositories to the channel that will perform the check on every project.
		for _, v := range results.Projects {
			log.Debugf("Gitlab Projects %s", v.PathWithNamespace)
			// Join file raw URL.
			u, err := url.Parse(domain.RawBaseUrl)
			if err != nil {
				return link, err
			}
			u.Path = path.Join(u.Path, v.PathWithNamespace, "raw", v.DefaultBranch, viper.GetString("CRAWLED_FILENAME"))

			if v.DefaultBranch != "" {
				repositories <- Repository{
					Name:       v.PathWithNamespace,
					FileRawURL: u.String(),
					Domain:     domain,
					Headers:    headers,
				}
			}
		}
		// Add repositories to the channel that will perform the check on every sharedd project.
		for _, v := range results.SharedProjects {
			log.Debugf("Gitlab SharedProjects %s", v.PathWithNamespace)
			// Join file raw URL.
			u, err := url.Parse(domain.RawBaseUrl)
			if err != nil {
				return link, err
			}
			u.Path = path.Join(u.Path, v.PathWithNamespace, "raw", v.DefaultBranch, viper.GetString("CRAWLED_FILENAME"))

			if v.DefaultBranch != "" {
				repositories <- Repository{
					Name:       v.PathWithNamespace,
					FileRawURL: u.String(),
					Domain:     domain,
					Headers:    headers,
				}
			}
		}
		// if last page for this organization, the Link is empty.
		if len(resp.Headers.Get("Link")) == 0 {
			return "", nil
		}

		// Return next url
		parsedLink := httpclient.NextHeaderLink(resp.Headers.Get("Link"))
		if parsedLink == "" {
			return "", nil
		}

		return parsedLink, nil
	}
}

// RegisterSingleGitlabAPI register the crawler function for single Bitbucket API.
func RegisterSingleGitlabAPI() SingleHandler {
	return func(domain Domain, link string, repositories chan Repository) error {
		// Set BasicAuth header
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

		var fullURL string
		// Starting URL. Generate using go templates.
		fullURL = domain.ApiRepoURL
		data := struct{ Name string }{Name: url.QueryEscape(fullName)}
		// Create a new template and parse the Url into it.
		t := template.Must(template.New("url").Parse(fullURL))
		buf := new(bytes.Buffer)
		// Execute the template: add "data" data in "url".
		t.Execute(buf, data)
		fullURL = buf.String()

		// Get single Repo
		resp, err := httpclient.GetURL(fullURL, headers)
		if err != nil {
			return err
		}
		if resp.Status.Code != http.StatusOK {
			log.Warnf("Request returned: %s", string(resp.Body))
			return errors.New("request returned an incorrect http.Status: " + resp.Status.Text)
		}

		// Fill response as list of values (repositories data).
		var result GitlabRepo
		err = json.Unmarshal(resp.Body, &result)
		if err != nil {
			return err
		}

		// Join file raw URL.
		u, err = url.Parse(domain.RawBaseUrl)
		if err != nil {
			return err
		}
		u.Path = path.Join(u.Path, result.PathWithNamespace, "raw", result.DefaultBranch, viper.GetString("CRAWLED_FILENAME"))

		// If the repository was never used, the Mainbranch is empty ("")
		if result.DefaultBranch != "" {
			repositories <- Repository{
				Name:       result.PathWithNamespace,
				FileRawURL: u.String(),
				Domain:     domain,
				Headers:    headers,
			}
		} else {
			return errors.New("repository is: empty")
		}

		return nil
	}
}
