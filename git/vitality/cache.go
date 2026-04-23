package vitality

import "time"

type Cache struct {
	LastUpdated      time.Time
	OldestCommitDate time.Time
	FirstEntryDate   time.Time
	Authors          []string
	Entries          []DayEntry
	Tags             []TagEntry
}

type DayEntry struct {
	Delta   uint32
	Commits uint32
	Merges  uint32
	Authors []uint16
}

type TagEntry struct {
	Delta uint32
	Count uint32
}
