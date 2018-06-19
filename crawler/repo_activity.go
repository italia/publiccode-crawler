package crawler

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// CalculateRepoActivity return the repository activity index calculated on the git clone.
// It follows the document https://lg-acquisizione-e-riuso-software-per-la-pa.readthedocs.io/
// In reference to section: fase-2-2-valutazione-soluzioni-riusabili-per-la-pa
func CalculateRepoActivity(domain Domain, hostname string, name string) (float64, error) {
	if domain.Host == "" {
		return 0, errors.New("cannot calculate repository activity without domain host")
	}
	if name == "" {
		return 0, errors.New("cannot  calculate repository activity without name")
	}

	vendor, repo := splitFullName(name)

	path := filepath.Join("./data", hostname, vendor, repo, "gitClone")

	// MkdirAll will create all the folder path, if not exists.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return 0, err
	}

	// Repository activity score.
	var (
		userCommunity  float64 // max 25
		codeActivity   float64 // max 30
		releaseHistory float64 // max 15
		longevity      float64 // max 25

		repoActivity float64 // max 95
	)

	// Open and load the git repo path.
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Error(err)
	}

	// Total authors. (userCommunity index)
	// 0 to 25 = + #authors
	// > 25 = +25 userCommunity
	// Retrieves the commit history.
	userCommunity, err = calculateUserCommunityIndex(r)
	if err != nil {
		log.Error(err)
	}

	// Oldest (first) commit. (logevity index)
	// 0 to 6 months = +5
	// 6 to 12 months = +10
	// 12 to 24 months = +15
	// 24 to 36 months = +20
	// > 36 months = +25 longevity points
	longevity, err = calculateLongevityIndex(r)
	if err != nil {
		log.Error(err)
	}

	// Tags. (releaseHistory)
	// 0 to 1 tag = +5
	// 1 to 5 tags = +10
	// > 5 tags = +15 releaseHistory points
	releaseHistory, err = calculateReleaseHistoryIndex(r)
	if err != nil {
		log.Error(err)
	}

	// Commits last year. (codeActivity)
	// 0 to 100 commits = +5
	// 1 to 200 commits = +10
	// > 200 commmits = +15 codeActivity points
	// And merges last year. (codeActivity)
	// 0 to 10 merges = +5
	// 10 to 20 merges = +10
	// > 20 merges = +15 codeActivity points
	codeActivity, err = calculateCodeActivityIndex(r)
	if err != nil {
		log.Error(err)
	}

	// Calculate repoActivity index. (sum of other indexes)
	repoActivity = userCommunity + codeActivity + releaseHistory + longevity

	log.Debugf("Repoactivity: %f", repoActivity)
	log.Debugf("Repoactivity (userCommunity): %f", userCommunity)
	log.Debugf("Repoactivity (codeActivity): %f", codeActivity)
	log.Debugf("Repoactivity (releaseHistory): %f", releaseHistory)
	log.Debugf("Repoactivity (longevity): %f", longevity)

	return repoActivity, err
}

// extractLastYearMerges returns a slice of commits of merges from last year.
func extractLastYearMerges(cIter object.CommitIter) ([]object.Commit, error) {
	lastYear := time.Now().AddDate(-1, 0, 0) // last year
	var commits []object.Commit

	err := cIter.ForEach(func(c *object.Commit) error {
		// if Numparents > 1 is a merge, otherwise is a single commit.
		if c.Author.When.After(lastYear) && c.NumParents() > 1 {
			commits = append(commits, *c)
		}
		return nil
	})
	if err != nil {
		log.Error(err)
	}
	return commits, nil
}

// extractLastYearCommits returns a slice of commits from last year.
func extractLastYearCommits(cIter object.CommitIter) ([]object.Commit, error) {
	lastYear := time.Now().AddDate(-1, 0, 0) // last year
	var commits []object.Commit

	err := cIter.ForEach(func(c *object.Commit) error {
		if c.Author.When.After(lastYear) {
			commits = append(commits, *c)
		}
		return nil
	})
	if err != nil {
		log.Error(err)
	}
	return commits, nil
}

