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

// GitlabGroups is the complete result from the Gitlab API respose.
type GitlabGroups struct {
	ID                        int                   `json:"id"`
	WebURL                    string                `json:"web_url"`
	Name                      string                `json:"name"`
	Path                      string                `json:"path"`
	Description               string                `json:"description"`
	Visibility                string                `json:"visibility"`
	LfsEnabled                bool                  `json:"lfs_enabled"`
	AvatarURL                 string                `json:"avatar_url"`
	RequestAccessEnabled      bool                  `json:"request_access_enabled"`
	FullName                  string                `json:"full_name"`
	FullPath                  string                `json:"full_path"`
	ParentID                  string                `json:"parent_id"`
	Projects                  []GitlabProject       `json:"projects"`
	SharedProjects            []GitlabSharedProject `json:"shared_projects"`
	LdapCn                    interface{}           `json:"ldap_cn"`
	LdapAccess                interface{}           `json:"ldap_access"`
	SharedRunnersMinutesLimit int                   `json:"shared_runners_minutes_limit"`
}

// GitlabRepo is a complete result from the Gitlab API respose for a single repository.
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

// GitlabProject is a software project hosted on Gitlab.
type GitlabProject struct {
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
}

// GitlabSharedProject is a software project hosted on Gitlab, owned by a group and shared with someone.
type GitlabSharedProject struct {
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
	ForkedFromProject GitlabRepo  `json:"forked_from_project,omitempty"`
	ImportStatus      string      `json:"import_status"`
	OpenIssuesCount   int         `json:"open_issues_count,omitempty"`
	PublicJobs        bool        `json:"public_jobs"`
	CiConfigPath      interface{} `json:"ci_config_path"`
	SharedWithGroups  []struct {
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
}

// RegisterGitlabAPI register the crawler function for Gitlab API.
func RegisterGitlabAPI() OrganizationHandler {
	return func(domain Domain, link string, repositories chan Repository, pa PA, wg *sync.WaitGroup) (string, error) {
		log.Debugf("RegisterGitlabAPI: %s ", link)

		// Set BasicAuth header.
		headers := make(map[string]string)
		if domain.BasicAuth != nil {
			n, err := generateRandomInt(len(domain.BasicAuth))
			if err != nil {
				return link, err
			}
			headers["Authorization"] = domain.BasicAuth[n]
		}

		// Parse url.
		u, err := url.Parse(link)
		if err != nil {
			return link, err
		}
		// Set domain host to new host.
		domain.Host = u.Hostname()

		// Get List of repositories.
		resp, err := httpclient.GetURL(link, headers)
		if err != nil {
			return link, err
		}
		if resp.Status.Code != http.StatusOK {
			log.Warnf("Request returned: %s", string(resp.Body))
			return "", errors.New("request returned an incorrect http.Status: " + resp.Status.Text)
		}

		// Fill response as list of values (repositories data).
		var results GitlabGroups
		err = json.Unmarshal(resp.Body, &results)
		if err != nil {
			return link, err
		}

		// Add repositories to the channel that will perform the check on every project.
		err = addGitlabProjectsToRepositories(results.Projects, domain, pa, headers, repositories)
		if err != nil {
			return link, err
		}
		// Add repositories to the channel that will perform the check on every sharedd project.
		err = addGitlabSharedProjectsToRepositories(results.SharedProjects, domain, pa, headers, repositories)
		if err != nil {
			return link, err
		}

		// if last page for this organization, the Link is empty.
		if len(resp.Headers.Get("Link")) == 0 {
			return "", nil
		}

		// Return next url.
		nextLink := httpclient.HeaderLink(resp.Headers.Get("Link"), "next")
		if nextLink == "" {
			return "", nil
		}

		return nextLink, nil
	}
}

// RegisterSingleGitlabAPI register the crawler function for single Bitbucket API.
func RegisterSingleGitlabAPI() SingleRepoHandler {
	return func(domain Domain, link string, repositories chan Repository) error {
		// Set BasicAuth header
		headers := make(map[string]string)
		if domain.BasicAuth != nil {
			n, err := generateRandomInt(len(domain.BasicAuth))
			if err != nil {
				return err
			}
			headers["Authorization"] = "Basic " + domain.BasicAuth[n]
		}

		// Parse url.
		u, err := url.Parse(link)
		if err != nil {
			log.Error(err)
			return err
		}

		// Set domain host to new host.
		domain.Host = u.Hostname()

		// Dirty concatenation. With the normal URL String() the escaped characters are escaped two times.
		repoString := strings.Trim(u.Path, "/")
		fullURL := "https://" + u.Hostname() + "/api/v4/projects/" + url.QueryEscape(repoString)

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

		// Join file raw URL string.
		_, err = generateGitlabRawURL(result.WebURL, result.DefaultBranch)
		if err != nil {
			return err
		}

		// Marshal all the repository metadata.
		metadata, err := json.Marshal(result)
		if err != nil {
			log.Errorf("gitlab metadata: %v", err)
			return err
		}

		// If the repository was never used, the Mainbranch is empty ("")
		if result.DefaultBranch != "" {
			repositories <- Repository{
				Name:        result.PathWithNamespace,
				FileRawURL:  u.String(),
				GitCloneURL: result.HTTPURLToRepo,
				GitBranch:   result.DefaultBranch,
				Hostname:    u.Hostname(),
				Domain:      domain,
				Headers:     headers,
				Metadata:    metadata,
			}
		} else {
			return errors.New("repository is empty." + result.WebURL)
		}

		return nil
	}
}

// generateGitlabRawURL returns the file Gitlab specific file raw url.
func generateGitlabRawURL(baseURL, defaultBranch string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, "raw", defaultBranch, viper.GetString("CRAWLED_FILENAME"))

	return u.String(), err
}

