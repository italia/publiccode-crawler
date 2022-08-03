package crawler

import (
	"io/ioutil"
	"net/url"
	"testing"

	log "github.com/sirupsen/logrus"
)

// GenerateGitlabAPIURL returns the api url of given Gitlab organization link.
// IN: https://gitlab.org/blockninja
// OUT:https://gitlab.com/api/v4/groups/blockninja
func TestGenerateGitlabAPIURL(t *testing.T) {
	// Disablle log output for this function
	log.SetOutput(ioutil.Discard)

	links := []struct {
		in  string
		out string
	}{
		{"https://gitlab.com/blockninja", "https://gitlab.com/api/v4/groups/blockninja"},
	}

	for _, l := range links {
		genURL := GenerateGitlabAPIURL()

		u, _ := url.Parse(l.in)
		if out, err := genURL(*u); out[0].String() != l.out {
			t.Logf("Expected %s == %s: %v ", out[0].String(), l.out, err)
			t.Fail()
		}
	}

}
