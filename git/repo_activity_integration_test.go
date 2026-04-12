package git

import (
	"os"
	"path/filepath"
	"testing"

	gogit "github.com/go-git/go-git/v5"
	"github.com/italia/publiccode-crawler/v4/common"
	"github.com/italia/publiccode-crawler/v4/git/vitality"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeCache marshals a Cache and writes it to the expected path under dataDir.
func writeCache(t *testing.T, dataDir, host, owner, name string, cache vitality.Cache) {
	t.Helper()

	dir := filepath.Join(dataDir, "repos", host, owner, name)
	require.NoError(t, os.MkdirAll(dir, 0o755))

	data, err := vitality.Marshal(cache)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(filepath.Join(dir, "vitality.json"), data, 0o644))
}

// setupSelfRepoCache runs buildVitalityCache on the publiccode-crawler repo
// itself and writes the resulting cache to a temp data dir.
func setupSelfRepoCache(t *testing.T) (dataDir string, repo common.Repository) {
	t.Helper()

	host := "github.com"
	owner := "italia"
	name := "publiccode-crawler"

	r, err := gogit.PlainOpen(".")
	require.NoError(t, err)

	cache, err := buildVitalityCache(r, nil)
	require.NoError(t, err)

	dataDir = t.TempDir()
	writeCache(t, dataDir, host, owner, name, cache)

	repo = common.Repository{Name: owner + "/" + name}
	repo.URL.Host = host

	return dataDir, repo
}

// TestCalculateRepoActivity exercises buildVitalityCache and CalculateRepoActivity
// end to end against the publiccode-crawler repo itself, with a fixed now so
// the score stays deterministic.
func TestCalculateRepoActivity(t *testing.T) {
	t.Chdir("..")

	dataDir, repo := setupSelfRepoCache(t)
	viper.Set("DATADIR", dataDir)

	total, index, err := CalculateRepoActivity(repo, 365, fixedNow)
	require.NoError(t, err)

	assert.Equal(t, float64(93), total)
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
func TestCalculateRepoActivityLongevityBoundary(t *testing.T) {
	t.Chdir("..")

	host := "github.com"
	owner := "boundary"
	name := "repo"

	dataDir := t.TempDir()

	cache := vitality.Cache{
		LastUpdated:      fixedNow,
		OldestCommitDate: fixedNow.AddDate(0, 0, -365),
		Entries: []vitality.DayEntry{
			{
				Date:    fixedNow.AddDate(0, 0, -365),
				Commits: 1,
				Merges:  0,
				Authors: []string{"mario@example.com"},
			},
		},
	}

	writeCache(t, dataDir, host, owner, name, cache)

	repo := common.Repository{Name: owner + "/" + name}
	repo.URL.Host = host
	viper.Set("DATADIR", dataDir)

	total, index, err := CalculateRepoActivity(repo, 1, fixedNow)
	require.NoError(t, err)

	assert.Equal(t, float64(56), total)
	assert.Len(t, index, 1)
}
