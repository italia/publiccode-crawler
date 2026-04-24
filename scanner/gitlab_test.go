package scanner

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsGitlabGroup(t *testing.T) {
	tests := []struct {
		rawURL string
		want   bool
	}{
		{"https://gitlab.com/mygroup", true},
		{"https://gitlab.com/mygroup/mysubgroup", true},
		{"https://mygitlab.example.com/mygroup", true},
		{"https://mygitlab.example.com/", false},
		{"https://mygitlab.example.com", false},
	}

	for _, tc := range tests {
		parsed, err := url.Parse(tc.rawURL)
		require.NoError(t, err)

		assert.Equal(t, tc.want, isGitlabGroup(*parsed), tc.rawURL)
	}
}

func TestGenerateGitlabRawURL(t *testing.T) {
	tests := []struct {
		baseURL       string
		defaultBranch string
		want          string
	}{
		{
			"https://gitlab.com/mygroup/myrepo",
			"main",
			"https://gitlab.com/mygroup/myrepo/raw/main/publiccode.yml",
		},
		{
			"https://mygitlab.example.com/org/suborg/repo",
			"develop",
			"https://mygitlab.example.com/org/suborg/repo/raw/develop/publiccode.yml",
		},
	}

	for _, tc := range tests {
		got, err := generateGitlabRawURL(tc.baseURL, tc.defaultBranch)

		require.NoError(t, err)
		assert.Equal(t, tc.want, got)
	}
}
