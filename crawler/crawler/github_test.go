package crawler

import (
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"
)

// GenerateGithubAPIURL returns the api url of given Gitlab organization link.
// IN: https://github.com/italia
// OUT:https://api.github.com/orgs/italia/repos
func TestGenerateGithubAPIURL(t *testing.T) {
	// Disablle log output for this function
	log.SetOutput(ioutil.Discard)

	links := []struct {
		in  string
		out string
	}{
		{"https://github.com/italia", "https://api.github.com/orgs/italia/repos"},
		{":unparsable", ":unparsable"},
	}

	for _, l := range links {
		genURL := GenerateGithubAPIURL()
		if out, err := genURL(l.in); out[0] != l.out {
			t.Logf("Expected %s == %s: %v ", out[0], l.out, err)
			t.Fail()
		}
	}

}
