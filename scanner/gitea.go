package scanner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/italia/publiccode-crawler/v4/common"
	log "github.com/sirupsen/logrus"
)

var errNotFound = errors.New("not found")

type GiteaScanner struct{}

func NewGiteaScanner() Scanner {
	return GiteaScanner{}
}

type giteaRepo struct {
	Name          string `json:"name"`
	FullName      string `json:"full_name"` //nolint:tagliatelle // Gitea API uses snake_case
	Private       bool   `json:"private"`
	Archived      bool   `json:"archived"`
	DefaultBranch string `json:"default_branch"` //nolint:tagliatelle // Gitea API uses snake_case
	HTMLURL       string `json:"html_url"`       //nolint:tagliatelle // Gitea API uses snake_case
	CloneURL      string `json:"clone_url"`      //nolint:tagliatelle // Gitea API uses snake_case
	Empty         bool   `json:"empty"`
}

type giteaSearchResult struct {
	Data []giteaRepo `json:"data"`
}

// ScanGroupOfRepos scans a Gitea or Forgejo org/user represented by u,
// or all public repos on the instance if u is a root URL.
func (scanner GiteaScanner) ScanGroupOfRepos(
	groupURL url.URL, publisher common.Publisher, repositories chan common.Repository,
) error {
	log.Debugf("GiteaScanner.ScanGroupOfRepos(%s)", groupURL.String())

	owner := strings.Trim(groupURL.Path, "/")

	const limit = 50

	fetchPage := giteaOrgReposPage
	if owner == "" {
		fetchPage = giteaInstanceReposPage
	}

	for page := 1; ; page++ {
		repos, err := fetchPage(&groupURL, owner, limit, page)
		if err != nil {
			return fmt.Errorf("GiteaScanner: %w", err)
		}

		for _, repo := range repos {
			if err := addGiteaRepo(nil, repo, publisher, repositories); err != nil {
				return err
			}
		}

		if len(repos) < limit {
			break
		}
	}

	return nil
}

// ScanRepo scans a single Gitea or Forgejo repository represented by repoURL.
func (scanner GiteaScanner) ScanRepo(
	repoURL url.URL, publisher common.Publisher, repositories chan common.Repository,
) error {
	log.Debugf("GiteaScanner.ScanRepo(%s)", repoURL.String())

	repoPath := strings.TrimSuffix(strings.Trim(repoURL.Path, "/"), ".git")

	parts := strings.SplitN(repoPath, "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("GiteaScanner: invalid repo URL: %s", repoURL.String())
	}

	owner, repoName := parts[0], parts[1]
	apiURL := fmt.Sprintf("%s://%s/api/v1/repos/%s/%s", repoURL.Scheme, repoURL.Host, owner, repoName)

	repo, err := giteaFetchRepo(apiURL)
	if err != nil {
		return fmt.Errorf("GiteaScanner: %w", err)
	}

	return addGiteaRepo(&repoURL, repo, publisher, repositories)
}

func addGiteaRepo(
	originalURL *url.URL, repo giteaRepo,
	publisher common.Publisher, repositories chan common.Repository,
) error {
	if repo.Private || repo.Archived || repo.Empty || repo.DefaultBranch == "" {
		return nil
	}

	canonicalURL, err := url.Parse(repo.CloneURL)
	if err != nil {
		return fmt.Errorf("GiteaScanner: parse clone URL %s: %w", repo.CloneURL, err)
	}

	if originalURL == nil {
		originalURL = canonicalURL
	}

	rawURL := fmt.Sprintf("%s/raw/branch/%s/publiccode.yml",
		strings.TrimSuffix(repo.HTMLURL, ".git"),
		repo.DefaultBranch,
	)

	repositories <- common.Repository{
		Name:         repo.FullName,
		FileRawURL:   rawURL,
		URL:          *originalURL,
		CanonicalURL: *canonicalURL,
		GitBranch:    repo.DefaultBranch,
		Publisher:    publisher,
	}

	return nil
}

func giteaOrgReposPage(base *url.URL, owner string, limit, page int) ([]giteaRepo, error) {
	apiURL := fmt.Sprintf("%s://%s/api/v1/orgs/%s/repos?limit=%d&page=%d",
		base.Scheme, base.Host, url.PathEscape(owner), limit, page)

	repos, err := giteaFetchRepoList(apiURL)
	if !errors.Is(err, errNotFound) {
		return repos, err
	}

	// Fall back to user endpoint if the owner is not an org.
	userURL := fmt.Sprintf("%s://%s/api/v1/users/%s/repos?limit=%d&page=%d",
		base.Scheme, base.Host, url.PathEscape(owner), limit, page)

	return giteaFetchRepoList(userURL)
}

func giteaInstanceReposPage(base *url.URL, _ string, limit, page int) ([]giteaRepo, error) {
	apiURL := fmt.Sprintf("%s://%s/api/v1/repos/search?limit=%d&page=%d",
		base.Scheme, base.Host, limit, page)

	req, err := giteaNewRequest(apiURL)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", apiURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: status %d", apiURL, resp.StatusCode)
	}

	var result giteaSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode search result: %w", err)
	}

	return result.Data, nil
}

func giteaFetchRepoList(apiURL string) ([]giteaRepo, error) {
	req, err := giteaNewRequest(apiURL)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", apiURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: status %d", apiURL, resp.StatusCode)
	}

	var repos []giteaRepo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("decode repo list: %w", err)
	}

	return repos, nil
}

func giteaFetchRepo(apiURL string) (giteaRepo, error) {
	req, err := giteaNewRequest(apiURL)
	if err != nil {
		return giteaRepo{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return giteaRepo{}, fmt.Errorf("GET %s: %w", apiURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return giteaRepo{}, fmt.Errorf("GET %s: status %d", apiURL, resp.StatusCode)
	}

	var repo giteaRepo
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return giteaRepo{}, fmt.Errorf("decode repo: %w", err)
	}

	return repo, nil
}

func giteaNewRequest(apiURL string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("GiteaScanner: new request: %w", err)
	}

	if token := os.Getenv("GITEA_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	return req, nil
}
