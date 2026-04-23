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

type dayData struct {
	date    time.Time
	commits uint32
	merges  uint32
	authors []string
}

func buildVitalityCache(r *gogit.Repository, existing *vitality.Cache) (vitality.Cache, error) {
	var cache vitality.Cache
	cache.LastUpdated = time.Now().UTC()

	authorIndex := map[string]uint16{}
	existingDays := map[dayKey]dayData{}
	existingTags := map[dayKey]uint32{}

	if existing != nil {
		for _, a := range existing.Authors {
			authorIndex[a] = uint16(len(authorIndex)) //nolint:gosec // Authors capped at 65535
		}

		cache.Authors = append(cache.Authors, existing.Authors...)
		cache.OldestCommitDate = existing.OldestCommitDate
		existingDays = decodeExistingDays(existing)
		existingTags = decodeExistingTags(existing)
	}

	commits, err := extractAllCommits(r)
	if err != nil {
		return cache, err
	}

	newDays, err := aggregateCommits(commits, &cache, authorIndex)
	if err != nil {
		return cache, err
	}

	cache.Entries = mergeDayEntries(existingDays, newDays, authorIndex, &cache)

	tags, err := extractAllTagsCommit(r)
	if err != nil {
		return cache, err
	}

	cache.Tags = mergeTagEntries(existingTags, tags, cache.FirstEntryDate)

	return cache, nil
}

func decodeExistingDays(existing *vitality.Cache) map[dayKey]dayData {
	out := map[dayKey]dayData{}
	cur := existing.FirstEntryDate

	for _, e := range existing.Entries {
		cur = cur.AddDate(0, 0, int(e.Delta))
		authors := make([]string, 0, len(e.Authors))

		for _, id := range e.Authors {
			if int(id) < len(existing.Authors) {
				authors = append(authors, existing.Authors[id])
			}
		}

		out[dayKey{cur.Year(), cur.Month(), cur.Day()}] = dayData{
			date: cur, commits: e.Commits, merges: e.Merges, authors: authors,
		}
	}

	return out
}

func decodeExistingTags(existing *vitality.Cache) map[dayKey]uint32 {
	out := map[dayKey]uint32{}
	cur := existing.FirstEntryDate

	for _, t := range existing.Tags {
		cur = cur.AddDate(0, 0, int(t.Delta))
		out[dayKey{cur.Year(), cur.Month(), cur.Day()}] = t.Count
	}

	return out
}

func aggregateCommits(
	commits []*object.Commit, cache *vitality.Cache, authorIndex map[string]uint16,
) (map[dayKey]*dayData, error) {
	newDays := map[dayKey]*dayData{}

	for _, c := range commits {
		t := c.Author.When.UTC()

		if cache.OldestCommitDate.IsZero() || t.Before(cache.OldestCommitDate) {
			cache.OldestCommitDate = t
		}

		if _, ok := authorIndex[c.Author.Email]; !ok {
			if len(authorIndex) >= 65535 {
				return nil, errors.New("vitality cache: more than 65535 unique authors")
			}

			authorIndex[c.Author.Email] = uint16(len(authorIndex)) //nolint:gosec // bounded by check above
			cache.Authors = append(cache.Authors, c.Author.Email)
		}

		key := dayKey{t.Year(), t.Month(), t.Day()}
		if newDays[key] == nil {
			newDays[key] = &dayData{
				date: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC),
			}
		}

		e := newDays[key]
		e.commits++

		if c.NumParents() > 1 {
			e.merges++
		}

		if !slices.Contains(e.authors, c.Author.Email) {
			e.authors = append(e.authors, c.Author.Email)
		}
	}

	return newDays, nil
}

func mergeDayEntries(
	existing map[dayKey]dayData, fresh map[dayKey]*dayData,
	authorIndex map[string]uint16, cache *vitality.Cache,
) []vitality.DayEntry {
	merged := maps.Clone(existing)
	for k, e := range fresh {
		merged[k] = *e
	}

	days := make([]dayData, 0, len(merged))
	for _, d := range merged {
		days = append(days, d)
	}

	sort.Slice(days, func(i, j int) bool {
		return days[i].date.Before(days[j].date)
	})

	if len(days) == 0 {
		return nil
	}

	cache.FirstEntryDate = days[0].date
	prev := days[0].date
	out := make([]vitality.DayEntry, 0, len(days))

	for _, d := range days {
		delta := uint32(d.date.Sub(prev).Hours() / 24)
		authorIDs := make([]uint16, len(d.authors))

		for i, a := range d.authors {
			authorIDs[i] = authorIndex[a]
		}

		out = append(out, vitality.DayEntry{
			Delta:   delta,
			Commits: d.commits,
			Merges:  d.merges,
			Authors: authorIDs,
		})
		prev = d.date
	}

	return out
}

func mergeTagEntries(existing map[dayKey]uint32, tags []*object.Commit, firstEntry time.Time) []vitality.TagEntry {
	tagDates := map[dayKey]time.Time{}
	newTags := map[dayKey]uint32{}

	for _, t := range tags {
		if t == nil {
			continue
		}

		d := t.Author.When.UTC()
		key := dayKey{d.Year(), d.Month(), d.Day()}
		tagDates[key] = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
		newTags[key]++
	}

	merged := maps.Clone(existing)
	maps.Copy(merged, newTags)

	for k := range merged {
		if _, ok := tagDates[k]; !ok {
			tagDates[k] = time.Date(k.y, k.m, k.d, 0, 0, 0, 0, time.UTC)
		}
	}

	type tagDay struct {
		date  time.Time
		count uint32
	}

	tagDays := make([]tagDay, 0, len(merged))
	for k, count := range merged {
		tagDays = append(tagDays, tagDay{date: tagDates[k], count: count})
	}

	sort.Slice(tagDays, func(i, j int) bool {
		return tagDays[i].date.Before(tagDays[j].date)
	})

	if len(tagDays) == 0 {
		return nil
	}

	prev := firstEntry
	out := make([]vitality.TagEntry, 0, len(tagDays))

	for _, t := range tagDays {
		delta := uint32(t.date.Sub(prev).Hours() / 24)
		out = append(out, vitality.TagEntry{
			Delta: delta,
			Count: t.count,
		})
		prev = t.date
	}

	return out
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
