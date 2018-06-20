package main

type PublicCode struct {
	PubliccodeYamlVersion string `json:"publiccode-yaml-version"`

	Name             string `json:"name"`
	ApplicationSuite string `json:"application-suite"`
	URL              string `json:"url"`
	LandingURL       string `json:"landing-url"`

	IsBasedOn       []string `json:"is-based-on"`
	SoftwareVersion string   `json:"software-version"`
	ReleaseDate     string   `json:"release-date"`
	Logo            string   `json:"logo"`
	MonochromeLogo  string   `json:"monochrome-logo"`

	InputTypes  []string `json:"input-types"`
	OutputTypes []string `json:"output-types"`

	Platforms []string `json:"platforms"`

	Tags []string `json:"tags"`

	FreeTags map[string][]string `json:"free-tags"`

	UsedBy []string `json:"used-by"`

	Roadmap string `json:"roadmap"`

	DevelopmentStatus string `json:"development-status"`

	// Vitalityscore
	VitalityScore     int   `json:"vitalityScore"`
	VitalityDataChart []int `json:"vitalityDataChart"`

	RelatedSoftware string //TODO

	SoftwareType string `json:"software-type"`

	IntendedAudience struct {
		OnlyFor              []string `json:"intended-audience-only-for"`
		Countries            []string `json:"intended-audience-countries"`
		UnsupportedCountries []string `json:"intended-audience-unsupported-countries"`
	} `json:"intendedAudience"`

	Description map[string]Desc `json:"description"`
	OldVariants []OldVariant    `json:"old-variant"`

	Legal struct {
		License            string `json:"legal-license"`
		MainCopyrightOwner string `json:"legal-main-copyright-owner"`
		RepoOwner          string `json:"legal-repo-owner"`
		AuthorsFile        string `json:"legal-authors-file"`
	} `json:"legal"`

	Maintenance struct {
		Type        string       `json:"maintainance-type"`
		Contractors []Contractor `json:"maintainance-contractors"`
		Contacts    []Contact    `json:"maintainance-contacts"`
	} `json:"maintenance"`

	Localisation struct {
		LocalisationReady  bool     `json:"localisation-localisation-ready"`
		AvailableLanguages []string `json:"localisation-available-languages"`
	} `json:"localisation"`

	Dependencies struct {
		Open        []Dependency `json:"dependencies-open"`
		Proprietary []Dependency `json:"dependencies-proprietary"`
		Hardware    []Dependency `json:"dependencies-hardware"`
	} `json:"dependencies"`
}

// Desc is a general description of the software.
// Reference: https://github.com/publiccodenet/publiccode.yml/blob/develop/schema.md#section-description
type Desc struct {
	LocalisedName    string   `json:"localisedName"`
	GenericName      string   `json:"genericName"`
	ShortDescription string   `json:"shortDescription"`
	LongDescription  string   `json:"longDescription"`
	Documentation    string   `json:"documentation"`
	APIDocumentation string   `json:"apiDocumentation"`
	FeatureList      []string `json:"featureList"`
	Screenshots      []string `json:"screenshots"`
	Videos           []string `json:"videos"`
	Awards           []string `json:"awards"`
}

// Contractor is an entity or entities, if any, that are currently contracted for maintaining the software.
// Reference: https://github.com/publiccodenet/publiccode.yml/blob/develop/schema.md#contractor
type Contractor struct {
	Name    string `json:"name"`
	Website string `json:"website"`
	Until   string `json:"until"`
}

// Contact is a contact info maintaining the software.
// Reference: https://github.com/publiccodenet/publiccode.yml/blob/develop/schema.md#contact
type Contact struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Affiliation string `json:"affiliation"`
	Phone       string `json:"phone"`
}

// Dependency describe system-level dependencies required to install and use this software.
// Reference: https://github.com/publiccodenet/publiccode.yml/blob/develop/schema.md#section-dependencies
type Dependency struct {
	Name       string `json:"name"`
	VersionMin string `json:"versionMin"`
	VersionMax string `json:"versionMax"`
	Optional   bool   `json:"optional"`
	Version    string `json:"version"`
}

type OldVariant struct {
	Name        string             `json:"name"`
	URL         string             `json:"url"`
	Description map[string]OldDesc `json:"description"`
}
type OldDesc struct {
	LocalisedName string   `json:"localisedName"`
	GenericName   string   `json:"genericName"`
	FeatureList   []string `json:"featureList"`
}
