package git

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/italia/publiccode-crawler/v4/common"
	"github.com/italia/publiccode-crawler/v4/git/vitality"
	"github.com/italia/publiccode-crawler/v4/metrics"
	"github.com/spf13/viper"
)

func vitilityCachePath(hostname, vendor, repo string) string {
	return filepath.Join(viper.GetString("DATADIR"), "repos", hostname, vendor, repo, "vitality.json")
}

func CloneRepository(hostname, name, gitURL, index string) error {
	if name == "" {
		return errors.New("cannot save a file without name")
	}

	if gitURL == "" {
		return errors.New("cannot clone a repository without git URL")
	}

	vendor, repo := common.SplitFullName(name)
	cachePath := vitilityCachePath(hostname, vendor, repo)

	existing := loadExistingCache(cachePath)

	tmpDir, err := os.MkdirTemp("", "vitality-*")
	if err != nil {
		return fmt.Errorf("cannot create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	args := []string{"clone", "--filter=blob:none", "--bare"}
	if existing != nil && !existing.LastUpdated.IsZero() {
		args = append(args, "--shallow-since="+existing.LastUpdated.Format("2006-01-02"))
	}

	args = append(args, gitURL, tmpDir)

	cmd := exec.CommandContext(context.Background(), "git", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("cannot clone %s: %w: %s", gitURL, err, out)
	}

	r, err := gogit.PlainOpen(tmpDir)
	if err != nil {
		return fmt.Errorf("cannot open clone of %s: %w", gitURL, err)
	}

	metrics.GetCounter("repository_cloned", index).Inc()

	cache, err := buildVitalityCache(r, existing)
	if err != nil {
		return fmt.Errorf("cannot build vitality cache: %w", err)
	}

	data, err := vitality.Marshal(cache)
	if err != nil {
		return fmt.Errorf("cannot marshal vitality cache: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(cachePath), 0o755); err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0o600)
}

func loadExistingCache(path string) *vitality.Cache {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	c, err := vitality.Unmarshal(data)
	if err != nil {
		return nil
	}

	return &c
}

type dayKey struct {
	y int
	m time.Month
	d int
}

func buildVitalityCache(r *gogit.Repository, existing *vitality.Cache) (vitality.Cache, error) {
	var cache vitality.Cache
	cache.LastUpdated = time.Now().UTC()

	existingDays := map[dayKey]vitality.DayEntry{}
	existingTags := map[dayKey]vitality.TagEntry{}

	if existing != nil {
		cache.OldestCommitDate = existing.OldestCommitDate

		for _, e := range existing.Entries {
			existingDays[dayKey{e.Date.Year(), e.Date.Month(), e.Date.Day()}] = e
		}

		for _, t := range existing.Tags {
			existingTags[dayKey{t.Date.Year(), t.Date.Month(), t.Date.Day()}] = t
		}
	}

	commits, err := extractAllCommits(r)
	if err != nil {
		return cache, err
	}

	newDays := map[dayKey]*vitality.DayEntry{}

	for _, c := range commits {
		t := c.Author.When.UTC()

		if cache.OldestCommitDate.IsZero() || t.Before(cache.OldestCommitDate) {
			cache.OldestCommitDate = t
		}

		key := dayKey{t.Year(), t.Month(), t.Day()}
		if newDays[key] == nil {
			newDays[key] = &vitality.DayEntry{
				Date: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC),
			}
		}

		e := newDays[key]
		e.Commits++

		if c.NumParents() > 1 {
			e.Merges++
		}

		if !slices.Contains(e.Authors, c.Author.Email) {
			e.Authors = append(e.Authors, c.Author.Email)
		}
	}

	merged := maps.Clone(existingDays)
	for k, e := range newDays {
		merged[k] = *e
	}

	days := make([]vitality.DayEntry, 0, len(merged))
	for _, d := range merged {
		days = append(days, d)
	}

	sort.Slice(days, func(i, j int) bool {
		return days[i].Date.Before(days[j].Date)
	})

	cache.Entries = days

	tags, err := extractAllTagsCommit(r)
	if err != nil {
		return cache, err
	}

	newTags := map[dayKey]vitality.TagEntry{}

	for _, t := range tags {
		if t == nil {
			continue
		}

		d := t.Author.When.UTC()
		key := dayKey{d.Year(), d.Month(), d.Day()}
		entry := newTags[key]
		entry.Date = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
		entry.Count++
		newTags[key] = entry
	}

	mergedTags := maps.Clone(existingTags)
	maps.Copy(mergedTags, newTags)

	tagDays := make([]vitality.TagEntry, 0, len(mergedTags))
	for _, t := range mergedTags {
		tagDays = append(tagDays, t)
	}

	sort.Slice(tagDays, func(i, j int) bool {
		return tagDays[i].Date.Before(tagDays[j].Date)
	})

	cache.Tags = tagDays

	return cache, nil
}

func extractAllCommits(r *gogit.Repository) ([]*object.Commit, error) {
	var commits []*object.Commit

	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	cIter, err := r.Log(&gogit.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}

	err = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)

		return nil
	})

	return commits, err
}

func extractAllTagsCommit(r *gogit.Repository) ([]*object.Commit, error) {
	var allTags []*object.Commit

	tagrefs, err := r.Tags()
	if err != nil {
		return nil, err
	}

	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		if !t.Hash().IsZero() {
			tagObject, _ := r.CommitObject(t.Hash())
			if tagObject != nil {
				allTags = append(allTags, tagObject)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return allTags, nil
}
