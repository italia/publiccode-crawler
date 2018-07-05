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
const Version = "http://w3id.org/publiccode/version/0.1"

// Publiccode is a publiccode.yml file definition.
// Reference: https://github.com/publiccodenet/publiccode.yml
type PublicCode struct {
	PubliccodeYamlVersion string `yaml:"publiccode-yaml-version"`

	Name             string   `yaml:"name"`
	ApplicationSuite string   `yaml:"applicationSuite"`
	URL              *url.URL `yaml:"url"`
	LandingURL       *url.URL `yaml:"landingURL"`

	IsBasedOn       []string  `yaml:"isBasedOn"`
	SoftwareVersion string    `yaml:"softwareVersion"`
	ReleaseDate     time.Time `yaml:"releaseDate"`
	Logo            string    `yaml:"logo"`
	MonochromeLogo  string    `yaml:"monochromeLogo"`

	InputTypes  []string `yaml:"inputTypes"`
	OutputTypes []string `yaml:"outputTypes"`

	Platforms []string `yaml:"platforms"`

	Tags []string `yaml:"tags"`

	FreeTags map[string][]string `yaml:"freeTags"`

	UsedBy []string `yaml:"usedBy"`

	Roadmap *url.URL `yaml:"roadmap"`

	DevelopmentStatus string `yaml:"developmentStatus"`

	SoftwareType string `yaml:"softwareType"`

	IntendedAudience struct {
		OnlyFor              []string `yaml:"onlyFor"`
		Countries            []string `yaml:"countries"`
		UnsupportedCountries []string `yaml:"unsupportedCountries"`
	} `yaml:"intendedAudience"`

	Description map[string]Desc `yaml:"description"`

	Legal struct {
		License            string `yaml:"license"`
		MainCopyrightOwner string `yaml:"mainCopyrightOwner"`
		RepoOwner          string `yaml:"repoOwner"`
		AuthorsFile        string `yaml:"authorsFile"`
	} `yaml:"legal"`

	Maintenance struct {
		Type        string       `yaml:"type"`
		Contractors []Contractor `yaml:"contractors"`
		Contacts    []Contact    `yaml:"contacts"`
	} `yaml:"maintenance"`

	Localisation struct {
		LocalisationReady  bool     `yaml:"localisationReady"`
		AvailableLanguages []string `yaml:"availableLanguages"`
	} `yaml:"localisation"`

	Dependencies struct {
		Open        []Dependency `yaml:"open"`
		Proprietary []Dependency `yaml:"proprietary"`
		Hardware    []Dependency `yaml:"hardware"`
	} `yaml:"dependencies"`

	It ExtensionIT `yaml:"it"`
}

// Desc is a general description of the software.
// Reference: https://github.com/publiccodenet/publiccode.yml/blob/develop/schema.md#section-description
type Desc struct {
	LocalisedName    string     `yaml:"localisedName"`
	GenericName      string     `yaml:"genericName"`
	ShortDescription string     `yaml:"shortDescription"`
	LongDescription  string     `yaml:"longDescription"`
	Documentation    *url.URL   `yaml:"documentation"`
	APIDocumentation *url.URL   `yaml:"apiDocumentation"`
	FeatureList      []string   `yaml:"featureList"`
	Screenshots      []string   `yaml:"screenshots"`
	Videos           []*url.URL `yaml:"videos"`
	Awards           []string   `yaml:"awards"`
}

// Contractor is an entity or entities, if any, that are currently contracted for maintaining the software.
// Reference: https://github.com/publiccodenet/publiccode.yml/blob/develop/schema.md#contractor
type Contractor struct {
	Name    string    `yaml:"name"`
	Website *url.URL  `yaml:"website"`
	Until   time.Time `yaml:"until"`
}

// Contact is a contact info maintaining the software.
// Reference: https://github.com/publiccodenet/publiccode.yml/blob/develop/schema.md#contact
type Contact struct {
	Name        string `yaml:"name"`
	Email       string `yaml:"email"`
	Affiliation string `yaml:"affiliation"`
	Phone       string `yaml:"phone"`
}

// Dependency describe system-level dependencies required to install and use this software.
// Reference: https://github.com/publiccodenet/publiccode.yml/blob/develop/schema.md#section-dependencies
type Dependency struct {
	Name       string `yaml:"name"`
	VersionMin string `yaml:"versionMin"`
	VersionMax string `yaml:"versionMax"`
	Optional   bool   `yaml:"optional"`
	Version    string `yaml:"version"`
}
