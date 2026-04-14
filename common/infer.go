package common

import (
	"net/url"

	"github.com/alranel/go-vcsurl/v2"
)

// InferDriver returns the catalog driver name for a URL based on its hostname.
// Returns an empty string if the platform is not recognized.
func InferDriver(u url.URL) string {
	switch {
	case vcsurl.IsGitHub(&u):
		return "github"
	case vcsurl.IsGitLab(&u):
		return "gitlab"
	case vcsurl.IsBitBucket(&u):
		return "bitbucket"
	case vcsurl.IsGitea(&u) || vcsurl.IsForgeJo(&u):
		return "gitea"
	default:
		return ""
	}
}
