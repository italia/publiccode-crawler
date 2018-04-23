package publiccode

import (
	"net/url"
	"time"
)

// BaseDir is the Base Directory of the PublicCode file.
// If local will be empty: ""
// If remote will be the url of the repository
var BaseDir = ""

// Version of the PublicCode specs.
// Source https://github.com/publiccodenet/publiccode.yml
const Version = "0.1"

// A PublicCode represents a standard metadata description for public software and policy repositories.
type PublicCode struct {
	Version     string
	Url         *url.URL
	UpstreamUrl []*url.URL

	Legal struct {
		License            string
		MainCopyrightOwner string
		AuthorsFile        string
		RepoOwner          string
	}

	Maintenance struct {
		Type              string
		Until             time.Time
		Maintainer        []string
		TechnicalContacts []Contact
	}

	Description struct {
		Name        string
		Logo        []string
		Shortdesc   []Desc
		LongDesc    []Desc
		Screenshots []string
		Videos      []*url.URL
		Version     string
		Released    time.Time
		Platforms   string
	}

	Meta struct {
		Scope    []string
		PaType   []string
		Category string
		Tags     []string
		UsedBy   []string
	}

	Dependencies struct {
		Hardware    []string
		Open        []string
		Proprietary []string
	}
}

// A Contact represents all the standard informations about a technical contact.
type Contact struct {
	Name        string
	Email       string
	Affiliation string
}

// A Desc represents a generic description with multiple languages.
type Desc struct {
	En string
	It string
}
