package scanner_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/italia/publiccode-crawler/v4/common"
	"github.com/italia/publiccode-crawler/v4/scanner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func giteaRepoJSON(name, fullName, defaultBranch string, private, archived, empty bool) map[string]any {
	return map[string]any{
		"name":           name,
		"full_name":      fullName,
		"private":        private,
		"archived":       archived,
		"empty":          empty,
		"default_branch": defaultBranch,
		"html_url":       "http://placeholder/" + fullName,
		"clone_url":      "http://placeholder/" + fullName + ".git",
	}
}

// newGiteaTestServer returns a test server that handles Gitea API requests.
func newGiteaTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func giteaPublisher() common.Publisher {
	return common.Publisher{ID: "test", Name: "Test"}
}

func TestGiteaScanner_ScanRepo_success(t *testing.T) {
	repo := giteaRepoJSON("myrepo", "myorg/myrepo", "main", false, false, false)

	ts := newGiteaTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/repos/myorg/myrepo", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(repo)
	})
	defer ts.Close()

	repoURL, err := url.Parse(ts.URL + "/myorg/myrepo")
	require.NoError(t, err)

	repositories := make(chan common.Repository, 1)

	sc := scanner.NewGiteaScanner()
	err = sc.ScanRepo(*repoURL, giteaPublisher(), repositories)

	require.NoError(t, err)
	require.Len(t, repositories, 1)

	got := <-repositories
	assert.Equal(t, "myorg/myrepo", got.Name)
	assert.Equal(t, "main", got.GitBranch)
	assert.Equal(t, "http://placeholder/myorg/myrepo/raw/branch/main/publiccode.yml", got.FileRawURL)
}

func TestGiteaScanner_ScanRepo_dotGitSuffix(t *testing.T) {
	repo := giteaRepoJSON("myrepo", "myorg/myrepo", "main", false, false, false)

	ts := newGiteaTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/repos/myorg/myrepo", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(repo)
	})
	defer ts.Close()

	repoURL, err := url.Parse(ts.URL + "/myorg/myrepo.git")
	require.NoError(t, err)

	repositories := make(chan common.Repository, 1)

	sc := scanner.NewGiteaScanner()
	err = sc.ScanRepo(*repoURL, giteaPublisher(), repositories)

	require.NoError(t, err)
	assert.Len(t, repositories, 1)
}

func TestGiteaScanner_ScanRepo_private(t *testing.T) {
	repo := giteaRepoJSON("myrepo", "myorg/myrepo", "main", true, false, false)

	ts := newGiteaTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(repo)
	})
	defer ts.Close()

	repoURL, err := url.Parse(ts.URL + "/myorg/myrepo")
	require.NoError(t, err)

	repositories := make(chan common.Repository, 1)

	sc := scanner.NewGiteaScanner()
	err = sc.ScanRepo(*repoURL, giteaPublisher(), repositories)

	require.NoError(t, err)
	assert.Empty(t, repositories)
}

func TestGiteaScanner_ScanRepo_archived(t *testing.T) {
	repo := giteaRepoJSON("myrepo", "myorg/myrepo", "main", false, true, false)

	ts := newGiteaTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(repo)
	})
	defer ts.Close()

	repoURL, err := url.Parse(ts.URL + "/myorg/myrepo")
	require.NoError(t, err)

	repositories := make(chan common.Repository, 1)

	sc := scanner.NewGiteaScanner()
	err = sc.ScanRepo(*repoURL, giteaPublisher(), repositories)

	require.NoError(t, err)
	assert.Empty(t, repositories)
}

func TestGiteaScanner_ScanRepo_empty(t *testing.T) {
	repo := giteaRepoJSON("myrepo", "myorg/myrepo", "", false, false, true)

	ts := newGiteaTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(repo)
	})
	defer ts.Close()

	repoURL, err := url.Parse(ts.URL + "/myorg/myrepo")
	require.NoError(t, err)

	repositories := make(chan common.Repository, 1)

	sc := scanner.NewGiteaScanner()
	err = sc.ScanRepo(*repoURL, giteaPublisher(), repositories)

	require.NoError(t, err)
	assert.Empty(t, repositories)
}

func TestGiteaScanner_ScanRepo_apiError(t *testing.T) {
	ts := newGiteaTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer ts.Close()

	repoURL, err := url.Parse(ts.URL + "/myorg/myrepo")
	require.NoError(t, err)

	repositories := make(chan common.Repository, 1)

	sc := scanner.NewGiteaScanner()
	err = sc.ScanRepo(*repoURL, giteaPublisher(), repositories)

	require.Error(t, err)
	assert.Empty(t, repositories)
}

