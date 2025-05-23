package git

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/italia/publiccode-crawler/v4/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// RangesData contains the data loaded from vitality-ranges.yml.
type RangesData []Ranges

// Ranges are the ranges for a specific parameter (userCommunity, codeActivity, releaseHistory, longevity).
type Ranges struct {
	Name   string
	Ranges []Range
}

// Range is a range between will be assigned Points value.
type Range struct {
	Min    float64
	Max    float64
	Points float64
}

// CalculateRepoActivity return the repository activity index and the vitality slice calculated on the git clone.
// It follows the document https://lg-acquisizione-e-riuso-software-per-la-pa.readthedocs.io/
// In reference to section: 2.5.2. Fase 2.2: Valutazione soluzioni riusabili per la PA.
func CalculateRepoActivity(repository common.Repository, days int) (float64, map[int]float64, error) {
	if repository.Name == "" {
		return 0, nil, errors.New("cannot  calculate repository activity without name")
	}

	vendor, repo := common.SplitFullName(repository.Name)

	path := filepath.Join(viper.GetString("DATADIR"), "repos", repository.URL.Host, vendor, repo, "gitClone")

	// MkdirAll will create all the folder path, if not exists.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return 0, nil, err
	}
	// Repository activity score.
	var (
		userCommunity  float64
		codeActivity   float64
		releaseHistory float64
		longevity      float64

		repoActivity float64
	)

	// Open and load the git repo path.
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Error(err)

		return 0, nil, err
	}

	// Extract all the commits.
	commits, err := extractAllCommits(r)
	if err != nil {
		log.Error(err)
	}

	// List commits before a number of days: commitsLastDays[from days before today][]commits
	commitsLastDays := extractCommitsLastDays(days, commits)

	// List commits in a day: commitsPerDay[day][]commits
	commitsPerDay := extractCommitsPerDay(days, commits)

	// Extract all tags.
	tags, err := extractAllTagsCommit(r)
	if err != nil {
		log.Error(err)
	}

	// List tags in a day: tagsPerDay[day][]commits
	tagsPerDays := extractTagsPerDay(days, tags)

	// For every day (and before) calculate the Vitality index.
	vitalityIndex := map[int]float64{}

	// Longevity is the repository age.
	longevity, err = calculateLongevityIndex(r)
	if err != nil {
		log.Warn(err)
	}

	for i := range days {
		userCommunity = ranges("userCommunity", userCommunityLastDays(commitsLastDays[i]))

		codeActivity = ranges("codeActivity", activityLastDays(commitsPerDay[i]))
		releaseHistory = ranges("releaseHistory", releaseHistoryLastDays(tagsPerDays[i]))

		repoActivity = userCommunity + codeActivity + releaseHistory + ranges("longevity", longevity)
		if repoActivity > 100 {
			repoActivity = 100
		}

		vitalityIndex[i] = repoActivity
	}

	vitalityIndexTotal := meanActivity(vitalityIndex)
	if vitalityIndexTotal > 100 {
		vitalityIndexTotal = float64(100)
	}

	return float64(int(vitalityIndexTotal)), vitalityIndex, nil
}

// userCommunityLastDays returns the number of unique commits authors.
func userCommunityLastDays(commits []*object.Commit) float64 {
	// Prepare single author map.
	totalAuthors := map[string]int{}
	// Iterates over the commits and extract infos.
	for _, c := range commits {
		totalAuthors[c.Author.Email]++
	}

	return float64(len(totalAuthors))
}

// activityLastDays: # commits and # merges.
func activityLastDays(commits []*object.Commit) float64 {
	numberCommits := float64(len(commits))
	numberMerges := 0

	for _, c := range commits {
		if c.NumParents() > 1 {
			numberMerges++
		}
	}

	return numberCommits + float64(numberMerges)
}

// releaseHistoryLastDays: number of releases.
func releaseHistoryLastDays(tags []*object.Commit) float64 {
	return float64(len(tags))
}

