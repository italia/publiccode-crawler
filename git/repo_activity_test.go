package git

import (
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

// fixedNow is a fixed point in time used to make all tests deterministic.
var fixedNow = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func makeCommit(email string, when time.Time, numParents int) *object.Commit {
	c := &object.Commit{
		Author: object.Signature{Email: email, When: when},
	}
	for range numParents {
		c.ParentHashes = append(c.ParentHashes, plumbing.ZeroHash)
	}
	return c
}

func daysAgo(n int) time.Time {
	return fixedNow.AddDate(0, 0, -n)
}

func TestUserCommunityEmpty(t *testing.T) {
	assert.Equal(t, float64(0), userCommunityLastDays(nil))
}

func TestUserCommunityDeduplicatesEmails(t *testing.T) {
	commits := []*object.Commit{
		makeCommit("mario@example.com", daysAgo(1), 1),
		makeCommit("mario@example.com", daysAgo(2), 1),
		makeCommit("mario@example.com", daysAgo(3), 1),
	}
	assert.Equal(t, float64(1), userCommunityLastDays(commits))
}

func TestUserCommunityCountsDistinctEmails(t *testing.T) {
	commits := []*object.Commit{
		makeCommit("mario@example.com", daysAgo(1), 1),
		makeCommit("luigi@example.com", daysAgo(1), 1),
		makeCommit("mario@example.com", daysAgo(2), 1),
	}
	assert.Equal(t, float64(2), userCommunityLastDays(commits))
}

func TestActivityEmpty(t *testing.T) {
	assert.Equal(t, float64(0), activityLastDays(nil))
}

func TestActivityCountsCommits(t *testing.T) {
	commits := []*object.Commit{
		makeCommit("mario@example.com", daysAgo(1), 1),
		makeCommit("luigi@example.com", daysAgo(1), 1),
	}
	assert.Equal(t, float64(2), activityLastDays(commits))
}

func TestActivityMergeCountedDouble(t *testing.T) {
	// merge commit (numParents=2) counts as 1 commit + 1 merge = 2
	commits := []*object.Commit{
		makeCommit("mario@example.com", daysAgo(1), 1),
		makeCommit("luigi@example.com", daysAgo(1), 2),
	}
	assert.Equal(t, float64(3), activityLastDays(commits))
}

func TestReleaseHistoryEmpty(t *testing.T) {
	assert.Equal(t, float64(0), releaseHistoryLastDays(nil))
}

func TestReleaseHistoryCounts(t *testing.T) {
	tags := []*object.Commit{
		makeCommit("", daysAgo(1), 1),
		makeCommit("", daysAgo(5), 1),
		makeCommit("", daysAgo(10), 1),
	}
	assert.Equal(t, float64(3), releaseHistoryLastDays(tags))
}

func TestExtractCommitsLastDaysAll(t *testing.T) {
	commits := []*object.Commit{
		makeCommit("mario@example.com", daysAgo(1), 1),
		makeCommit("luigi@example.com", daysAgo(7), 1),
		makeCommit("luca@example.com", daysAgo(30), 1),
	}
	result := extractCommitsLastDays(40, commits, fixedNow)

	assert.Len(t, result[0], 3)  // all before fixedNow
	assert.Len(t, result[2], 2)  // -7 and -30 visible, -1 is not before fixedNow-2
	assert.Len(t, result[8], 1)  // only -30 visible
	assert.Len(t, result[31], 0) // none visible
}

func TestExtractCommitsLastDaysBoundaryExclusive(t *testing.T) {
	// commit exactly 5 days ago: visible on day 4, NOT on day 5 (Before is strict)
	commits := []*object.Commit{
		makeCommit("mario@example.com", daysAgo(5), 1),
	}
	result := extractCommitsLastDays(10, commits, fixedNow)

	assert.Len(t, result[4], 1) // fixedNow-5 < fixedNow-4: visible
	assert.Len(t, result[5], 0) // fixedNow-5 is not before fixedNow-5: not visible
}

func TestExtractCommitsPerDayGroups(t *testing.T) {
	commits := []*object.Commit{
		makeCommit("mario@example.com", daysAgo(1), 1),
		makeCommit("luigi@example.com", daysAgo(1), 1),
		makeCommit("luca@example.com", daysAgo(7), 1),
	}
	result := extractCommitsPerDay(10, commits, fixedNow)

	assert.Len(t, result[1], 2)
	assert.Len(t, result[7], 1)
	assert.Len(t, result[0], 0)
}

func TestMeanActivity(t *testing.T) {
	assert.Equal(t, float64(20), meanActivity(map[int]float64{0: 10, 1: 20, 2: 30}))
}

func TestMeanActivitySingle(t *testing.T) {
	assert.Equal(t, float64(42), meanActivity(map[int]float64{0: 42}))
}
