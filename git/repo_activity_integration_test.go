package git

import (
	"os"
	"path/filepath"
	"testing"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/italia/publiccode-crawler/v4/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRepo builds an in-memory git repository with a fixed commit history
// relative to fixedNow so that vitality scores are deterministic.
//
// Commit schedule (relative to fixedNow):
//
//	-200d mario  (oldest, longevity = 200d = 20 pts)
//	-100d luigi
//	 -50d mario
//	 -10d luca
//	  -5d mario
//	  -2d luigi  (merge commit)
//	  -1d luca
func setupTestRepo(t *testing.T) (dataDir string, repo common.Repository) {
	t.Helper()

	host := "github.com"
	owner := "test"
	name := "repo"

	dataDir = t.TempDir()
	repoPath := filepath.Join(dataDir, "repos", host, owner, name, "gitClone")
	require.NoError(t, os.MkdirAll(repoPath, 0o755))

	r, err := gogit.PlainInit(repoPath, false)
	require.NoError(t, err)

	w, err := r.Worktree()
	require.NoError(t, err)

	commits := []struct {
		email    string
		daysBack int
		merge    bool
	}{
		{"mario@example.com", 200, false},
		{"luigi@example.com", 100, false},
		{"mario@example.com", 50, false},
		{"luca@example.com", 10, false},
		{"mario@example.com", 5, false},
		{"luigi@example.com", 2, true},
		{"luca@example.com", 1, false},
	}

	var prevHash [2]plumbing.Hash

	for i, c := range commits {
		f := filepath.Join(repoPath, "file.txt")
		require.NoError(t, os.WriteFile(f, []byte(c.email), 0o644))

		_, err = w.Add("file.txt")
		require.NoError(t, err)

		opts := &gogit.CommitOptions{
			Author: &object.Signature{
				Name:  c.email,
				Email: c.email,
				When:  fixedNow.AddDate(0, 0, -c.daysBack),
			},
		}
		if c.merge && i >= 2 {
			opts.Parents = []plumbing.Hash{prevHash[0], prevHash[1]}
		}

		h, err := w.Commit("commit "+c.email, opts)
		require.NoError(t, err)

		prevHash[1] = prevHash[0]
		prevHash[0] = h
	}

	repo = common.Repository{Name: owner + "/" + name}
	repo.URL.Host = host

	return dataDir, repo
}

// TestCalculateRepoActivity verifies the end-to-end vitality score against a
// precisely known commit history.
//
// Expected score derivation (days=365, fixedNow):
//
//	i=0..99    userCommunity=8 (3 authors) + codeActivity=2 + releaseHistory=20 + longevity=20 = 50
//	i=100..364 userCommunity=4 (<=1 author) + codeActivity=2 + releaseHistory=20 + longevity=20 = 46
//	mean = (100*50 + 265*46) / 365 = 47.09..., truncated to 47
func TestCalculateRepoActivity(t *testing.T) {
	t.Chdir("..")

	dataDir, repo := setupTestRepo(t)
	viper.Set("DATADIR", dataDir)

	total, index, err := CalculateRepoActivity(repo, 365, fixedNow)
	require.NoError(t, err)

	assert.Equal(t, float64(47), total)
	assert.Len(t, index, 365)
}

func TestCalculateRepoActivityMissingCache(t *testing.T) {
	viper.Set("DATADIR", t.TempDir())

	repo := common.Repository{Name: "owner/nonexistent"}
	repo.URL.Host = "github.com"

	_, _, err := CalculateRepoActivity(repo, 365, fixedNow)
	assert.Error(t, err)
}

// TestCalculateRepoActivityLongevityBoundary checks that a repo whose oldest
// commit is exactly 365 days old crosses into the [365, 730) longevity bucket
// (30 pts instead of 20 pts), raising the total score.
//
// Setup: 1 commit at -365d by mario, days=1.
//
//	userCommunity:  1 author  = 4 pts
//	codeActivity:   0 commits = 2 pts
//	releaseHistory: no tags   = 20 pts
//	longevity:      365d, range [365, 730) = 30 pts
//	total = 4+2+20+30 = 56
func TestCalculateRepoActivityLongevityBoundary(t *testing.T) {
	t.Chdir("..")

	host := "github.com"
	owner := "boundary"
	name := "repo"

	dataDir := t.TempDir()
	repoPath := filepath.Join(dataDir, "repos", host, owner, name, "gitClone")
	require.NoError(t, os.MkdirAll(repoPath, 0o755))

	r, err := gogit.PlainInit(repoPath, false)
	require.NoError(t, err)

	w, err := r.Worktree()
	require.NoError(t, err)

	f := filepath.Join(repoPath, "file.txt")
	require.NoError(t, os.WriteFile(f, []byte("x"), 0o644))

	_, err = w.Add("file.txt")
	require.NoError(t, err)

	_, err = w.Commit("init", &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  "mario",
			Email: "mario@example.com",
			When:  fixedNow.AddDate(0, 0, -365),
		},
	})
	require.NoError(t, err)

	repo := common.Repository{Name: owner + "/" + name}
	repo.URL.Host = host
	viper.Set("DATADIR", dataDir)

	total, index, err := CalculateRepoActivity(repo, 1, fixedNow)
	require.NoError(t, err)

	assert.Equal(t, float64(56), total)
	assert.Len(t, index, 1)
}
