package jekyll

import "github.com/italia/developers-italia-backend/crawler"

// PublicCode describe the data used in jekyll that includes publiccode and meta informations.
type Software struct {
	Name              string                  `json:"name"`
	ApplicationSuite  string                  `json:"applicationSuite"`
	URL               string                  `json:"url"`
	LandingURL        string                  `json:"landingURL"`
	IsBasedOn         []string                `json:"isBasedOn"`
	SoftwareVersion   string                  `json:"softwareVersion"`
	ReleaseDate       string                  `json:"releaseDate"`
	Logo              string                  `json:"logo"`
	MonochromeLogo    string                  `json:"monochromeLogo"`
	Platforms         []string                `json:"platforms"`
	Tags              []string                `json:"tags"`
	FreeTags          FreeTagsData            `json:"freeTags"`
	PopularTags       []string                `json:"popularTags"`
	ShareTags         []string                `json:"shareTags"`
	UsedBy            []string                `json:"usedBy"`
	Roadmap           string                  `json:"roadmap"`
	DevelopmentStatus string                  `json:"developmentStatus"`
	VitalityScore     float64                 `json:"vitalityScore"`
	VitalityDataChart []int                   `json:"vitalityDataChart"`
	SoftwareType      SoftwareTypeData        `json:"softwareType"`
	IntendedAudience  IntendedAudienceData    `json:"intendedAudience"`
	Description       map[string]crawler.Desc `json:"description"`
	OldVariant        []OldVariantData        `json:"oldVariant"`
	OldFeatureList    OldFeatureListData      `json:"oldFeatureList"`
	TagsRelate        []string                `json:"tagsRelate"`

	RelatedSoftwares []RelatedSoftware `json:"relatedSoftwares"`

	Legal       LegalData `json:"legal"`
	Maintenance struct {
		Type        string               `json:"type"`
		Contractors []crawler.Contractor `json:"contractors"`
		Contacts    []crawler.Contact    `json:"contacts"`
	} `json:"maintenance"`
	Localisation LocalisationData `json:"localisation"`
	Dependencies DependenciesData `json:"dependencies"`
	It           ExtensionIT      `json:"it"`
}

type LegalData struct {
	License            string `json:"license"`
	MainCopyrightOwner string `json:"mainCopyrightOwner"`
	RepoOwner          string `json:"repoOwner"`
	AuthorsFile        string `json:"authorsFile"`
}

type LocalisationData struct {
	LocalisationReady  bool     `json:"localisationReady"`
	AvailableLanguages []string `json:"availableLanguages"`
}

type DependenciesData struct {
	Open        []crawler.Dependency `json:"open"`
	Proprietary []crawler.Dependency `json:"proprietary"`
	Hardware    []crawler.Dependency `json:"hardware"`
}

type FreeTagsData struct {
	Eng []string `json:"eng"`
	Ita []string `json:"ita"`
}

type SoftwareTypeData struct {
	Type        string   `json:"type"`
	IsRelatedTo []string `json:"isRelatedTo"`
}

type IntendedAudienceData struct {
	OnlyFor              []string `json:"onlyFor"`
	Countries            []string `json:"countries"`
	UnsupportedCountries []string `json:"unsupportedCountries"`
}

type OldFeatureListData struct {
	Eng []string `json:"eng"`
	Ita []string `json:"ita"`
}

type OldVariantData struct {
	Name          string  `json:"name"`
	VitalityScore float64 `json:"vitalityScore"`
	Legal         struct {
		RepoOwner string `json:"repoOwner"`
	} `json:"legal"`
	Eng struct {
		URL      string   `json:"url"`
		Features []string `json:"features"`
	} `json:"eng"`
	Ita struct {
		URL      string   `json:"url"`
		Features []string `json:"features"`
	} `json:"ita"`
}

type RelatedSoftware struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Eng   struct {
		LocalisedName string `json:"localisedName"`
		URL           string `json:"url"`
	} `json:"eng"`
	Ita struct {
		LocalisedName string `json:"localisedName"`
		URL           string `json:"url"`
	} `json:"ita"`
}

type ExtensionIT struct {
	Accessibile    bool          `json:"accessibile"`
	Interoperabile bool          `json:"interoperabile"`
	Riuso          ItRiusoData   `json:"riuso"`
	Spid           bool          `json:"spid"`
	Pagopa         bool          `json:"pagopa"`
	Cie            bool          `json:"cie"`
	Anpr           bool          `json:"anpr"`
	DesignKit      DesignKitData `json:"designKit"`
}

type ItRiusoData struct {
	CodiceIPA string `json:"codiceIPA"`
}

type DesignKitData struct {
	Seo bool `json:"seo"`
	UI  bool `json:"ui"`
	Web bool `json:"web"`
}
