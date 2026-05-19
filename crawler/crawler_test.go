package crawler

import (
	"net/url"
	"testing"

	"github.com/italia/publiccode-crawler/v4/common"
	publiccode "github.com/italia/publiccode-parser-go/v5"
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

func newPublicCode(t *testing.T, repoURL, organisationURI string) publiccode.PublicCodeV0 {
	t.Helper()

	u, err := url.Parse(repoURL)
	assert.NoError(t, err)

	return publiccode.PublicCodeV0{
		URL: (*publiccode.URL)(u),
		Organisation: &publiccode.OrganisationV0{
			URI: organisationURI,
		},
	}
}

func TestValidateFile_NoAlternativeIDSkipsCheck(t *testing.T) {
	pc := newPublicCode(t,
		"https://github.com/mastodon/mastodon",
		"https://joinmastodon.org",
	)
	publisher := common.Publisher{
		ID:   "9b9aa9c5-30b1-4e56-b3c0-5b6e2e6d3b22",
		Name: "Mastodon community",
	}

	err := validateFile("urn:x-italian-pa:", publisher, pc,
		"https://raw.githubusercontent.com/mastodon/mastodon/main/publiccode.yml")

	assert.NoError(t, err)
}

// With no namespace, the expected value is the bare alternativeId.
func TestValidateFile_NoNamespaceMatchesAlternativeID(t *testing.T) {
	pc := newPublicCode(t,
		"https://github.com/foo/bar",
		"pcm",
	)
	publisher := common.Publisher{ID: "pcm", AlternativeID: "pcm", Name: "PCM"}

	err := validateFile("", publisher, pc,
		"https://raw.githubusercontent.com/foo/bar/main/publiccode.yml")

	assert.NoError(t, err)
}

func TestValidateFile_NoNamespaceMismatch(t *testing.T) {
	pc := newPublicCode(t,
		"https://github.com/foo/bar",
		"https://example.com",
	)
	publisher := common.Publisher{ID: "pcm", AlternativeID: "pcm", Name: "PCM"}

	err := validateFile("", publisher, pc,
		"https://raw.githubusercontent.com/foo/bar/main/publiccode.yml")

	assert.Error(t, err)
}

func TestValidateFile_OrganisationMatchesPublisher(t *testing.T) {
	pc := newPublicCode(t,
		"https://github.com/foo/bar",
		"urn:x-italian-pa:pcm",
	)
	publisher := common.Publisher{ID: "pcm", AlternativeID: "pcm", Name: "PCM"}

	err := validateFile("urn:x-italian-pa:", publisher, pc,
		"https://raw.githubusercontent.com/foo/bar/main/publiccode.yml")

	assert.NoError(t, err)
}

func TestValidateFile_OrganisationMismatch(t *testing.T) {
	pc := newPublicCode(t,
		"https://github.com/foo/bar",
		"urn:x-italian-pa:wrong",
	)
	publisher := common.Publisher{ID: "pcm", AlternativeID: "pcm", Name: "PCM"}

	err := validateFile("urn:x-italian-pa:", publisher, pc,
		"https://raw.githubusercontent.com/foo/bar/main/publiccode.yml")

	assert.Error(t, err)
}
