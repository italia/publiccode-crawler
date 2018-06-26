package crawler

import (
	"net/url"
	"time"
)

// Publiccode is a publiccode.yml file definition.
// Reference: https://github.com/publiccodenet/publiccode.yml
type PublicCode struct {
	PubliccodeYamlVersion string `yaml:"publiccode-yaml-version" json:"publiccode-yaml-version"`

	Name             string   `yaml:"name" json:"name"`
	ApplicationSuite string   `yaml:"applicationSuite" json:"applicationSuite"`
	URL              *url.URL `yaml:"url" json:"url"`
	LandingURL       *url.URL `yaml:"landingURL" json:"landingURL"`

	IsBasedOn       []string  `yaml:"isBasedOn" json:"isBasedOn"`
	SoftwareVersion string    `yaml:"softwareVersion" json:"softwareVersion"`
	ReleaseDate     time.Time `yaml:"releaseDate" json:"releaseDate"`
	Logo            string    `yaml:"logo" json:"logo"`
	MonochromeLogo  string    `yaml:"monochromeLogo" json:"monochromeLogo"`

	InputTypes  []string `yaml:"inputTypes" json:"inputTypes"`
	OutputTypes []string `yaml:"outputTypes" json:"outputTypes"`

	Platforms []string `yaml:"platforms" json:"platforms"`

	Tags []string `yaml:"tags" json:"tags"`

	FreeTags map[string][]string `yaml:"freeTags" json:"freeTags"`

	UsedBy []string `yaml:"usedBy" json:"usedBy"`

	Roadmap *url.URL `yaml:"roadmap" json:"roadmap"`

	DevelopmentStatus string `yaml:"developmentStatus" json:"developmentStatus"`

	SoftwareType string `yaml:"softwareType" json:"softwareType"`

	IntendedAudience struct {
		OnlyFor              []string `yaml:"onlyFor" json:"onlyFor"`
		Countries            []string `yaml:"countries" json:"countries"`
		UnsupportedCountries []string `yaml:"unsupportedCountries" json:"unsupportedCountries"`
	} `yaml:"intendedAudience" json:"intendedAudience"`

	Description map[string]Desc `yaml:"description" json:"description"`

	Legal struct {
		License            string `yaml:"license" json:"license"`
		MainCopyrightOwner string `yaml:"mainCopyrightOwner" json:"mainCopyrightOwner"`
		RepoOwner          string `yaml:"repoOwner" json:"repoOwner"`
		AuthorsFile        string `yaml:"authorsFile" json:"authorsFile"`
	} `yaml:"legal" json:"legal"`

	Maintenance struct {
		Type        string       `yaml:"type" json:"type"`
		Contractors []Contractor `yaml:"contractors" json:"contractors"`
		Contacts    []Contact    `yaml:"contacts" json:"contacts"`
	} `yaml:"maintenance" json:"maintenance"`

	Localisation struct {
		LocalisationReady  bool     `yaml:"localisationReady" json:"localisationReady"`
		AvailableLanguages []string `yaml:"availableLanguages" json:"availableLanguages"`
	} `yaml:"localisation" json:"localisation"`

	Dependencies struct {
		Open        []Dependency `yaml:"open" json:"open"`
		Proprietary []Dependency `yaml:"proprietary" json:"proprietary"`
		Hardware    []Dependency `yaml:"hardware" json:"hardware"`
	} `yaml:"dependencies" json:"dependencies"`

	It ExtensionIT `yaml:"it" json:"it"`
}

// Desc is a general description of the software.
// Reference: https://github.com/publiccodenet/publiccode.yml/blob/develop/schema.md#section-description
type Desc struct {
	LocalisedName    string     `yaml:"localisedName" json:"localisedName"`
	GenericName      string     `yaml:"genericName" json:"genericName"`
	ShortDescription string     `yaml:"shortDescription" json:"shortDescription"`
	LongDescription  string     `yaml:"longDescription" json:"longDescription"`
	Documentation    *url.URL   `yaml:"documentation" json:"documentation"`
	APIDocumentation *url.URL   `yaml:"apiDocumentation" json:"apiDocumentation"`
	FeatureList      []string   `yaml:"featureList" json:"featureList"`
	Screenshots      []string   `yaml:"screenshots" json:"screenshots"`
	Videos           []*url.URL `yaml:"videos" json:"videos"`
	Awards           []string   `yaml:"awards" json:"awards"`
}

// Contractor is an entity or entities, if any, that are currently contracted for maintaining the software.
// Reference: https://github.com/publiccodenet/publiccode.yml/blob/develop/schema.md#contractor
type Contractor struct {
	Name    string    `yaml:"name" json:"name"`
	Website *url.URL  `yaml:"website" json:"website"`
	Until   time.Time `yaml:"until" json:"until"`
}

// Contact is a contact info maintaining the software.
// Reference: https://github.com/publiccodenet/publiccode.yml/blob/develop/schema.md#contact
type Contact struct {
	Name        string `yaml:"name" json:"name"`
	Email       string `yaml:"email" json:"email"`
	Affiliation string `yaml:"affiliation" json:"affiliation"`
	Phone       string `yaml:"phone" json:"phone"`
}

// Dependency describe system-level dependencies required to install and use this software.
// Reference: https://github.com/publiccodenet/publiccode.yml/blob/develop/schema.md#section-dependencies
type Dependency struct {
	Name       string `yaml:"name" json:"name"`
	VersionMin string `yaml:"versionMin" json:"versionMin"`
	VersionMax string `yaml:"versionMax" json:"versionMax"`
	Optional   bool   `yaml:"optional" json:"optional"`
	Version    string `yaml:"version" json:"optional"`
}

// ExtensionIT is the country-specific extension for Italy.
// Reference: https://github.com/publiccodenet/publiccode.yml/blob/develop/schema.it.md
type ExtensionIT struct {
	Conforme struct {
		Accessibile    bool `yaml:"accessibile" json:"accessibile"`
		Interoperabile bool `yaml:"interoperabile" json:"interoperabile"`
		Sicuro         bool `yaml:"sicuro" json:"sicuro"`
		Privacy        bool `yaml:"privacy" json:"privacy"`
	} `yaml:"conforme" json:"conforme"`

	Riuso struct {
		CodiceIPA string `yaml:"codiceIPA" json:"codiceIPA"`
	} `yaml:"riuso" json:"riuso"`

	Spid   bool `yaml:"spid" json:"spid"`
	Pagopa bool `yaml:"pagopa" json:"pagopa"`
	Cie    bool `yaml:"cie" json:"cie"`
	Anpr   bool `yaml:"anpr" json:"anpr"`

	Ecosistemi []string `yaml:"ecosistemi" json:"ecosistemi"`

	DesignKit struct {
		Seo     bool `yaml:"seo" json:"seo"`
		UI      bool `yaml:"ui" json:"ui"`
		Web     bool `yaml:"web" json:"web"`
		Content bool `yaml:"content" json:"content"`
	} `yaml:"designKit" json:"designKit"`
}
