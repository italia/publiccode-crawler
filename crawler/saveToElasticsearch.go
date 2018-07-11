package crawler

import (
	"context"
	"net/url"

	"github.com/italia/developers-italia-backend/ipa"
	"github.com/italia/developers-italia-backend/metrics"
	pcode "github.com/italia/developers-italia-backend/publiccode.yml-parser-go"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

// SaveToES save the chosen data []byte in elasticsearch
func SaveToES(fileRawURL string, domain Domain, name string, activityIndex float64, vitality []int, data []byte, index string, elasticClient *elastic.Client) error {
	// Starting with elastic.v5, you must pass a context to execute each service.
	ctx := context.Background()

	// Generate publiccode data using the parser.
	pc := pcode.PublicCode{}
	err := pcode.Parse(data, &pc)
	if err != nil {
		log.Errorf("Error parsing publiccode.yml for %s: %v", name, err)
	}

	// Add a document to the index.
	file := PublicCodeES{
		FileRawURL:            fileRawURL,
		ItRiusoCodiceIPALabel: ipa.GetAdministrationName(pc.It.Riuso.CodiceIPA),

		Name:             pc.Name,
		ApplicationSuite: pc.ApplicationSuite,
		URL:              pc.URL.String(),
		LandingURL:       pc.LandingURL.String(),

		IsBasedOn:       pc.IsBasedOn,
		SoftwareVersion: pc.SoftwareVersion,
		ReleaseDate:     pc.ReleaseDate.Format("2006-01-02"),
		Logo:            pc.Logo,
		MonochromeLogo:  pc.MonochromeLogo,
		InputTypes:      pc.InputTypes,
		OutputTypes:     pc.OutputTypes,

		Platforms: pc.Platforms,

		Tags: pc.Tags,

		FreeTags: pc.FreeTags,

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
		LegalAuthorsFile:        pc.Legal.AuthorsFile,

		MaintenanceType:        pc.Maintenance.Type,
		MaintenanceContractors: []Contractor{},
		MaintenanceContacts:    []Contact{},

		LocalisationLocalisationReady:  pc.Localisation.LocalisationReady,
		LocalisationAvailableLanguages: pc.Localisation.AvailableLanguages,

		DependenciesOpen:        []Dependency{},
		DependenciesProprietary: []Dependency{},
		DependenciesHardware:    []Dependency{},

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
			FeatureList:      pc.Description[lang].FeatureList,
			Screenshots: func(screenshots []string) []string {
				var s []string
				s = append(s, screenshots...)
				return s
			}(pc.Description[lang].Screenshots),
			Videos: func(videos []*url.URL) []string {
				var v []string
				for _, video := range videos {
					v = append(v, video.String())
				}
				return v
			}(pc.Description[lang].Videos),
			Awards: pc.Description[lang].Awards,
		}

	}
	for _, dependency := range pc.Dependencies.Open {
		file.DependenciesOpen = append(file.DependenciesOpen, Dependency{
			Name:       dependency.Name,
			VersionMin: dependency.VersionMin,
			VersionMax: dependency.VersionMax,
			Optional:   dependency.Optional,
			Version:    dependency.Version,
		})
	}
	for _, dependency := range pc.Dependencies.Proprietary {
		file.DependenciesProprietary = append(file.DependenciesProprietary, Dependency{
			Name:       dependency.Name,
			VersionMin: dependency.VersionMin,
			VersionMax: dependency.VersionMax,
			Optional:   dependency.Optional,
			Version:    dependency.Version,
		})
	}
	for _, dependency := range pc.Dependencies.Hardware {
		file.DependenciesHardware = append(file.DependenciesHardware, Dependency{
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
		Id(domain.Host + "/" + name + "_" + index).
		BodyJson(file).
		Do(ctx)
	if err != nil {
		return err
	}

	metrics.GetCounter("repository_file_indexed", index).Inc()

	return nil
}
