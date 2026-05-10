package vitality

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// version prefixes every cache file. Bump on any incompatible binary layout
// change so old caches are rejected instead of misread.
const version uint8 = 1

// epoch is the cache's date anchor: uint16 days from this point covers
// 1980-01-01 through ~2159-06-06.
var epoch = time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)

func daysFromEpoch(date time.Time) uint16 {
	days := date.UTC().Unix()/86400 - epoch.Unix()/86400
	if days < 0 {
		log.Warnf("vitality cache: date %s is before epoch %s; clamping",
			date.UTC().Format(time.RFC3339), epoch.Format(time.RFC3339))

		return 0
	}

	if days > 65535 {
		log.Warnf("vitality cache: date %s is after max representable date; clamping",
			date.UTC().Format(time.RFC3339))

		return 65535
	}

	return uint16(days)
}

func daysToTime(days uint16) time.Time {
	return epoch.AddDate(0, 0, int(days))
}

func Marshal(cache Cache) ([]byte, error) {
	var buf bytes.Buffer
	write := func(v any) { binary.Write(&buf, binary.LittleEndian, v) } //nolint:errcheck

	buf.WriteByte(version)

	write(daysFromEpoch(cache.LastUpdated))
	write(daysFromEpoch(cache.OldestCommitDate))
	write(daysFromEpoch(cache.FirstEntryDate))

	if len(cache.Authors) > 65535 {
		return nil, fmt.Errorf("vitality marshal: too many authors (%d)", len(cache.Authors))
	}

	write(uint16(len(cache.Authors))) //nolint:gosec // bounded above
	for _, author := range cache.Authors {
		buf.WriteString(author)
		buf.WriteByte(0)
	}

	if len(cache.Tags) > 65535 {
		return nil, fmt.Errorf("vitality marshal: too many tag days (%d)", len(cache.Tags))
	}

	write(uint16(len(cache.Tags))) //nolint:gosec // bounded above
	for _, tag := range cache.Tags {
		if tag.Delta > 65535 || tag.Count > 65535 {
			return nil, fmt.Errorf("vitality marshal: tag delta %d or count %d out of range", tag.Delta, tag.Count)
		}

		write(uint16(tag.Delta))
		write(uint16(tag.Count))
	}

	if len(cache.Entries) > 65535 {
		return nil, fmt.Errorf("vitality marshal: too many entries (%d)", len(cache.Entries))
	}

	write(uint16(len(cache.Entries))) //nolint:gosec // bounded above
	for _, entry := range cache.Entries {
		if entry.Delta > 65535 || entry.Commits > 65535 || entry.Merges > 65535 {
			return nil, fmt.Errorf("vitality marshal: entry out of range (delta=%d commits=%d merges=%d)",
				entry.Delta, entry.Commits, entry.Merges)
		}

		write(uint16(entry.Delta))
		write(uint16(entry.Commits))
		write(uint16(entry.Merges))

		if len(entry.Authors) > 65535 {
			return nil, fmt.Errorf("vitality marshal: too many authors in entry (%d)", len(entry.Authors))
		}

		write(uint16(len(entry.Authors))) //nolint:gosec // bounded above
		for _, author := range entry.Authors {
			write(author)
		}
	}

	return buf.Bytes(), nil
}

//nolint:funlen // Linear binary decoder, splitting hurts readability.
func Unmarshal(data []byte) (Cache, error) {
	if len(data) == 0 {
		return Cache{}, errors.New("vitality unmarshal: empty data")
	}

	reader := bytes.NewReader(data)
	read := func(v any) error { return binary.Read(reader, binary.LittleEndian, v) }

	var cache Cache

	ver, err := reader.ReadByte()
	if err != nil {
		return cache, fmt.Errorf("vitality unmarshal: read version: %w", err)
	}

	if ver != version {
		return cache, fmt.Errorf("vitality unmarshal: unsupported version %d (expected %d)", ver, version)
	}

	var lastUpdated, oldestCommit, firstDay uint16
	if err := read(&lastUpdated); err != nil {
		return cache, fmt.Errorf("vitality unmarshal lastUpdated: %w", err)
	}

	if err := read(&oldestCommit); err != nil {
		return cache, fmt.Errorf("vitality unmarshal oldestCommit: %w", err)
	}

	if err := read(&firstDay); err != nil {
		return cache, fmt.Errorf("vitality unmarshal firstDay: %w", err)
	}

	cache.LastUpdated = daysToTime(lastUpdated)
	cache.OldestCommitDate = daysToTime(oldestCommit)
	cache.FirstEntryDate = daysToTime(firstDay)

	var numAuthors uint16
	if err := read(&numAuthors); err != nil {
		return cache, fmt.Errorf("vitality unmarshal num_authors: %w", err)
	}

	cache.Authors = make([]string, numAuthors)
	for idx := range cache.Authors {
		var nameBuf bytes.Buffer

		for {
			nextByte, err := reader.ReadByte()
			if err != nil {
				return cache, fmt.Errorf("vitality unmarshal author[%d]: %w", idx, err)
			}

			if nextByte == 0 {
				break
			}

			nameBuf.WriteByte(nextByte)
		}

		cache.Authors[idx] = nameBuf.String()
	}

	var numTags uint16
	if err := read(&numTags); err != nil {
		return cache, fmt.Errorf("vitality unmarshal num_tags: %w", err)
	}

	cache.Tags = make([]TagEntry, numTags)
	for idx := range cache.Tags {
		var delta, count uint16
		if err := read(&delta); err != nil {
			return cache, fmt.Errorf("vitality unmarshal tag[%d].delta: %w", idx, err)
		}

		if err := read(&count); err != nil {
			return cache, fmt.Errorf("vitality unmarshal tag[%d].count: %w", idx, err)
		}

		cache.Tags[idx] = TagEntry{Delta: uint32(delta), Count: uint32(count)}
	}

	var numEntries uint16
	if err := read(&numEntries); err != nil {
		return cache, fmt.Errorf("vitality unmarshal num_entries: %w", err)
	}

	cache.Entries = make([]DayEntry, numEntries)
	for idx := range cache.Entries {
		var delta, commits, merges, numA uint16
		if err := read(&delta); err != nil {
			return cache, fmt.Errorf("vitality unmarshal entry[%d].delta: %w", idx, err)
		}

		if err := read(&commits); err != nil {
			return cache, fmt.Errorf("vitality unmarshal entry[%d].commits: %w", idx, err)
		}

		if err := read(&merges); err != nil {
			return cache, fmt.Errorf("vitality unmarshal entry[%d].merges: %w", idx, err)
		}

		if err := read(&numA); err != nil {
			return cache, fmt.Errorf("vitality unmarshal entry[%d].num_authors: %w", idx, err)
		}

		authors := make([]uint16, numA)
		if err := binary.Read(reader, binary.LittleEndian, authors); err != nil {
			return cache, fmt.Errorf("vitality unmarshal entry[%d].authors: %w", idx, err)
		}

		cache.Entries[idx] = DayEntry{
			Delta:   uint32(delta),
			Commits: uint32(commits),
			Merges:  uint32(merges),
			Authors: authors,
		}
	}

	return cache, nil
}
