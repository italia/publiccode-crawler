package vitality

import "time"

type Cache struct {
	LastUpdated      time.Time
	OldestCommitDate time.Time
	Entries          []DayEntry
	Tags             []TagEntry
}

type DayEntry struct {
	Date    time.Time
	Commits uint32
	Merges  uint32
	Authors []string
}

type TagEntry struct {
	Date  time.Time
	Count uint32
}
