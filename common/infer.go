package common

import (
	"net/url"

	"github.com/alranel/go-vcsurl/v2"
)

// InferVCSDriver returns the VCS driver name for a URL based on its hostname.
// Returns an empty string if the platform is not recognized.
//
// The returned name matches the values that VCS sources use in
// CatalogSource.Driver ("github", "gitlab", "bitbucket", "gitea"). Non-VCS
// drivers like "json" are never inferred and must be set explicitly.
func InferVCSDriver(repoURL url.URL) string {
	switch {
	case vcsurl.IsGitHub(&repoURL):
		return "github"
	case vcsurl.IsGitLab(&repoURL):
		return "gitlab"
	case vcsurl.IsBitBucket(&repoURL):
		return "bitbucket"
	case vcsurl.IsGitea(&repoURL) || vcsurl.IsForgeJo(&repoURL):
		return "gitea"
	default:
		return ""
	}
}
