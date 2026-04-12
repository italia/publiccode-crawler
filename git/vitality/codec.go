package vitality

import (
	"encoding/json"
	"fmt"
	"time"
)

const dateLayout = "2006-01-02"

type jsonCache struct {
	LastUpdated      string      `json:"last_updated"`
	OldestCommitDate string      `json:"oldest_commit_date"`
	Entries          []jsonEntry `json:"entries"`
	Tags             []jsonTag   `json:"tags"`
}

type jsonEntry struct {
	Date    string   `json:"date"`
	Commits uint32   `json:"commits"`
	Merges  uint32   `json:"merges"`
	Authors []string `json:"authors"`
}

type jsonTag struct {
	Date  string `json:"date"`
	Count uint32 `json:"count"`
}

func Marshal(c Cache) ([]byte, error) {
	j := jsonCache{
		LastUpdated:      c.LastUpdated.UTC().Format(dateLayout),
		OldestCommitDate: c.OldestCommitDate.UTC().Format(dateLayout),
	}

	for _, e := range c.Entries {
		authors := e.Authors
		if authors == nil {
			authors = []string{}
		}
		j.Entries = append(j.Entries, jsonEntry{
			Date:    e.Date.UTC().Format(dateLayout),
			Commits: e.Commits,
			Merges:  e.Merges,
			Authors: authors,
		})
	}

	for _, t := range c.Tags {
		j.Tags = append(j.Tags, jsonTag{
			Date:  t.Date.UTC().Format(dateLayout),
			Count: t.Count,
		})
	}

	return json.MarshalIndent(j, "", "  ")
}

func Unmarshal(data []byte) (Cache, error) {
	if len(data) == 0 {
		return Cache{}, fmt.Errorf("vitality unmarshal: empty data")
	}

	var j jsonCache
	if err := json.Unmarshal(data, &j); err != nil {
		return Cache{}, fmt.Errorf("vitality unmarshal: %w", err)
	}

	var c Cache
	var err error

	if c.LastUpdated, err = time.Parse(dateLayout, j.LastUpdated); err != nil {
		return c, fmt.Errorf("vitality unmarshal last_updated: %w", err)
	}
	if c.OldestCommitDate, err = time.Parse(dateLayout, j.OldestCommitDate); err != nil {
		return c, fmt.Errorf("vitality unmarshal oldest_commit_date: %w", err)
	}

	for _, e := range j.Entries {
		date, err := time.Parse(dateLayout, e.Date)
		if err != nil {
			return c, fmt.Errorf("vitality unmarshal entry date %q: %w", e.Date, err)
		}
		authors := e.Authors
		if authors == nil {
			authors = []string{}
		}
		c.Entries = append(c.Entries, DayEntry{
			Date:    date,
			Commits: e.Commits,
			Merges:  e.Merges,
			Authors: authors,
		})
	}

	for _, t := range j.Tags {
		date, err := time.Parse(dateLayout, t.Date)
		if err != nil {
			return c, fmt.Errorf("vitality unmarshal tag date %q: %w", t.Date, err)
		}
		c.Tags = append(c.Tags, TagEntry{
			Date:  date,
			Count: t.Count,
		})
	}

	return c, nil
}
