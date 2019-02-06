package crawler

import (
	"bytes"
	"context"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/italia/developers-italia-backend/crawler/ipa"
	"github.com/italia/developers-italia-backend/crawler/metrics"
	"github.com/dyatlov/go-oembed/oembed"
	"github.com/olivere/elastic"
	pcode "github.com/italia/publiccode-parser-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type administration struct {
	Name      string `json:"it-riuso-codiceIPA-label"`
	CodiceIPA string `json:"it-riuso-codiceIPA"`
}

// SaveToES save the chosen data []byte in elasticsearch
func SaveToES(fileRawURL, hashedRepoURL string, name string, activityIndex float64, vitality []int, data []byte, index string, elasticClient *elastic.Client) error {
	// Starting with elastic.v5, you must pass a context to execute each service.
	ctx := context.Background()

	// Generate publiccode data using the parser.
	pc := pcode.PublicCode{}
	err := pcode.Parse(data, &pc)
	if err != nil {
		log.Errorf("Error parsing publiccode.yml for %s.", name)
	}

	// Extract raw base url.
	rawBaseDir := strings.TrimRight(fileRawURL, viper.GetString("CRAWLED_FILENAME"))

	// Add a document to the index.
	file := PublicCodeES{
		FileRawURL:            fileRawURL,
		ID:                    hashedRepoURL,
		CrawlTime:             time.Now().String(),
		ItRiusoCodiceIPALabel: ipa.GetAdministrationName(pc.It.Riuso.CodiceIPA),

		Name:             pc.Name,
		ApplicationSuite: pc.ApplicationSuite,
		URL:              pc.URL.String(),
		LandingURL:       pc.LandingURL.String(),

		IsBasedOn:       pc.IsBasedOn,
		SoftwareVersion: pc.SoftwareVersion,
		ReleaseDate:     pc.ReleaseDate.Format("2006-01-02"),
		Logo:            concatenateLink(rawBaseDir, pc.Logo),
		MonochromeLogo:  concatenateLink(rawBaseDir, pc.MonochromeLogo),
		InputTypes:      pc.InputTypes,
		OutputTypes:     pc.OutputTypes,

		Platforms: pc.Platforms,

		Tags: pc.Tags,

		UsedBy: pc.UsedBy,

		Roadmap: pc.Roadmap.String(),

		DevelopmentStatus: pc.DevelopmentStatus,

		VitalityScore:     activityIndex,
		VitalityDataChart: vitality,

		RelatedSoftware: nil,

		SoftwareType: pc.SoftwareType,

		IntendedAudienceOnlyFor:              pc.IntendedAudience.OnlyFor,
		IntendedAudienceCountries:            pc.IntendedAudience.Countries,
		IntendedAudienceUnsupportedCountries: pc.IntendedAudience.UnsupportedCountries,

		Description: map[string]Desc{},
		OldVariants: []OldVariant{},

		LegalLicense:            pc.Legal.License,
		LegalMainCopyrightOwner: pc.Legal.MainCopyrightOwner,
		LegalRepoOwner:          pc.Legal.RepoOwner,
		LegalAuthorsFile:        concatenateLink(rawBaseDir, pc.Legal.AuthorsFile),

		MaintenanceType:        pc.Maintenance.Type,
		MaintenanceContractors: []Contractor{},
		MaintenanceContacts:    []Contact{},

		LocalisationLocalisationReady:  pc.Localisation.LocalisationReady,
		LocalisationAvailableLanguages: pc.Localisation.AvailableLanguages,

		DependsOnOpen:        []Dependency{},
		DependsOnProprietary: []Dependency{},
		DependsOnHardware:    []Dependency{},

		ItConformeAccessibile:    pc.It.Conforme.Accessibile,
		ItConformeInteroperabile: pc.It.Conforme.Interoperabile,
		ItConformeSicuro:         pc.It.Conforme.Sicuro,
		ItConformePrivacy:        pc.It.Conforme.Privacy,

		ItRiusoCodiceIPA: pc.It.Riuso.CodiceIPA,

		ItSpid:   pc.It.Spid,
		ItPagopa: pc.It.Pagopa,
		ItCie:    pc.It.Cie,
		ItAnpr:   pc.It.Anpr,

		ItEcosistemi: pc.It.Ecosistemi,

		ItDesignKitSeo:     pc.It.DesignKit.Seo,
		ItDesignKitUI:      pc.It.DesignKit.UI,
		ItDesignKitWeb:     pc.It.DesignKit.Web,
		ItDesignKitContent: pc.It.DesignKit.Content,
	}
	for _, contractor := range pc.Maintenance.Contractors {
		file.MaintenanceContractors = append(file.MaintenanceContractors, Contractor{
			Name:    contractor.Name,
			Website: contractor.Website.String(),
			Until:   contractor.Until.Format("2006-01-02"),
		})
	}
	for _, contact := range pc.Maintenance.Contacts {
		file.MaintenanceContacts = append(file.MaintenanceContacts, Contact{
			Name:        contact.Name,
			Email:       contact.Email,
			Affiliation: contact.Affiliation,
			Phone:       contact.Phone,
		})
	}
	for lang := range pc.Description {
		file.Description[lang] = Desc{
			LocalisedName:    pc.Description[lang].LocalisedName,
			GenericName:      pc.Description[lang].GenericName,
			ShortDescription: pc.Description[lang].ShortDescription,
			LongDescription:  pc.Description[lang].LongDescription,
			Documentation:    pc.Description[lang].Documentation.String(),
			APIDocumentation: pc.Description[lang].APIDocumentation.String(),
			Features:         pc.Description[lang].Features,
			Screenshots:      concatenateLinks(rawBaseDir, pc.Description[lang].Screenshots),
			Videos: func(videos []*url.URL) []string {
				var v []string
				for _, video := range videos {
					v = append(v, getOembedInfo("video", video.String()))
				}
				return v
			}(pc.Description[lang].Videos),
			Awards:   pc.Description[lang].Awards,
			FreeTags: pc.Description[lang].FreeTags,
		}

	}
	for _, dependency := range pc.DependsOn.Open {
		file.DependsOnOpen = append(file.DependsOnOpen, Dependency{
			Name:       dependency.Name,
			VersionMin: dependency.VersionMin,
			VersionMax: dependency.VersionMax,
			Optional:   dependency.Optional,
			Version:    dependency.Version,
		})
	}
	for _, dependency := range pc.DependsOn.Proprietary {
		file.DependsOnProprietary = append(file.DependsOnProprietary, Dependency{
			Name:       dependency.Name,
			VersionMin: dependency.VersionMin,
			VersionMax: dependency.VersionMax,
			Optional:   dependency.Optional,
			Version:    dependency.Version,
		})
	}
	for _, dependency := range pc.DependsOn.Hardware {
		file.DependsOnHardware = append(file.DependsOnHardware, Dependency{
			Name:       dependency.Name,
			VersionMin: dependency.VersionMin,
			VersionMax: dependency.VersionMax,
			Optional:   dependency.Optional,
			Version:    dependency.Version,
		})
	}
	// Put publiccode data in ES.
	_, err = elasticClient.Index().
		Index(index).
		Type("software").
		Id(hashedRepoURL).
		BodyJson(file).
		Do(ctx)
	if err != nil {
		return err
	}

	metrics.GetCounter("repository_file_indexed", index).Inc()

	// Add administration data.
	if file.ItRiusoCodiceIPA != "" {

		// Put administrations data in ES.
		_, err = elasticClient.Index().
			Index("administration").
			Type("administration").
			Id(file.ItRiusoCodiceIPA).
			BodyJson(administration{
				Name:      file.ItRiusoCodiceIPALabel,
				CodiceIPA: file.ItRiusoCodiceIPA,
			}).
			Do(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// getOembedInfo retrive the oembed info from a link.
// Reference: https://oembed.com/providers.json
func getOembedInfo(t, link string) string { // nolint: unparam
	html := ""
	// Fail fast on empty links.
	if link == "" {
		return html
	}

	// Load oembed library and providers.js.
	oe := oembed.NewOembed()
	dataFile, err := Asset("data/oembed_providers.json")
	if err != nil {
		log.Errorf("Error retrieving assets in getOembedInfo.")
		return html
	}
	providers := dataFile
	err = oe.ParseProviders(bytes.NewReader(providers))
	if err != nil {
		log.Errorf("Error parsing providers in getOembedInfo.")
		return html
	}

	item := oe.FindItem(link)
	if item != nil {
		// Extract infos.
		info, err := item.FetchOembed(oembed.Options{URL: link})
		if err != nil {
			log.Errorf("Error fetching oembed in getOembedInfo.")
			return html
		}

		if info.Status >= 300 {
			log.Errorf("Error retrieving info in getOembedInfo.")
			return html
		}

		log.Debugf("Successfully extracted oembed data.")
		html = info.HTML
		return html
	}

	return html
}

// concatenateLink returns the host path joined with the file name.
func concatenateLink(host, file string) string {
	u, err := url.Parse(host)
	if err != nil {
		return ""
	}

	u.Path = path.Join(u.Path, file)

	return u.String()
}

// concatenateLinks returns a list of host paths joined with the file name.
func concatenateLinks(host string, files []string) []string {
	var concatenateFiles []string

	for _, file := range files {
		concatenateFiles = append(concatenateFiles, concatenateLink(host, file))
	}

	return concatenateFiles
}
