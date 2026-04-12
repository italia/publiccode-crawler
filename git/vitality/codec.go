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
	D0  string      `json:"d0"`
	A   []string    `json:"a"`
	E   []jsonEntry `json:"e"`
	T   []jsonTag   `json:"t"`
}

type jsonEntry struct {
	Delta uint32  `json:"delta"`
	C     uint32  `json:"c"`
	M     uint32  `json:"m"`
	A     []uint8 `json:"a"`
}

type jsonTag struct {
	Delta uint32 `json:"delta"`
	N     uint32 `json:"n"`
}

func Marshal(c Cache) ([]byte, error) {
	j := jsonCache{
		LU:  c.LastUpdated.UTC().Format(dateLayout),
		OCD: c.OldestCommitDate.UTC().Format(dateLayout),
		D0:  c.FirstEntryDate.UTC().Format(dateLayout),
		A:   c.Authors,
	}

	for _, e := range c.Entries {
		a := e.Authors
		if a == nil {
			a = []uint8{}
		}
		j.E = append(j.E, jsonEntry{Delta: e.Delta, C: e.Commits, M: e.Merges, A: a})
	}

	for _, t := range c.Tags {
		j.T = append(j.T, jsonTag{Delta: t.Delta, N: t.Count})
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
	if c.FirstEntryDate, err = time.Parse(dateLayout, j.D0); err != nil {
		return c, fmt.Errorf("vitality unmarshal d0: %w", err)
	}

	c.Authors = j.A

	for _, e := range j.E {
		a := e.A
		if a == nil {
			a = []uint8{}
		}
		c.Entries = append(c.Entries, DayEntry{Delta: e.Delta, Commits: e.C, Merges: e.M, Authors: a})
	}

	for _, t := range j.T {
		c.Tags = append(c.Tags, TagEntry{Delta: t.Delta, Count: t.N})
	}

	return c, nil
}

func EntryDate(c Cache, i int) time.Time {
	d := c.FirstEntryDate
	for j := 0; j <= i; j++ {
		d = d.AddDate(0, 0, int(c.Entries[j].Delta))
	}

	return d
}

func TagDate(c Cache, i int) time.Time {
	d := c.FirstEntryDate
	for j := 0; j <= i; j++ {
		d = d.AddDate(0, 0, int(c.Tags[j].Delta))
	}

	return d
}
