package crawler

import (
	"io/ioutil"
    "net/url"
	"testing"

	publiccode "github.com/italia/publiccode-parser-go/v2"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

var whitelist string = `
-
  name: testit
  repos: 
    - "https://github.com/test/testrepo"
-
  codice-iPA:
  name: testit
  repos: 
    - "https://github.com/test/testrepo"
-
  codice-iPA: test
  name: testit
  repos: 
    - "https://github.com/test/testrepo"
`

// validateRemoteFile will parse and validate
// crawled publiccode.yml. It will thrown errors
// if parse fails and if IPA code mismatch between
// whithelist file and publiccode itself
func TestIPAMatch(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	var pas []PA
	var parser publiccode.Parser

    u, _ := url.Parse("https://github.com/a/b/blob/main/publiccode.yml")
    parser.PublicCode.URL = (*publiccode.URL)(u)

	err := yaml.Unmarshal([]byte(whitelist), &pas)
	if err != nil {
		t.Errorf("error on unmarsalling whitelist %s", err)
	}

	// should not throw error codiceIPA key is equal
	// on both sides
	for _, pa := range pas {
		parser.PublicCode.It.Riuso.CodiceIPA = pa.CodiceIPA
		err = validateFile(pa, parser, "https://raw.githubusercontent.com/a/b/main/publiccode.yml")
		if err != nil {
			t.Errorf("error comparing IPA codes %s", err)
		}
	}

	// it should thowns errors since they always mismatch
	for _, pa := range pas {
		parser.PublicCode.It.Riuso.CodiceIPA = pa.CodiceIPA + "x"
		err = validateFile(pa, parser, "https://raw.githubusercontent.com/a/b/main/publiccode.yml")
		if err == nil {
			t.Errorf("error comparing IPA codes %v", err)
		}
	}
}

func createFakeRepo(name, gitCloneURL string) (r Repository) {
	r.Name = name
	r.GitCloneURL = gitCloneURL
	return
}

func TestRemovingRepoAsBlacklisted(t *testing.T) {
	var c Crawler
	// Faking repositories
	c.repositories = make(chan Repository, 3)
	c.repositories <- createFakeRepo("repo1", "https://github.com/italia/repo1.git")
	c.repositories <- createFakeRepo("repo2", "https://github.com/italia/repo2.git")
	c.repositories <- createFakeRepo("repo3", "https://github.com/italia/repo3.git")
	close(c.repositories)

	// Faking blacklist entries
	var repoListed = make(map[string]string)
	repoListed["https://github.com/italia/repo1.git"] = "https://github.com/italia/repo1"
	repoListed["https://github.com/italia/repo3.git"] = "https://github.com/italia/repo3"

	toBeRemoved := c.removeBlackListedFromRepositories(repoListed)

	assert.Len(t, toBeRemoved, 2)
	for _, entry := range toBeRemoved {
		assert.NotEmpty(t, repoListed[appendGitExt(entry)])
	}
}
