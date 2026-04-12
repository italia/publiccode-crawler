package vitality

import "time"

type Cache struct {
	LastUpdated      time.Time
	OldestCommitDate time.Time
	FirstEntryDate   time.Time
	Entries          []DayEntry
	Tags             []TagEntry
}

type DayEntry struct {
	Delta   uint32
	Commits uint32
	Merges  uint32
	Authors []string
}

type TagEntry struct {
	Delta uint32
	Count uint32
}
