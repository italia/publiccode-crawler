package crawler

// SoftwareES represents a software record in Elasticsearch
type SoftwareES struct {
	FileRawURL            string       `json:"fileRawURL"`
	ID                    string       `json:"id"`
	CrawlTime             string       `json:"crawltime"`
	ItRiusoCodiceIPALabel string       `json:"it-riuso-codiceIPA-label"`
	PublicCode            PublicCodeES `json:"publiccode"`
	VitalityScore         float64      `json:"vitalityScore"`
	VitalityDataChart     []int        `json:"vitalityDataChart"`
}

// PublicCodeES describe the data in ElasticSearch that includes publiccode.
type PublicCodeES struct {
	Name             string `json:"name"`
	ApplicationSuite string `json:"applicationSuite,omitempty"`
	URL              string `json:"url"`
	LandingURL       string `json:"landingURL"`

	IsBasedOn       []string `json:"isBasedOn"`
	SoftwareVersion string   `json:"softwareVersion"`
	ReleaseDate     string   `json:"releaseDate"`
	Logo            string   `json:"logo"`
	MonochromeLogo  string   `json:"monochromeLogo"`

	InputTypes  []string `json:"inputTypes"`
	OutputTypes []string `json:"outputTypes"`

	Platforms []string `json:"platforms"`

	Tags []string `json:"tags"`

	UsedBy []string `json:"used-by"`

	Roadmap string `json:"roadmap"`

	DevelopmentStatus string `json:"development-status"`

	SoftwareType string `json:"software-type"`

	IntendedAudienceOnlyFor              []string `json:"intended-audience-only-for"`
	IntendedAudienceCountries            []string `json:"intended-audience-countries"`
	IntendedAudienceUnsupportedCountries []string `json:"intended-audience-unsupported-countries"`

	Description map[string]Desc `json:"description"`

	LegalLicense            string `json:"legal-license"`
	LegalMainCopyrightOwner string `json:"legal-main-copyright-owner"`
	LegalRepoOwner          string `json:"legal-repo-owner"`
	LegalAuthorsFile        string `json:"legal-authors-file"`

	MaintenanceType        string       `json:"maintenance-type"`
	MaintenanceContractors []Contractor `json:"maintenance-contractors"`
	MaintenanceContacts    []Contact    `json:"maintenance-contacts"`

	LocalisationLocalisationReady  bool     `json:"localisation-localisation-ready"`
	LocalisationAvailableLanguages []string `json:"localisation-available-languages"`

	DependsOnOpen        []Dependency `json:"dependencies-open"`
	DependsOnProprietary []Dependency `json:"dependencies-proprietary"`
	DependsOnHardware    []Dependency `json:"dependencies-hardware"`

	// Italian extension.
	ItConformeAccessibile    bool `json:"it-conforme-accessibile"`
	ItConformeInteroperabile bool `json:"it-conforme-interoperabile"`
	ItConformeSicuro         bool `json:"it-conforme-sicuro"`
	ItConformePrivacy        bool `json:"it-conforme-privacy"`

	ItRiusoCodiceIPA string `json:"it-riuso-codice-ipa"`

	ItSpid   bool `json:"it-spid"`
	ItPagopa bool `json:"it-pagopa"`
	ItCie    bool `json:"it-cie"`
	ItAnpr   bool `json:"it-anpr"`

	ItEcosistemi []string `json:"it-ecosistemi"`

	ItDesignKitSeo     bool `json:"it-design-kit-seo"`
	ItDesignKitUI      bool `json:"it-design-kit-ui"`
	ItDesignKitWeb     bool `json:"it-design-kit-web"`
	ItDesignKitContent bool `json:"it-design-kit-content"`
}

// Desc is a general description of the software.
// Reference: http://w3id.org/publiccode/version/0.1/schema.html#section-description
type Desc struct {
	LocalisedName    string   `json:"localisedName"`
	GenericName      string   `json:"genericName"`
	ShortDescription string   `json:"shortDescription"`
	LongDescription  string   `json:"longDescription"`
	Documentation    string   `json:"documentation"`
	APIDocumentation string   `json:"apiDocumentation"`
	Features         []string `json:"features"`
	Screenshots      []string `json:"screenshots"`
	Videos           []string `json:"videos"`
	Awards           []string `json:"awards"`
	FreeTags         []string `json:"freeTags"`
}

// Contractor is an entity or entities, if any, that are currently contracted for maintaining the software.
// Reference: http://w3id.org/publiccode/version/0.1/schema.html#contractor
type Contractor struct {
	Name    string `json:"name"`
	Website string `json:"website"`
	Until   string `json:"until"`
}

// Contact is a contact info maintaining the software.
// Reference: http://w3id.org/publiccode/version/0.1/schema.html#contact
type Contact struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Affiliation string `json:"affiliation"`
	Phone       string `json:"phone"`
}

// Dependency describe system-level dependencies required to install and use this software.
// Reference: http://w3id.org/publiccode/version/0.1/schema.html#section-dependencies
type Dependency struct {
	Name       string `json:"name"`
	VersionMin string `json:"versionMin"`
	VersionMax string `json:"versionMax"`
	Optional   bool   `json:"optional"`
	Version    string `json:"version"`
}

// OldVariant describe some infos about variant of the software.
type OldVariant struct {
	Name        string             `json:"name"`
	URL         string             `json:"url"`
	Description map[string]OldDesc `json:"description"`
}

// OldDesc is a description of old OldVariants softwares.
type OldDesc struct {
	LocalisedName string   `json:"localisedName"`
	GenericName   string   `json:"genericName"`
	Features      []string `json:"features"`
}

//Audience describe who is the audience of this software.
type Audience struct {
	OnlyFor              []string `json:"intended-audience-only-for"`
	Countries            []string `json:"intended-audience-countries"`
	UnsupportedCountries []string `json:"intended-audience-unsupported-countries"`
}