// extractOldestCommitDate returns the oldest commit date.
func extractOldestCommitDate(cIter object.CommitIter) (time.Time, error) {
	// Iterates over the commits and extract infos.
	result := time.Now()
	err := cIter.ForEach(func(c *object.Commit) error {
		if c.Author.When.Before(result) {
			result = c.Author.When
		}
		return nil
	})
	if err != nil {
		log.Error(err)
	}
	return result, nil
}

// extractAuthorsCommits returns a map of email-number of commits for the repository.
func extractAuthorsCommits(cIter object.CommitIter) (map[string]int, error) {
	totalAuthors := make(map[string]int)

	// Iterates over the commits and extract infos.
	err := cIter.ForEach(func(c *object.Commit) error {
		totalAuthors[c.Author.Email]++
		return nil
	})
	if err != nil {
		log.Error(err)
	}

	return totalAuthors, nil

}

func between(v, min, max int) bool {
	if v > min && v <= max {
		return true
	}
	return false
}

func calculateUserCommunityIndex(r *git.Repository) (float64, error) {
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

	authors, err := extractAuthorsCommits(cIter)
	if err != nil {
		log.Error(err)
	}
	if len(authors) < 25 {
		return float64(len(authors)), err
	} else if len(authors) > 25 {
		return 25, err
	}

	return 0, err
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
	creationDate, err := extractOldestCommitDate(cIter)
	if err != nil {
		log.Error(err)
	}
	if time.Since(creationDate).Hours() < 24*182.5 { // 6 months
		return 5, err
	} else if time.Since(creationDate).Hours() < 24*365 { // 1 year
		return 10, err
	} else if time.Since(creationDate).Hours() < 24*730 { // 2 year
		return 15, err
	} else if time.Since(creationDate).Hours() < 24*1095 { // 3 year
		return 20, err
	} else if time.Since(creationDate).Hours() > 24*1095 { // > 3 years
		return 25, err
	}

	return 0, err
}

func calculateReleaseHistoryIndex(r *git.Repository) (float64, error) {
	tags, err := r.TagObjects()
	if err != nil {
		log.Error(err)
		return 0, err
	}
	sumTags := 0
	err = tags.ForEach(func(t *object.Tag) error {
		sumTags++
		return nil
	})
	if err != nil {
		log.Error(err)
		return 0, err
	}
	if between(sumTags, 0, 1) { // 1 tag
		return 5, err
	} else if between(sumTags, 1, 5) { // < 5 tags
		return 10, err
	} else if sumTags > 5 {
		return 15, err // > 5 tags
	}

	return 0, err

}

func calculateCodeActivityIndex(r *git.Repository) (float64, error) {
	var codeActivity float64

	ref, err := r.Head()
	if err != nil {
		log.Error(err)
		return 0, err
	}

	// Commits.
	// Return to HEAD.
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Error(err)
		return codeActivity, err
	}
	commitsLastYear, err := extractLastYearCommits(cIter)
	if err != nil {
		log.Error(err)
		return codeActivity, err
	}
	if len(commitsLastYear) < 100 { // 100 commits last year
		codeActivity += 5
	} else if len(commitsLastYear) < 200 { // 200 commits last year
		codeActivity += 10
	} else {
		codeActivity += 15 // > 200 commits last year
	}

	// Merges.
	// Return to HEAD.
	cIter, err = r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Error(err)
		return codeActivity, err
	}
	// Merges last year.
	mergesLastyear, err := extractLastYearMerges(cIter)
	if err != nil {
		log.Error(err)
		return codeActivity, err
	}
	if len(mergesLastyear) < 10 { // 10 merges last year
		codeActivity += 5
	} else if len(mergesLastyear) < 20 { // 20 merges last year
		codeActivity += 10
	} else {
		codeActivity += 15 // > 20 merges last year
	}

	return codeActivity, err

}
