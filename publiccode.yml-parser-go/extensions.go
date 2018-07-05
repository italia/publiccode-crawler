package publiccode

// Country-specific extensions
//
// While the standard is structured to be meaningful on an international level,
// there are additional information that can be added that makes sense in specific
// countries, such as declaring compliance with local laws or regulations.
// The provided extension mechanism is the usage of country-specific sections.
//
// All country-specific sections are contained in a section named with
// the two-letter lowercase ISO 3166-1 alpha-2 country code.
//
// For instance "spid" is a property for Italian software declaring whether
// the software is integrated with the Italian Public Identification System.
// Reference: https://github.com/publiccodenet/publiccode.yml/blob/develop/schema.md

// ExtensionIT is the country-specific extension for Italy.
// Reference: https://github.com/publiccodenet/publiccode.yml/blob/develop/schema.it.md
type ExtensionIT struct {
	Conforme struct {
		Accessibile    bool `yaml:"accessibile"`
		Interoperabile bool `yaml:"interoperabile"`
		Sicuro         bool `yaml:"sicuro"`
		Privacy        bool `yaml:"privacy"`
	} `yaml:"conforme"`

	Riuso struct {
		CodiceIPA string `yaml:"codiceIPA"`
	} `yaml:"riuso"`

	Spid   bool `yaml:"spid"`
	Pagopa bool `yaml:"pagopa"`
	Cie    bool `yaml:"cie"`
	Anpr   bool `yaml:"anpr"`

	Ecosistemi []string `yaml:"ecosistemi"`

	DesignKit struct {
		Seo     bool `yaml:"seo"`
		UI      bool `yaml:"ui"`
		Web     bool `yaml:"web"`
		Content bool `yaml:"content"`
	} `yaml:"designKit"`
}
