package jekyll

import (
	"context"
	"net/url"
	"os"
	"path"
	"reflect"
	"strings"

	yaml "github.com/ghodss/yaml"
	"github.com/italia/developers-italia-backend/crawler"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// AllSoftwareYML generate the softwares.yml file
func AllSoftwareYML(filename string, numberOfSimilarSoftware int, elasticClient *elastic.Client) error {
	log.Debug("Generating softwares.yml")
	// Create file if not exists.
	if _, err := os.Stat(filename); os.IsExist(err) {
		err := os.Remove(filename)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	// Open file.
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close() // nolint: errcheck

	// Publiccodes data.
	var softwares []Software

	// Extract all the documents.
	searchResult, err := elasticClient.Search().
		Index("publiccode").               // search in index "publiccode"
		Query(elastic.NewMatchAllQuery()). // specify the query
		Pretty(true).                      // pretty print request and response JSON
		From(0).Size(10000).               // get first 10k elements. The limit can be changed in ES.
		Do(context.Background())           // execute
	if err != nil {
		log.Error(err)
	}

	var pctype crawler.PublicCodeES
	for _, item := range searchResult.Each(reflect.TypeOf(pctype)) {
		i := item.(crawler.PublicCodeES)

		rawBaseDir := strings.TrimRight(i.FileRawURL, viper.GetString("CRAWLED_FILENAME"))

		softwareExtracted := Software{
			Name:              i.Name,
			ApplicationSuite:  i.ApplicationSuite,
			URL:               i.URL,
			LandingURL:        i.LandingURL,
			IsBasedOn:         i.IsBasedOn,
			SoftwareVersion:   i.SoftwareVersion,
			ReleaseDate:       i.ReleaseDate,
			Logo:              concatenateLink(rawBaseDir, i.Logo),
			MonochromeLogo:    concatenateLink(rawBaseDir, i.MonochromeLogo),
			Platforms:         i.Platforms,
			Tags:              i.Tags,
			FreeTags:          populateFreeTags(i.FreeTags),
			PopularTags:       i.Tags, // todo
			ShareTags:         i.Tags, // todo
			UsedBy:            i.UsedBy,
			Roadmap:           i.Roadmap,
			DevelopmentStatus: i.DevelopmentStatus,
			VitalityScore:     i.VitalityScore,
			VitalityDataChart: i.VitalityDataChart,
			SoftwareType: SoftwareTypeData{
				Type: i.SoftwareType,
			},
			IntendedAudience: IntendedAudienceData{
				OnlyFor:              i.IntendedAudienceOnlyFor,
				Countries:            i.IntendedAudienceCountries,
				UnsupportedCountries: i.IntendedAudienceUnsupportedCountries,
			},
			Description:    i.Description,
			OldVariant:     []OldVariantData{},
			OldFeatureList: OldFeatureListData{},
			TagsRelate:     i.Tags,
			Legal: LegalData{
				License:            i.LegalLicense,
				MainCopyrightOwner: i.LegalMainCopyrightOwner,
				RepoOwner:          i.LegalRepoOwner,
				AuthorsFile:        concatenateLink(rawBaseDir, i.LegalAuthorsFile),
			},
			Localisation: LocalisationData{
				LocalisationReady:  i.LocalisationLocalisationReady,
				AvailableLanguages: i.LocalisationAvailableLanguages,
			},
			Dependencies: DependenciesData{
				Open:        i.DependenciesOpen,
				Proprietary: i.DependenciesProprietary,
				Hardware:    i.DependenciesProprietary,
			},
			It: ExtensionIT{
				Accessibile:    i.ItConformeAccessibile,
				Interoperabile: i.ItConformeInteroperabile,
				Riuso: ItRiusoData{
					CodiceIPA: i.ItRiusoCodiceIPA,
				},
				Spid:   i.ItSpid,
				Pagopa: i.ItPagopa,
				Cie:    i.ItCie,
				Anpr:   i.ItAnpr,
				DesignKit: DesignKitData{
					Seo: i.ItDesignKitSeo,
					UI:  i.ItDesignKitUI,
					Web: i.ItDesignKitWeb,
				},
			},
		}

		// Populate complex data.
		softwareExtracted.Maintenance.Contractors = i.MaintenanceContractors
		softwareExtracted.Maintenance.Contacts = i.MaintenanceContacts
		softwareExtracted.Maintenance.Type = i.MaintenanceType

		for lang := range softwareExtracted.Description {
			for n := range softwareExtracted.Description[lang].Screenshots {
				softwareExtracted.Description[lang].Screenshots[n] = concatenateLink(rawBaseDir, softwareExtracted.Description[lang].Screenshots[n])
			}
		}

		// Search similar softwares for this software and add them to olSoftwares.
		similarSoftware := findSimilarSoftwares(i.Tags, numberOfSimilarSoftware, elasticClient)
		for _, v := range similarSoftware {
			// Remove the extracted software.
			if v.URL != softwareExtracted.URL {
				similarBaseDir := strings.TrimRight(v.FileRawURL, viper.GetString("CRAWLED_FILENAME"))
				related := RelatedSoftware{
					Name:  v.Name,
					Image: concatenateLink(similarBaseDir, v.Logo),
				}
				if d, ok := v.Description["eng"]; ok {
					related.Eng.LocalisedName = d.LocalisedName
					related.Eng.URL = v.URL
				}
				if d, ok := v.Description["ita"]; ok {
					related.Ita.LocalisedName = d.LocalisedName
					related.Ita.URL = v.URL
				}

				softwareExtracted.RelatedSoftwares = append(softwareExtracted.RelatedSoftwares, related)
			}
		}

		// Search softwares basedOn this one.
		isBasedOnSoftware := findIsBasedOnSoftwares(i, elasticClient)
		softwareExtracted.OldVariant = append(softwareExtracted.OldVariant, findOldVariant(isBasedOnSoftware, softwareExtracted)...)

		// Diff features.
		diffFeatures := findDiffFeatures(softwareExtracted)
		softwareExtracted.OldFeatureList.Eng = diffFeatures.Eng
		softwareExtracted.OldFeatureList.Ita = diffFeatures.Ita

		// Append.
		softwares = append(softwares, softwareExtracted)
	}

	// Marshal yml.
	d, err := yaml.Marshal(&softwares)
	if err != nil {
		log.Error(err)
	}
	//Append data to file.
	if _, err = f.WriteString(string(d)); err != nil {
		return err
	}

	return err
}

func findOldVariant(isBasedOnSoftware []crawler.PublicCodeES, softwareExtracted Software) []OldVariantData {
	for _, v := range isBasedOnSoftware {
		// Remove the extracted software.
		if v.URL != softwareExtracted.URL {
			basedOn := OldVariantData{
				Name:          v.Name,
				VitalityScore: v.VitalityScore,
			}
			basedOn.Legal.RepoOwner = v.LegalRepoOwner

			if d, ok := v.Description["eng"]; ok {
				basedOn.Eng.Features = d.FeatureList
				basedOn.Eng.URL = v.URL
			}
			if d, ok := v.Description["ita"]; ok {
				basedOn.Ita.Features = d.FeatureList
				basedOn.Ita.URL = v.URL
			}

			softwareExtracted.OldVariant = append(softwareExtracted.OldVariant, basedOn)
		}
	}
	return softwareExtracted.OldVariant
}

func findDiffFeatures(softwareExtracted Software) OldFeatureListData {
	var diffFeatures OldFeatureListData

	for _, variant := range softwareExtracted.OldVariant {
		// Diff for eng.
		for _, oldFeature := range variant.Eng.Features {
			if !contains(softwareExtracted.Description["eng"].FeatureList, oldFeature) {
				diffFeatures.Eng = append(diffFeatures.Eng, oldFeature)
			}
		}
		//Diff for ita.
		for _, oldFeature := range variant.Ita.Features {
			if !contains(softwareExtracted.Description["ita"].FeatureList, oldFeature) {
				diffFeatures.Ita = append(diffFeatures.Ita, oldFeature)
			}
		}
	}

	return diffFeatures
}

func findSimilarSoftwares(tags []string, numberOfSimilarSoftware int, elasticClient *elastic.Client) []crawler.PublicCodeES {
	var pcs []crawler.PublicCodeES

	// Generate query.
	query := elastic.NewBoolQuery()
	for _, tag := range tags {
		query = query.Should(elastic.NewTermQuery("tags", tag))
	}

	searchResult, err := elasticClient.Search().
		Index("publiccode").                   // search in index "publiccode"
		Query(query).                          // specify the query
		From(0).Size(numberOfSimilarSoftware). // take documents from 0-numberOfSimilarSoftware
		Pretty(true).                          // pretty print request and response JSON
		Do(context.Background())               // execute
	if err != nil {
		log.Error(err)
	}
	var pctype crawler.PublicCodeES
	for _, item := range searchResult.Each(reflect.TypeOf(pctype)) {
		i := item.(crawler.PublicCodeES)
		pcs = append(pcs, i)
	}

	return pcs

}

func findIsBasedOnSoftwares(document crawler.PublicCodeES, elasticClient *elastic.Client) []crawler.PublicCodeES {
	var pcs []crawler.PublicCodeES

	// Extract all the documents. It should filter only the ones with isBaseOn=url
	searchResult, err := elasticClient.Search().
		Index("publiccode").               // search in index "publiccode"
		Query(elastic.NewMatchAllQuery()). // specify the query
		Pretty(true).                      // pretty print request and response JSON
		From(0).Size(10000).               // get first 10k elements. The limit can be changed in ES.
		Do(context.Background())           // execute
	if err != nil {
		log.Error(err)
	}

	var pctype crawler.PublicCodeES
	// Range over isBasedOn
	for _, item := range searchResult.Each(reflect.TypeOf(pctype)) {
		i := item.(crawler.PublicCodeES)
		// If isBasedOn contains url, append to returned software.
		for _, name := range i.IsBasedOn {
			if name == document.URL {
				pcs = append(pcs, i)
			}
		}
		// And viceversa.
		if contains(document.IsBasedOn, i.URL) {
			pcs = append(pcs, i)
		}
	}

	return pcs
}

// contains returns true if the slice of strings contains the searched string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func concatenateLink(host, file string) string {
	u, err := url.Parse(host)
	if err != nil {
		return ""
	}

	u.Path = path.Join(u.Path, file)

	return u.String()
}

func populateFreeTags(freeTags map[string][]string) FreeTagsData {
	var tags FreeTagsData
	for lang, v := range freeTags {
		if lang == "ita" {
			tags.Ita = append(tags.Ita, v...)
		}
		if lang == "eng" {
			tags.Eng = append(tags.Eng, v...)
		}
	}
	return tags
}
