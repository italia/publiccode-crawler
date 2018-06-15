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

// CalculateRepoActivity return the repository activity calculated on the git clone.
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
	repoActivity := 0.0

	// Open and load the git repo path.
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Error(err)
	}

	// Retrieves the branch pointed by HEAD.
	ref, err := r.Head()
	if err != nil {
		log.Error(err)
	}

	// Retrieves the commit history.
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Error(err)
	}

	// Total authors.
	authors, err := extractAuthorsCommits(cIter)
	if err != nil {
		log.Error(err)
	}
	if between(len(authors), 0, 5) {
		repoActivity += 5
	} else if between(len(authors), 5, 10) {
		repoActivity += 10
	} else if between(len(authors), 10, 15) {
		repoActivity += 15
	} else if between(len(authors), 15, 20) {
		repoActivity += 20
	} else {
		repoActivity += 25
	}

	// Return to HEAD.
	cIter, err = r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Error(err)
	}
	// Oldest (first) commit.
	creationDate, err := extractOldestCommitDate(cIter)
	if err != nil {
		log.Error(err)
	}
	if time.Now().Sub(creationDate).Hours() < 24*182.5 { // 6 months
		repoActivity += 1
	} else if time.Now().Sub(creationDate).Hours() < 24*365 { // 1 year
		repoActivity += 10
	} else if time.Now().Sub(creationDate).Hours() < 24*730 { // 2 year
		repoActivity += 15
	} else if time.Now().Sub(creationDate).Hours() < 24*1095 { // 3 year
		repoActivity += 20
	} else {
		repoActivity += 25
	}

	// Tags.
	tags, err := r.TagObjects()
	if err != nil {
		log.Error(err)
	}
	sumTags := 0
	err = tags.ForEach(func(t *object.Tag) error {
		sumTags++
		return nil
	})
	if between(sumTags, 0, 1) { // 1 tag
		repoActivity += 5
	} else if between(sumTags, 1, 5) { // 2 tag
		repoActivity += 10
	} else {
		repoActivity += 15 // > 2 tag
	}

	// Return to HEAD.
	cIter, err = r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Error(err)
	}
	// Commits last year.
	commitsLastYear, err := extractLastYearCommits(cIter)
	if err != nil {
		log.Error(err)
	}
	if len(commitsLastYear) < 100 { // 100 commits last year
		repoActivity += 5
	} else if len(commitsLastYear) < 200 { // 200 commits last year
		repoActivity += 10
	} else {
		repoActivity += 15 // > 200 commits last year
	}

	// Merges.
	cIter, err = r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Error(err)
	}
	// Merges last year.
	mergesLastyear, err := extractLastYearMerges(cIter)
	if err != nil {
		log.Error(err)
	}
	if len(mergesLastyear) < 10 { // 10 merges last year
		repoActivity += 5
	} else if len(mergesLastyear) < 20 { // 20 merges last year
		repoActivity += 10
	} else {
		repoActivity += 15 // > 20 merges last year
	}

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

// extractCommitsPerMonth return a map[month]numberOfCommits.
func extractCommitsPerMonth(cIter object.CommitIter) (map[string]int, error) {
	// Authors per month.
	commitsMonth := make(map[string]int)
	// Iterates over the commits and extract infos.
	err := cIter.ForEach(func(c *object.Commit) error {
		monthYear := c.Author.When.Format("01-2006")
		commitsMonth[monthYear]++
		return nil
	})
	if err != nil {
		log.Error(err)
	}

	return commitsMonth, err
}

// extractNumberAuthorsPerMonth return a map[month]numberOfAuthors.
func extractNumberAuthorsPerMonth(cIter object.CommitIter) (map[string]int, error) {
	// Authors per month.
	authorsPerMonth := make(map[string]int)
	differentAuthorsMonth := map[string]map[string]int{}

	// Iterates over the commits and extract infos.
	err := cIter.ForEach(func(c *object.Commit) error {
		monthYear := c.Author.When.Format("01-2006")
		// Number of commits per author.
		if differentAuthorsMonth[monthYear] == nil {
			differentAuthorsMonth[monthYear] = make(map[string]int)
		}
		differentAuthorsMonth[monthYear][c.Author.Email] = differentAuthorsMonth[monthYear][c.Author.Email] + 1

		// Number authors per month.
		for month, authors := range differentAuthorsMonth {
			authorsPerMonth[month] = len(authors)
		}

		return nil
	})

	return authorsPerMonth, err
}

func between(v, min, max int) bool {
	if v > min && v <= max {
		return true
	}
	return false
}
