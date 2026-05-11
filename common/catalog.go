package common

import "net/url"

// CatalogSource is one of a catalog's enumeration points. By definition a
// source produces a list of repositories, so there is no Group flag.
// Driver may be a code-host driver ("github", "gitlab", "bitbucket",
// "gitea"), an enumeration driver ("json"), or an upstream API
// ("software-catalog-api").
type CatalogSource struct {
	URL    url.URL
	Driver string
	Args   []string
}

type Catalog struct {
	ID                  string
	Name                string
	PublishersNamespace string
	Sources             []CatalogSource
}
