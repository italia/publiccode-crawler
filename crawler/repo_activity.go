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

// CalculateRepoActivity return the repository activity index and the vitality slice calculated on the git clone.
// It follows the document https://lg-acquisizione-e-riuso-software-per-la-pa.readthedocs.io/
// In reference to section: fase-2-2-valutazione-soluzioni-riusabili-per-la-pa
func CalculateRepoActivity(domain Domain, hostname string, name string) (float64, []int, error) {
	if domain.Host == "" {
		return 0, []int{}, errors.New("cannot calculate repository activity without domain host")
	}
	if name == "" {
		return 0, []int{}, errors.New("cannot  calculate repository activity without name")
	}

	vendor, repo := splitFullName(name)

	path := filepath.Join("./data", hostname, vendor, repo, "gitClone")

	// MkdirAll will create all the folder path, if not exists.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return 0, []int{}, err
	}

	// Repository activity score.
	var (
		userCommunity  float64 // max 30
		codeActivity   float64 // max 30
		releaseHistory float64 // max 15
		longevity      float64 // max 25

		repoActivity float64 // max 100
	)

	// Open and load the git repo path.
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Error(err)
	}

	// Total authors. (userCommunity index)
	userCommunity, err = calculateUserCommunityIndex(r)
	if err != nil {
		log.Error(err)
	}

	// Oldest (first) commit. (logevity index)
	longevity, err = calculateLongevityIndex(r)
	if err != nil {
		log.Error(err)
	}

	// Tags. (releaseHistory index)
	// 0 to 1 tag = +5
	// 1 to 5 tags = +10
	// > 5 tags = +15 releaseHistory points
	releaseHistory, err = calculateReleaseHistoryIndex(r)
	if err != nil {
		log.Error(err)
	}

	// Commits and merges last year. (codeActivity index)
	codeActivity, err = calculateCodeActivityIndex(r)
	if err != nil {
		log.Error(err)
	}

	// Commits and merges vitality for year. (commits per months)
	vitality, err := calculateVitality(r)
	if err != nil {
		log.Error(err)
	}

	// Calculate repoActivity index. (sum of all the others indexes)
	repoActivity = userCommunity + codeActivity + releaseHistory + longevity

	return repoActivity, vitality, err
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

// between returns true if the first integer parameter is between ]min,max].
func between(v, min, max int) bool {
	if v > min && v <= max {
		return true
	}
	return false
}

// calculateUserCommunityIndex find out how many authors there are in the repository.
// 0 to 30 = + #authors
// > 30 = +30 userCommunity
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
	if len(authors) < 30 {
		return float64(len(authors)), err
	} else if len(authors) > 30 {
		return 30, err
	}

	return 0, err
}

// calculateUserCommunityIndex find out the index based on the age of the repository.
// 0 to 6 months = +5
// 6 to 12 months = +10
// 12 to 24 months = +15
// 24 to 36 months = +20
// > 36 months = +25 longevity points
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

// calculateReleaseHistoryIndex find out the index based on the number of tags released.
// 0 to 1 tag = +5
// 1 to 5 tags = +10
// > 5 tags = +15 releaseHistory points
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

// calculateCodeActivityIndex find out the index based on the number of commits and merges from the repository.
// 0 to 100 commits = +5
// 1 to 200 commits = +10
// > 200 commmits = +15 codeActivity points
// And merges last year. (codeActivity index)
// 0 to 10 merges = +5
// 10 to 20 merges = +10
// > 20 merges = +15 codeActivity points
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

// calculateVitality returns monthly commits for the last year.
func calculateVitality(r *git.Repository) ([]int, error) {
	vitality := make([]int, 12)

	ref, err := r.Head()
	if err != nil {
		log.Error(err)
		return vitality, err
	}

	// Commits.
	// Return to HEAD.
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Error(err)
		return vitality, err
	}
	commitsLastYear, err := extractLastYearCommits(cIter)
	if err != nil {
		log.Error(err)
		return vitality, err
	}

	for _, commit := range commitsLastYear {
		vitality[commit.Author.When.Month()-1]++
	}
	// Rotate for sorting.
	monthNow := time.Now().Month()
	rotateL(vitality, int(monthNow))

	return vitality, err
}

// https://play.golang.org/p/UwrHJYskNS
func rotateL(a []int, i int) {
	// Ensure the shift amount is less than the length of the array,
	// and that it is positive.
	i = i % len(a)
	if i < 0 {
		i += len(a)
	}

	for c := 0; c < gcd(i, len(a)); c++ {
		t := a[c]
		j := c

		for {
			k := j + i
			// loop around if we go past the end of the slice
			if k >= len(a) {
				k -= len(a)
			}
			// end when we get to where we started
			if k == c {
				break
			}
			// move the element directly into its final position
			a[j] = a[k]
			j = k
		}

		a[j] = t
	}
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}

	return a
}
