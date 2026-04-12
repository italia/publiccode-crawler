package vitality

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// epoch is the cache's date anchor: uint16 days from this point covers
// 1980-01-01 through ~2159-06-06.
var epoch = time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)

func daysFromEpoch(t time.Time) uint16 {
	d := t.UTC().Unix()/86400 - epoch.Unix()/86400
	if d < 0 {
		log.Warnf("vitality cache: date %s is before epoch %s; clamping",
			t.UTC().Format(time.RFC3339), epoch.Format(time.RFC3339))

		return 0
	}

	if d > 65535 {
		log.Warnf("vitality cache: date %s is after max representable date; clamping",
			t.UTC().Format(time.RFC3339))

		return 65535
	}

	return uint16(d)
}

func daysToTime(d uint16) time.Time {
	return epoch.AddDate(0, 0, int(d))
}

func Marshal(c Cache) ([]byte, error) {
	var buf bytes.Buffer
	w := func(v any) { binary.Write(&buf, binary.LittleEndian, v) } //nolint:errcheck

	w(daysFromEpoch(c.LastUpdated))
	w(daysFromEpoch(c.OldestCommitDate))
	w(daysFromEpoch(c.FirstEntryDate))

	if len(c.Authors) > 65535 {
		return nil, fmt.Errorf("vitality marshal: too many authors (%d)", len(c.Authors))
	}

	w(uint16(len(c.Authors))) //nolint:gosec // bounded above
	for _, a := range c.Authors {
		buf.WriteString(a)
		buf.WriteByte(0)
	}

	if len(c.Tags) > 65535 {
		return nil, fmt.Errorf("vitality marshal: too many tag days (%d)", len(c.Tags))
	}

	w(uint16(len(c.Tags))) //nolint:gosec // bounded above
	for _, t := range c.Tags {
		if t.Delta > 65535 || t.Count > 65535 {
			return nil, fmt.Errorf("vitality marshal: tag delta %d or count %d out of range", t.Delta, t.Count)
		}

		w(uint16(t.Delta))
		w(uint16(t.Count))
	}

	if len(c.Entries) > 65535 {
		return nil, fmt.Errorf("vitality marshal: too many entries (%d)", len(c.Entries))
	}

	w(uint16(len(c.Entries))) //nolint:gosec // bounded above
	for _, e := range c.Entries {
		if e.Delta > 65535 || e.Commits > 65535 || e.Merges > 65535 {
			return nil, fmt.Errorf("vitality marshal: entry out of range (delta=%d commits=%d merges=%d)",
				e.Delta, e.Commits, e.Merges)
		}

		w(uint16(e.Delta))
		w(uint16(e.Commits))
		w(uint16(e.Merges))

		if len(e.Authors) > 65535 {
			return nil, fmt.Errorf("vitality marshal: too many authors in entry (%d)", len(e.Authors))
		}

		w(uint16(len(e.Authors))) //nolint:gosec // bounded above
		for _, a := range e.Authors {
			w(a)
		}
	}

	return buf.Bytes(), nil
}

func Unmarshal(data []byte) (Cache, error) {
	if len(data) == 0 {
		return Cache{}, errors.New("vitality unmarshal: empty data")
	}

	r := bytes.NewReader(data)
	rd := func(v any) error { return binary.Read(r, binary.LittleEndian, v) }

	var c Cache

	var lu, ocd, d0 uint16
	if err := rd(&lu); err != nil {
		return c, fmt.Errorf("vitality unmarshal lu: %w", err)
	}

	if err := rd(&ocd); err != nil {
		return c, fmt.Errorf("vitality unmarshal ocd: %w", err)
	}

	if err := rd(&d0); err != nil {
		return c, fmt.Errorf("vitality unmarshal d0: %w", err)
	}

	c.LastUpdated = daysToTime(lu)
	c.OldestCommitDate = daysToTime(ocd)
	c.FirstEntryDate = daysToTime(d0)

	var numAuthors uint16
	if err := rd(&numAuthors); err != nil {
		return c, fmt.Errorf("vitality unmarshal num_authors: %w", err)
	}

	c.Authors = make([]string, numAuthors)
	for i := range c.Authors {
		var sb bytes.Buffer

		for {
			b, err := r.ReadByte()
			if err != nil {
				return c, fmt.Errorf("vitality unmarshal author[%d]: %w", i, err)
			}

			if b == 0 {
				break
			}

			sb.WriteByte(b)
		}

		c.Authors[i] = sb.String()
	}

	var numTags uint16
	if err := rd(&numTags); err != nil {
		return c, fmt.Errorf("vitality unmarshal num_tags: %w", err)
	}

	c.Tags = make([]TagEntry, numTags)
	for i := range c.Tags {
		var delta, count uint16
		if err := rd(&delta); err != nil {
			return c, fmt.Errorf("vitality unmarshal tag[%d].delta: %w", i, err)
		}

		if err := rd(&count); err != nil {
			return c, fmt.Errorf("vitality unmarshal tag[%d].count: %w", i, err)
		}

		c.Tags[i] = TagEntry{Delta: uint32(delta), Count: uint32(count)}
	}

	var numEntries uint16
	if err := rd(&numEntries); err != nil {
		return c, fmt.Errorf("vitality unmarshal num_entries: %w", err)
	}

	c.Entries = make([]DayEntry, numEntries)
	for i := range c.Entries {
		var delta, commits, merges, numA uint16
		if err := rd(&delta); err != nil {
			return c, fmt.Errorf("vitality unmarshal entry[%d].delta: %w", i, err)
		}

		if err := rd(&commits); err != nil {
			return c, fmt.Errorf("vitality unmarshal entry[%d].commits: %w", i, err)
		}

		if err := rd(&merges); err != nil {
			return c, fmt.Errorf("vitality unmarshal entry[%d].merges: %w", i, err)
		}

		if err := rd(&numA); err != nil {
			return c, fmt.Errorf("vitality unmarshal entry[%d].num_authors: %w", i, err)
		}

		authors := make([]uint16, numA)
		if err := binary.Read(r, binary.LittleEndian, authors); err != nil {
			return c, fmt.Errorf("vitality unmarshal entry[%d].authors: %w", i, err)
		}

		c.Entries[i] = DayEntry{Delta: uint32(delta), Commits: uint32(commits), Merges: uint32(merges), Authors: authors}
	}

	return c, nil
}
