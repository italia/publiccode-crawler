package jekyll

import "github.com/italia/developers-italia-backend/crawler"

// PublicCode describe the data used in jekyll that includes publiccode and meta informations.
type Software struct {
	Name              string                  `json:"name,omitempty"`
	ApplicationSuite  string                  `json:"applicationSuite,omitempty"`
	URL               string                  `json:"url,omitempty"`
	LandingURL        string                  `json:"landingURL,omitempty"`
	IsBasedOn         []string                `json:"isBasedOn,omitempty"`
	SoftwareVersion   string                  `json:"softwareVersion,omitempty"`
	ReleaseDate       string                  `json:"releaseDate,omitempty"`
	Logo              string                  `json:"logo,omitempty"`
	MonochromeLogo    string                  `json:"monochromeLogo,omitempty"`
	Platforms         []string                `json:"platforms,omitempty"`
	Tags              []string                `json:"tags,omitempty"`
	FreeTags          FreeTagsData            `json:"freeTags,omitempty"`
	PopularTags       []string                `json:"popularTags,omitempty"`
	ShareTags         []string                `json:"shareTags,omitempty"`
	UsedBy            []string                `json:"usedBy,omitempty"`
	Roadmap           string                  `json:"roadmap,omitempty"`
	DevelopmentStatus string                  `json:"developmentStatus,omitempty"`
	VitalityScore     float64                 `json:"vitalityScore,omitempty"`
	VitalityDataChart []int                   `json:"vitalityDataChart,omitempty"`
	SoftwareType      SoftwareTypeData        `json:"softwareType,omitempty"`
	IntendedAudience  IntendedAudienceData    `json:"intendedAudience,omitempty"`
	Description       map[string]crawler.Desc `json:"description,omitempty"`
	OldVariant        []OldVariantData        `json:"oldVariant,omitempty"`
	OldFeatureList    OldFeatureListData      `json:"oldFeatureList,omitempty"`
	TagsRelate        []string                `json:"tagsRelate,omitempty"`

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

type LegalData struct {
	License            string `json:"license,omitempty"`
	MainCopyrightOwner string `json:"mainCopyrightOwner,omitempty"`
	RepoOwner          string `json:"repoOwner,omitempty"`
	AuthorsFile        string `json:"authorsFile,omitempty"`
}

type LocalisationData struct {
	LocalisationReady  bool     `json:"localisationReady,omitempty"`
	AvailableLanguages []string `json:"availableLanguages,omitempty"`
}

type DependenciesData struct {
	Open        []crawler.Dependency `json:"open,omitempty"`
	Proprietary []crawler.Dependency `json:"proprietary,omitempty"`
	Hardware    []crawler.Dependency `json:"hardware,omitempty"`
}

type FreeTagsData struct {
	Eng []string `json:"eng,omitempty"`
	Ita []string `json:"ita,omitempty"`
}

type SoftwareTypeData struct {
	Type        string   `json:"type,omitempty"`
	IsRelatedTo []string `json:"isRelatedTo,omitempty"`
}

type IntendedAudienceData struct {
	OnlyFor              []string `json:"onlyFor,omitempty"`
	Countries            []string `json:"countries,omitempty"`
	UnsupportedCountries []string `json:"unsupportedCountries,omitempty"`
}

type OldFeatureListData struct {
	Eng []string `json:"eng,omitempty"`
	Ita []string `json:"ita,omitempty"`
}

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

type ItRiusoData struct {
	CodiceIPA string `json:"codiceIPA,omitempty"`
}

type DesignKitData struct {
	Seo bool `json:"seo,omitempty"`
	UI  bool `json:"ui,omitempty"`
	Web bool `json:"web,omitempty"`
}
