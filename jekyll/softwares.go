package jekyll

import (
	"context"
	"os"
	"reflect"

	yaml "github.com/ghodss/yaml"
	"github.com/italia/developers-italia-backend/crawler"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
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
	file.Close() // nolint: errcheck
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

		softwareExtracted := Software{
			Name:             i.Name,
			ApplicationSuite: i.ApplicationSuite,
			URL:              i.URL,
			LandingURL:       i.LandingURL,
			IsBasedOn:        i.IsBasedOn,
			SoftwareVersion:  i.SoftwareVersion,
			ReleaseDate:      i.ReleaseDate,
			Logo:             i.Logo,
			MonochromeLogo:   i.MonochromeLogo,
			Platforms:        i.Platforms,
			Tags:             i.Tags,
			FreeTags: func(freeTags map[string][]string) FreeTagsData {
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
			}(i.FreeTags),
			PopularTags:       []string{"todo", "popularTagsToCalculate"},
			ShareTags:         []string{"todo", "shareTagsToCalculate"},
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
			OldVariant:     []OldVariantData{},   //todo
			OldFeatureList: OldFeatureListData{}, //todo
			TagsRelate:     []string{"todo", "tagsRelated"},
			Legal: LegalData{
				License:            i.LegalLicense,
				MainCopyrightOwner: i.LegalMainCopyrightOwner,
				RepoOwner:          i.LegalRepoOwner,
				AuthorsFile:        i.LegalAuthorsFile,
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

		// Search similar softwares for this software and add them to olSoftwares.
		similarSoftware := findSimilarSoftwares(i.Tags, numberOfSimilarSoftware, elasticClient)
		for _, v := range similarSoftware {
			// Remove the extracted software.
			if v.URL != softwareExtracted.URL {
				related := RelatedSoftware{
					Name:  v.Name,
					Image: v.Logo,
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
		//
		// 	// Search softwares basedOn this one.
		// 	isBasedOnSoftware := findisBasedOnSoftwares(i.URL, elasticClient)
		//
		// 	// otherVariantsFeaturesList will be populated with every feature not present into this software.
		// 	var otherVariantsFeaturesList map[string][]string
		// 	var oldVariant []OldVariant
		//
		// 	for _, pc := range isBasedOnSoftware {
		// 		variant := OldVariant{
		// 			Name: i.Name,
		// 			URL:  i.URL,
		// 		}
		// 		// for every language in Description
		// 		for language, desc := range pc.Description {
		// 			// Prepare the list of otherVariantsFeaturesList.
		// 			for _, feature := range desc.FeatureList {
		// 				if !contains(i.Description[language].FeatureList, feature) {
		// 					otherVariantsFeaturesList[language] = append(otherVariantsFeaturesList[language], feature)
		// 				}
		// 			}
		// 		}
		// 		oldVariant = append(oldVariant, variant)
		//
		// 	}
		//

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

//
// func findisBasedOnSoftwares(url string, elasticClient *elastic.Client) []Software {
// 	var pcs []PublicCode
// 	query := elastic.NewBoolQuery()
// 	query = query.Must(elastic.NewTermQuery("is-based-on", url))
//
// 	searchResult, err := elasticClient.Search().
// 		Index("publiccode").     // search in index "publiccode"
// 		Query(query).            // specify the query
// 		Pretty(true).            // pretty print request and response JSON
// 		Do(context.Background()) // execute
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	var pctype PublicCode
// 	for _, item := range searchResult.Each(reflect.TypeOf(pctype)) {
// 		i := item.(PublicCode)
// 		pcs = append(pcs, i)
// 	}
//
// 	return pcs
// }
//
// // contains returns true if the slice of strings contains the searched string.
// func contains(slice []string, item string) bool {
// 	for _, s := range slice {
// 		if s == item {
// 			return true
// 		}
// 	}
// 	return false
// }
