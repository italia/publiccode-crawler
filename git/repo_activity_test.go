package git

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// fixedNow is a fixed point in time used to make all tests deterministic.
var fixedNow = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func TestMeanActivity(t *testing.T) {
	assert.Equal(t, float64(20), meanActivity(map[int]float64{0: 10, 1: 20, 2: 30}))
}

func TestMeanActivitySingle(t *testing.T) {
	assert.Equal(t, float64(42), meanActivity(map[int]float64{0: 42}))
}