func TestGiteaScanner_ScanRepo_invalidURL(t *testing.T) {
	repositories := make(chan common.Repository, 1)

	repoURL, _ := url.Parse("http://example.com/onlyone")

	sc := scanner.NewGiteaScanner()
	err := sc.ScanRepo(*repoURL, giteaPublisher(), repositories)

	require.Error(t, err)
}

func TestGiteaScanner_ScanGroupOfRepos_org(t *testing.T) {
	repos := []map[string]any{
		giteaRepoJSON("repo1", "myorg/repo1", "main", false, false, false),
		giteaRepoJSON("repo2", "myorg/repo2", "develop", false, false, false),
	}

	ts := newGiteaTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/orgs/myorg/repos":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(repos)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer ts.Close()

	groupURL, err := url.Parse(ts.URL + "/myorg")
	require.NoError(t, err)

	repositories := make(chan common.Repository, 10)

	sc := scanner.NewGiteaScanner()
	err = sc.ScanGroupOfRepos(*groupURL, giteaPublisher(), repositories)

	require.NoError(t, err)
	assert.Len(t, repositories, 2)
}

func TestGiteaScanner_ScanGroupOfRepos_fallbackToUser(t *testing.T) {
	repos := []map[string]any{
		giteaRepoJSON("repo1", "myuser/repo1", "main", false, false, false),
	}

	ts := newGiteaTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/orgs/myuser/repos":
			// Org not found — fall back to user endpoint.
			w.WriteHeader(http.StatusNotFound)
		case "/api/v1/users/myuser/repos":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(repos)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer ts.Close()

	groupURL, err := url.Parse(ts.URL + "/myuser")
	require.NoError(t, err)

	repositories := make(chan common.Repository, 10)

	sc := scanner.NewGiteaScanner()
	err = sc.ScanGroupOfRepos(*groupURL, giteaPublisher(), repositories)

	require.NoError(t, err)
	assert.Len(t, repositories, 1)
}

func TestGiteaScanner_ScanGroupOfRepos_instance(t *testing.T) {
	repos := []map[string]any{
		giteaRepoJSON("repo1", "org/repo1", "main", false, false, false),
	}

	ts := newGiteaTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/repos/search", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"data": repos})
	})
	defer ts.Close()

	// Root URL = scan whole instance.
	instanceURL, err := url.Parse(ts.URL)
	require.NoError(t, err)

	repositories := make(chan common.Repository, 10)

	sc := scanner.NewGiteaScanner()
	err = sc.ScanGroupOfRepos(*instanceURL, giteaPublisher(), repositories)

	require.NoError(t, err)
	assert.Len(t, repositories, 1)
}

func TestGiteaScanner_ScanGroupOfRepos_skipsPrivateAndArchived(t *testing.T) {
	repos := []map[string]any{
		giteaRepoJSON("pub", "org/pub", "main", false, false, false),
		giteaRepoJSON("priv", "org/priv", "main", true, false, false),
		giteaRepoJSON("arch", "org/arch", "main", false, true, false),
	}

	ts := newGiteaTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(repos)
	})
	defer ts.Close()

	groupURL, err := url.Parse(ts.URL + "/org")
	require.NoError(t, err)

	repositories := make(chan common.Repository, 10)

	sc := scanner.NewGiteaScanner()
	err = sc.ScanGroupOfRepos(*groupURL, giteaPublisher(), repositories)

	require.NoError(t, err)
	assert.Len(t, repositories, 1)

	got := <-repositories
	assert.Equal(t, "org/pub", got.Name)
}

func TestGiteaScanner_ScanGroupOfRepos_pagination(t *testing.T) {
	page1 := make([]map[string]any, 50)
	for i := range page1 {
		name := fmt.Sprintf("repo%d", i)
		page1[i] = giteaRepoJSON(name, "org/"+name, "main", false, false, false)
	}

	page2 := []map[string]any{
		giteaRepoJSON("last", "org/last", "main", false, false, false),
	}

	ts := newGiteaTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Query().Get("page") == "2" {
			_ = json.NewEncoder(w).Encode(page2)
		} else {
			_ = json.NewEncoder(w).Encode(page1)
		}
	})
	defer ts.Close()

	groupURL, err := url.Parse(ts.URL + "/org")
	require.NoError(t, err)

	repositories := make(chan common.Repository, 100)

	sc := scanner.NewGiteaScanner()
	err = sc.ScanGroupOfRepos(*groupURL, giteaPublisher(), repositories)

	require.NoError(t, err)
	assert.Len(t, repositories, 51)
}
