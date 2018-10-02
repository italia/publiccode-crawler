package jekyll

import "github.com/italia/developers-italia-backend/crawler"

// Software describe the data used in jekyll that includes publiccode and meta informations.
type Software struct {
	AmministrazioneLabel string                  `json:"it-riuso-codiceIPA-label,omitempty"`
	Id                   string                  `json:"id,omitempty"`
	Name                 string                  `json:"name,omitempty"`
	ApplicationSuite     string                  `json:"applicationSuite,omitempty"`
	URL                  string                  `json:"repoUrl,omitempty"`
	LandingURL           string                  `json:"landingURL,omitempty"`
	IsBasedOn            []string                `json:"isBasedOn,omitempty"`
	SoftwareVersion      string                  `json:"softwareVersion,omitempty"`
	ReleaseDate          string                  `json:"releaseDate,omitempty"`
	Logo                 string                  `json:"logo,omitempty"`
	MonochromeLogo       string                  `json:"monochromeLogo,omitempty"`
	Platforms            []string                `json:"platforms,omitempty"`
	Tags                 []string                `json:"tags,omitempty"`
	FreeTags             map[string][]string     `json:"freeTags,omitempty"`
	PopularTags          []string                `json:"popularTags,omitempty"`
	ShareTags            []string                `json:"shareTags,omitempty"`
	UsedBy               []string                `json:"usedBy,omitempty"`
	Roadmap              string                  `json:"roadmap,omitempty"`
	DevelopmentStatus    string                  `json:"developmentStatus,omitempty"`
	VitalityScore        float64                 `json:"vitalityScore,omitempty"`
	VitalityDataChart    []int                   `json:"vitalityDataChart,omitempty"`
	SoftwareType         SoftwareTypeData        `json:"softwareType,omitempty"`
	IntendedAudience     IntendedAudienceData    `json:"intendedAudience,omitempty"`
	Description          map[string]crawler.Desc `json:"description,omitempty"`
	OldVariant           []OldVariantData        `json:"oldVariant,omitempty"`
	OldFeatureList       map[string][]string     `json:"oldFeatureList,omitempty"`
	TagsRelate           []string                `json:"tagsRelate,omitempty"`

	RelatedSoftwares []RelatedSoftware `json:"relatedSoftwares,omitempty"`

	Legal       LegalData `json:"legal,omitempty"`
	Maintenance struct {
		Type        string               `json:"type,omitempty"`
		Contractors []crawler.Contractor `json:"contractors,omitempty"`
		Contacts    []crawler.Contact    `json:"contacts,omitempty"`
	} `json:"maintenance,omitempty"`
	Localisation LocalisationData `json:"localisation,omitempty"`
	Dependencies DependenciesData `json:"dependencies,omitempty"`
	It           ExtensionIT      `json:"it,omitempty"`
}

// LegalData contains the legal infos of the project.
type LegalData struct {
	License            string `json:"license,omitempty"`
	MainCopyrightOwner string `json:"mainCopyrightOwner,omitempty"`
	RepoOwner          string `json:"repoOwner,omitempty"`
	AuthorsFile        string `json:"authorsFile,omitempty"`
}

// LocalisationData contains available languages of the software and if it's ready to be localised.
type LocalisationData struct {
	LocalisationReady  bool     `json:"localisationReady,omitempty"`
	AvailableLanguages []string `json:"availableLanguages,omitempty"`
}

// DependenciesData contains infos about software and hardware dependencies.
type DependenciesData struct {
	Open        []crawler.Dependency `json:"open,omitempty"`
	Proprietary []crawler.Dependency `json:"proprietary,omitempty"`
	Hardware    []crawler.Dependency `json:"hardware,omitempty"`
}

// SoftwareTypeData contains the type of the software and related software urls.
type SoftwareTypeData struct {
	Type        string   `json:"type,omitempty"`
	IsRelatedTo []string `json:"isRelatedTo,omitempty"`
}

// IntendedAudienceData contains infos about actors intended to use or not the software.
type IntendedAudienceData struct {
	OnlyFor              []string `json:"onlyFor,omitempty"`
	Countries            []string `json:"countries,omitempty"`
	UnsupportedCountries []string `json:"unsupportedCountries,omitempty"`
}

// OldVariantData contains infos about other variants.
type OldVariantData struct {
	Name          string  `json:"name,omitempty"`
	VitalityScore float64 `json:"vitalityScore,omitempty"`
	Legal         struct {
		RepoOwner string `json:"repoOwner,omitempty"`
	} `json:"legal,omitempty"`
	Eng struct {
		URL      string   `json:"url,omitempty"`
		Features []string `json:"features,omitempty"`
	} `json:"eng,omitempty"`
	Ita struct {
		URL      string   `json:"url,omitempty"`
		Features []string `json:"features,omitempty"`
	} `json:"ita,omitempty"`
}

// RelatedSoftware contains info of software related to.
type RelatedSoftware struct {
	Name  string `json:"name,omitempty"`
	Image string `json:"image,omitempty"`
	Eng   struct {
		LocalisedName string `json:"localisedName,omitempty"`
		URL           string `json:"url,omitempty"`
	} `json:"eng,omitempty"`
	Ita struct {
		LocalisedName string `json:"localisedName,omitempty"`
		URL           string `json:"url,omitempty"`
	} `json:"ita,omitempty"`
}

// ExtensionIT describe the fields in italian extension.
type ExtensionIT struct {
	Accessibile    bool          `json:"accessibile,omitempty"`
	Interoperabile bool          `json:"interoperabile,omitempty"`
	Riuso          ItRiusoData   `json:"riuso,omitempty"`
	Spid           bool          `json:"spid,omitempty"`
	Pagopa         bool          `json:"pagopa,omitempty"`
	Cie            bool          `json:"cie,omitempty"`
	Anpr           bool          `json:"anpr,omitempty"`
	DesignKit      DesignKitData `json:"designKit,omitempty"`
}

// ItRiusoData is the codiceIPA of the software.
type ItRiusoData struct {
	CodiceIPA string `json:"codiceIPA,omitempty"`
}

// DesignKitData describe Design Kit variables for italian extension.
type DesignKitData struct {
	Seo bool `json:"seo,omitempty"`
	UI  bool `json:"ui,omitempty"`
	Web bool `json:"web,omitempty"`
}