// addGitlabProjectsToRepositories adds the projects from api response to repository channel.
func addGitlabProjectsToRepositories(projects []GitlabProject, domain Domain, pa PA, headers map[string]string, repositories chan Repository) error {
	for _, v := range projects {
		// Join file raw URL string.
		rawURL, err := generateGitlabRawURL(v.WebURL, v.DefaultBranch)
		if err != nil {
			return err
		}

		// Marshal all the repository metadata.
		metadata, err := json.Marshal(v)
		if err != nil {
			log.Errorf("gitlab metadata: %v", err)
			return err
		}

		if v.DefaultBranch != "" {
			repositories <- Repository{
				Name:        v.PathWithNamespace,
				Hostname:    domain.Host,
				FileRawURL:  rawURL,
				GitCloneURL: v.HTTPURLToRepo,
				GitBranch:   v.DefaultBranch,
				Domain:      domain,
				Pa:          pa,
				Headers:     headers,
				Metadata:    metadata,
			}
		}
	}

	return nil
}

// addGitlabSharedProjectsToRepositories adds the shared projects from api response to repository channel.
func addGitlabSharedProjectsToRepositories(projects []GitlabSharedProject, domain Domain, pa PA, headers map[string]string, repositories chan Repository) error {
	for _, v := range projects {
		// Join file raw URL string.
		rawURL, err := generateGitlabRawURL(v.WebURL, v.DefaultBranch)
		if err != nil {
			return err
		}

		// Marshal all the repository metadata.
		metadata, err := json.Marshal(v)
		if err != nil {
			log.Errorf("gitlab metadata: %v", err)
			return err
		}

		if v.DefaultBranch != "" {
			repositories <- Repository{
				Name:        v.PathWithNamespace,
				Hostname:    domain.Host,
				FileRawURL:  rawURL,
				GitCloneURL: v.HTTPURLToRepo,
				GitBranch:   v.DefaultBranch,
				Domain:      domain,
				Pa:          pa,
				Headers:     headers,
				Metadata:    metadata,
			}
		}
	}

	return nil
}

// GenerateGitlabAPIURL returns the api url of given Gitlab organization link.
// IN: https://gitlab.org/blockninja
// OUT:https://gitlab.com/api/v4/groups/blockninja
func GenerateGitlabAPIURL() GeneratorAPIURL {
	return func(in string) (string, error) {
		u, err := url.Parse(in)
		if err != nil {
			return in, err
		}
		u.Path = path.Join("api/v4/groups", u.Path)

		return u.String(), nil
	}
}

// IsGitlab returns "true" if the url can use Gitlab API.
func IsGitlab(link string) bool {
	if len(link) == 0 {
		log.Errorf("IsGitlab: empty link %s.", link)
		return false
	}

	u, err := url.Parse(link)
	if err != nil {
		log.Errorf("IsGitlab: impossible to parse %s.", link)
		return false
	}

	u.Path = "api/v4/templates/gitlab_ci_ymls"

	resp, err := httpclient.GetURL(u.String(), nil)
	if err != nil {
		log.Debugf("can %s use Gitlab API? No.", link)
		return false
	}
	if resp.Status.Code != http.StatusOK {
		log.Debugf("can %s use Gitlab API? No.", link)
		return false
	}

	log.Debugf("can %s use Gitlab API? Yes.", link)
	return true
}
