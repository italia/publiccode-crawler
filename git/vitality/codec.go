package vitality

import (
	"encoding/json"
	"fmt"
	"time"
)

const dateLayout = "2006-01-02"

type jsonCache struct {
	LU  string      `json:"lu"`
	OCD string      `json:"ocd"`
	E   []jsonEntry `json:"e"`
	T   []jsonTag   `json:"t"`
}

type jsonEntry struct {
	D string   `json:"d"`
	C uint32   `json:"c"`
	M uint32   `json:"m"`
	A []string `json:"a"`
}

type jsonTag struct {
	D string `json:"d"`
	N uint32 `json:"n"`
}

func Marshal(c Cache) ([]byte, error) {
	j := jsonCache{
		LU:  c.LastUpdated.UTC().Format(dateLayout),
		OCD: c.OldestCommitDate.UTC().Format(dateLayout),
	}

	for _, e := range c.Entries {
		a := e.Authors
		if a == nil {
			a = []string{}
		}
		j.E = append(j.E, jsonEntry{
			D: e.Date.UTC().Format(dateLayout),
			C: e.Commits,
			M: e.Merges,
			A: a,
		})
	}

	for _, t := range c.Tags {
		j.T = append(j.T, jsonTag{
			D: t.Date.UTC().Format(dateLayout),
			N: t.Count,
		})
	}

	return json.Marshal(j)
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

	if c.LastUpdated, err = time.Parse(dateLayout, j.LU); err != nil {
		return c, fmt.Errorf("vitality unmarshal lu: %w", err)
	}
	if c.OldestCommitDate, err = time.Parse(dateLayout, j.OCD); err != nil {
		return c, fmt.Errorf("vitality unmarshal ocd: %w", err)
	}

	for _, e := range j.E {
		date, err := time.Parse(dateLayout, e.D)
		if err != nil {
			return c, fmt.Errorf("vitality unmarshal entry date %q: %w", e.D, err)
		}
		a := e.A
		if a == nil {
			a = []string{}
		}
		c.Entries = append(c.Entries, DayEntry{
			Date:    date,
			Commits: e.C,
			Merges:  e.M,
			Authors: a,
		})
	}

	for _, t := range j.T {
		date, err := time.Parse(dateLayout, t.D)
		if err != nil {
			return c, fmt.Errorf("vitality unmarshal tag date %q: %w", t.D, err)
		}
		c.Tags = append(c.Tags, TagEntry{
			Date:  date,
			Count: t.N,
		})
	}

	return c, nil
}
