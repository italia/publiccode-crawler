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

	repository, err := gogit.PlainOpen(tmpDir)
	if err != nil {
		return fmt.Errorf("cannot open clone of %s: %w", gitURL, err)
	}

	metrics.GetCounter("repository_cloned", index).Inc()

	cache, err := buildVitalityCache(repository, existing)
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

	cache, err := vitality.Unmarshal(data)
	if err != nil {
		return nil
	}

	return &cache
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

func buildVitalityCache(repository *gogit.Repository, existing *vitality.Cache) (vitality.Cache, error) {
	var cache vitality.Cache
	cache.LastUpdated = time.Now().UTC()

	authorIndex := map[string]uint16{}
	existingDays := map[dayKey]dayData{}
	existingTags := map[dayKey]uint32{}

	if existing != nil {
		for _, author := range existing.Authors {
			authorIndex[author] = uint16(len(authorIndex)) //nolint:gosec // Authors capped at 65535
		}

		cache.Authors = append(cache.Authors, existing.Authors...)
		cache.OldestCommitDate = existing.OldestCommitDate
		existingDays = decodeExistingDays(existing)
		existingTags = decodeExistingTags(existing)
	}

	commits, err := extractAllCommits(repository)
	if err != nil {
		return cache, err
	}

	newDays, err := aggregateCommits(commits, &cache, authorIndex)
	if err != nil {
		return cache, err
	}

	cache.Entries = mergeDayEntries(existingDays, newDays, authorIndex, &cache)

	tags, err := extractAllTagsCommit(repository)
	if err != nil {
		return cache, err
	}

	cache.Tags = mergeTagEntries(existingTags, tags, cache.FirstEntryDate)

	return cache, nil
}

func decodeExistingDays(existing *vitality.Cache) map[dayKey]dayData {
	out := map[dayKey]dayData{}
	cur := existing.FirstEntryDate

	for _, entry := range existing.Entries {
		cur = cur.AddDate(0, 0, int(entry.Delta))
		authors := make([]string, 0, len(entry.Authors))

		for _, id := range entry.Authors {
			if int(id) < len(existing.Authors) {
				authors = append(authors, existing.Authors[id])
			}
		}

		out[dayKey{cur.Year(), cur.Month(), cur.Day()}] = dayData{
			date: cur, commits: entry.Commits, merges: entry.Merges, authors: authors,
		}
	}

	return out
}

func decodeExistingTags(existing *vitality.Cache) map[dayKey]uint32 {
	out := map[dayKey]uint32{}
	cur := existing.FirstEntryDate

	for _, tag := range existing.Tags {
		cur = cur.AddDate(0, 0, int(tag.Delta))
		out[dayKey{cur.Year(), cur.Month(), cur.Day()}] = tag.Count
	}

	return out
}

func aggregateCommits(
	commits []*object.Commit, cache *vitality.Cache, authorIndex map[string]uint16,
) (map[dayKey]*dayData, error) {
	newDays := map[dayKey]*dayData{}

	for _, commit := range commits {
		when := commit.Author.When.UTC()

		if cache.OldestCommitDate.IsZero() || when.Before(cache.OldestCommitDate) {
			cache.OldestCommitDate = when
		}

		if _, ok := authorIndex[commit.Author.Email]; !ok {
			if len(authorIndex) >= 65535 {
				return nil, errors.New("vitality cache: more than 65535 unique authors")
			}

			authorIndex[commit.Author.Email] = uint16(len(authorIndex)) //nolint:gosec // bounded by check above
			cache.Authors = append(cache.Authors, commit.Author.Email)
		}

		key := dayKey{when.Year(), when.Month(), when.Day()}
		if newDays[key] == nil {
			newDays[key] = &dayData{
				date: time.Date(when.Year(), when.Month(), when.Day(), 0, 0, 0, 0, time.UTC),
			}
		}

		entry := newDays[key]
		entry.commits++

		if commit.NumParents() > 1 {
			entry.merges++
		}

		if !slices.Contains(entry.authors, commit.Author.Email) {
			entry.authors = append(entry.authors, commit.Author.Email)
		}
	}

	return newDays, nil
}

func mergeDayEntries(
	existing map[dayKey]dayData, fresh map[dayKey]*dayData,
	authorIndex map[string]uint16, cache *vitality.Cache,
) []vitality.DayEntry {
	merged := maps.Clone(existing)
	for k, entry := range fresh {
		merged[k] = *entry
	}

	days := make([]dayData, 0, len(merged))
	for _, day := range merged {
		days = append(days, day)
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

	for _, day := range days {
		delta := uint32(day.date.Sub(prev).Hours() / 24)
		authorIDs := make([]uint16, len(day.authors))

		for idx, author := range day.authors {
			authorIDs[idx] = authorIndex[author]
		}

		out = append(out, vitality.DayEntry{
			Delta:   delta,
			Commits: day.commits,
			Merges:  day.merges,
			Authors: authorIDs,
		})
		prev = day.date
	}

	return out
}

func mergeTagEntries(existing map[dayKey]uint32, tags []*object.Commit, firstEntry time.Time) []vitality.TagEntry {
	tagDates := map[dayKey]time.Time{}
	newTags := map[dayKey]uint32{}

	for _, tag := range tags {
		if tag == nil {
			continue
		}

		when := tag.Author.When.UTC()
		key := dayKey{when.Year(), when.Month(), when.Day()}
		tagDates[key] = time.Date(when.Year(), when.Month(), when.Day(), 0, 0, 0, 0, time.UTC)
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

	for _, tag := range tagDays {
		delta := uint32(tag.date.Sub(prev).Hours() / 24)
		out = append(out, vitality.TagEntry{
			Delta: delta,
			Count: tag.count,
		})
		prev = tag.date
	}

	return out
}

func extractAllCommits(repository *gogit.Repository) ([]*object.Commit, error) {
	var commits []*object.Commit

	ref, err := repository.Head()
	if err != nil {
		return nil, err
	}

	cIter, err := repository.Log(&gogit.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}

	err = cIter.ForEach(func(commit *object.Commit) error {
		commits = append(commits, commit)

		return nil
	})

	return commits, err
}

func extractAllTagsCommit(repository *gogit.Repository) ([]*object.Commit, error) {
	var allTags []*object.Commit

	tagrefs, err := repository.Tags()
	if err != nil {
		return nil, err
	}

	err = tagrefs.ForEach(func(ref *plumbing.Reference) error {
		if !ref.Hash().IsZero() {
			tagObject, _ := repository.CommitObject(ref.Hash())
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
