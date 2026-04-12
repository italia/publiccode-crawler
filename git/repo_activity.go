package git

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/italia/publiccode-crawler/v4/common"
	"github.com/italia/publiccode-crawler/v4/git/vitality"
	log "github.com/sirupsen/logrus"
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
func CalculateRepoActivity(repository common.Repository, days int, now time.Time) (float64, map[int]float64, error) {
	if repository.Name == "" {
		return 0, nil, errors.New("cannot calculate repository activity without name")
	}

	vendor, repo := common.SplitFullName(repository.Name)
	path := vitilityCachePath(repository.URL.Host, vendor, repo)

	data, err := os.ReadFile(path)
	if err != nil {
		return 0, nil, fmt.Errorf("cannot read vitality cache: %w", err)
	}

	cache, err := vitality.Unmarshal(data)
	if err != nil {
		return 0, nil, fmt.Errorf("cannot unmarshal vitality cache: %w", err)
	}

	longevity := now.Sub(cache.OldestCommitDate).Hours() / 24

	then := time.Date(2005, time.January, 1, 1, 0, 0, 0, time.UTC)
	if longevity > now.Sub(then).Hours()/24 {
		return 0, nil, errors.New("first commit is too old, must be after 2005")
	}

	vitalityIndex := map[int]float64{}

	for i := range days {
		cutoff := now.AddDate(0, 0, -i)

		authorSet := map[uint8]struct{}{}
		var commits, merges float64

		cur := cache.FirstEntryDate
		for _, e := range cache.Entries {
			cur = cur.AddDate(0, 0, int(e.Delta))
			if cur.Before(cutoff) {
				for _, id := range e.Authors {
					authorSet[id] = struct{}{}
				}
			}
			if sameDay(cur, cutoff) {
				commits = float64(e.Commits)
				merges = float64(e.Merges)
			}
		}

		var tagCount float64
		cur = cache.FirstEntryDate
		for _, t := range cache.Tags {
			cur = cur.AddDate(0, 0, int(t.Delta))
			if sameDay(cur, cutoff) {
				tagCount = float64(t.Count)
			}
		}

		repoActivity := ranges("userCommunity", float64(len(authorSet))) +
			ranges("codeActivity", commits+merges) +
			ranges("releaseHistory", tagCount) +
			ranges("longevity", longevity)

		if repoActivity > 100 {
			repoActivity = 100
		}

		vitalityIndex[i] = repoActivity
	}

	total := meanActivity(vitalityIndex)
	if total > 100 {
		total = 100
	}

	return float64(int(total)), vitalityIndex, nil
}

func ranges(name string, value float64) float64 {
	data, err := os.ReadFile("vitality-ranges.yml")
	if err != nil {
		log.Error(err)
	}

	t := RangesData{}

	err = yaml.Unmarshal(data, &t)
	if err != nil {
		log.Errorf("error: %v", err)
	}

	for _, v := range t {
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

func sameDay(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}

func meanActivity(points map[int]float64) float64 {
	var total float64
	for _, point := range points {
		total += point
	}

	return total / float64(len(points))
}
