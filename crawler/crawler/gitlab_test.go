package crawler

import (
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"
)

// IsGitlab returns "true" if the url can use Gitlab API.
func TestIsGitlab(t *testing.T) {
	// Disablle log output for this function
	log.SetOutput(ioutil.Discard)

	links := []struct {
		in  string
		out bool
	}{
		// {"https://bitbucket.org/Soft", false},
		// {"https://gitlab.com/Soft", true},
		// {"https://github.com/Soft", false},
		// {"", false},
		// {"invalidUrl", false},
		// {"example.example", false},
		// {":unparsable", false},
	}

	for _, l := range links {
		if IsGitlab(l.in) != l.out {
			t.Logf("Expected %s == %t.", l.in, l.out)
			t.Fail()
		}
	}

}

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
		{":unparsable", ":unparsable"},
	}

	for _, l := range links {
		genURL := GenerateGitlabAPIURL()
		if out, err := genURL(l.in); out[0] != l.out {
			t.Logf("Expected %s == %s: %v ", out[0], l.out, err)
			t.Fail()
		}
	}

}