// Extract all commits referred to released Tags.
func extractAllTagsCommit(r *git.Repository) ([]*object.Commit, error) {
	var allTags []*object.Commit

	tagrefs, err := r.Tags()
	if err != nil {
		return nil, err
	}

	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		if !t.Hash().IsZero() {
			tagObject, _ := r.CommitObject(t.Hash())
			allTags = append(allTags, tagObject)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return allTags, nil
}

// extractAllCommits returns a slice of all the commits from the passed repository.
func extractAllCommits(r *git.Repository) ([]*object.Commit, error) {
	var commits []*object.Commit

	ref, err := r.Head()
	if err != nil {
		log.Error(err)

		return nil, err
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Error(err)

		return nil, err
	}

	err = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)

		return nil
	})
	if err != nil {
		log.Error(err)
	}

	return commits, nil
}

func calculateLongevityIndex(r *git.Repository) (float64, error) {
	ref, err := r.Head()
	if err != nil {
		log.Error(err)

		return 0, err
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Error(err)

		return 0, err
	}

	creationDate := extractOldestCommitDate(cIter)

	age := time.Since(creationDate).Hours() / 24

	// Git was invented in 2005. If some repo starts before, remove.
	then := time.Date(2005, time.January, 1, 1, 0, 0, 0, time.UTC)

	duration := time.Since(then).Hours()
	if age > duration/24 {
		return -1, errors.New("first commit is too old. Must be after the creation of git (2005)")
	}

	return age, err
}

// extractOldestCommitDate returns the oldest commit date.
func extractOldestCommitDate(cIter object.CommitIter) time.Time {
	// Iterates over the commits and extract infos.
	result := time.Now()
	_ = cIter.ForEach(func(c *object.Commit) error {
		if c.Author.When.Before(result) {
			result = c.Author.When
		}

		return nil
	})

	return result
}

func ranges(name string, value float64) float64 {
	data, err := os.ReadFile("vitality-ranges.yml")
	if err != nil {
		log.Error(err)
	}
	// Prepare the data structure for load.
	t := RangesData{}
	// Populate the yaml.
	err = yaml.Unmarshal(data, &t)
	if err != nil {
		log.Errorf("error: %v", err)
	}

	for _, v := range t {
		// Select the right ranges table.
		if v.Name == name {
			for _, r := range v.Ranges {
				if value >= r.Min && value < r.Max {
					return r.Points
				}
			}
		}
	}

	return 0
}

// extractCommitsLastDays returns a map of last days commits.
func extractCommitsLastDays(days int, commits []*object.Commit) map[int][]*object.Commit {
	commitsLastDays := map[int][]*object.Commit{}
	// Populate the slice of commits in every day.
	for i := range days {
		lastDays := time.Now().AddDate(0, 0, -i)
		// Append all the commits created before lastDays date.
		for _, c := range commits {
			if c.Author.When.Before(lastDays) {
				commitsLastDays[i] = append(commitsLastDays[i], c)
			}
		}
	}

	return commitsLastDays
}

// extractCommitsPerDay returns a map of number of commits per day, in the last [days].
func extractCommitsPerDay(days int, commits []*object.Commit) map[int][]*object.Commit {
	commitsPerDay := map[int][]*object.Commit{}
	// Populate the slice of commits in every day.
	for i := range days {
		lastDays := time.Now().AddDate(0, 0, -i)
		// Append all the commits created before lastDays date.
		for _, c := range commits {
			if c.Author.When.Day() == lastDays.Day() &&
				c.Author.When.Month() == lastDays.Month() && c.Author.When.Year() ==
				lastDays.Year() {
				commitsPerDay[i] = append(commitsPerDay[i], c)
			}
		}
	}

	return commitsPerDay
}

// extractTagsPerDay returns a map of #[days] commits where a tag is created.
func extractTagsPerDay(days int, tags []*object.Commit) map[int][]*object.Commit {
	tagsPerDays := map[int][]*object.Commit{}

	for i := range days {
		lastDays := time.Now().AddDate(0, 0, -i)
		// Append all the commits created before lastDays date.
		for _, t := range tags {
			if t != nil {
				if t.Author.When.Day() == lastDays.Day() &&
					t.Author.When.Month() == lastDays.Month() &&
					t.Author.When.Year() == lastDays.Year() {
					tagsPerDays[i] = append(tagsPerDays[i], t)
				}
			}
		}
	}

	return tagsPerDays
}

// meanActivity return the mean of all the points.
func meanActivity(points map[int]float64) float64 {
	var total float64
	for _, point := range points {
		total += point
	}

	return total / float64(len(points))
}
