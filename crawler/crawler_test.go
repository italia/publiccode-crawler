package crawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeNewAliases_addsNew(t *testing.T) {
	existing := []string{"https://github.com/org/repo.git"}
	newAliases := []string{"https://github.com/org/repo"}

	got := mergeNewAliases(existing, newAliases)

	assert.Equal(t, []string{
		"https://github.com/org/repo.git",
		"https://github.com/org/repo",
	}, got)
}

func TestMergeNewAliases_noDuplicates(t *testing.T) {
	existing := []string{"https://github.com/org/repo.git"}
	newAliases := []string{"https://github.com/org/repo.git"}

	got := mergeNewAliases(existing, newAliases)

	assert.Equal(t, []string{"https://github.com/org/repo.git"}, got)
}

func TestMergeNewAliases_emptyExisting(t *testing.T) {
	got := mergeNewAliases(nil, []string{"https://github.com/org/repo"})

	assert.Equal(t, []string{"https://github.com/org/repo"}, got)
}

func TestMergeNewAliases_emptyNew(t *testing.T) {
	existing := []string{"https://github.com/org/repo"}

	got := mergeNewAliases(existing, nil)

	assert.Equal(t, existing, got)
}

func TestMergeNewAliases_bothEmpty(t *testing.T) {
	got := mergeNewAliases(nil, nil)

	assert.Nil(t, got)
}

func TestMergeNewAliases_multipleNew(t *testing.T) {
	existing := []string{"a"}
	newAliases := []string{"b", "a", "c"}

	got := mergeNewAliases(existing, newAliases)

	assert.Equal(t, []string{"a", "b", "c"}, got)
}
