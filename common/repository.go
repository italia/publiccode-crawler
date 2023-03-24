package common

import (
	"net/url"
)

// Repository is a single code repository. FileRawURL contains the direct url to the raw file.
type Repository struct {
	Name         string
	URL          url.URL
	CanonicalURL url.URL
	FileRawURL   string
	GitBranch    string
	Publisher    Publisher
	Headers      map[string]string
}
