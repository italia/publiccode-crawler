package jekyll

import (
	"github.com/spf13/viper"
	"context"
	"os"
	"reflect"
	"sort"

	"github.com/italia/developers-italia-backend/crawler/crawler"
	yaml "github.com/ghodss/yaml"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

// AllSoftwareYML generate the softwares.yml file
func AllSoftwareYML(filename string, numberOfSimilarSoftware, numberOfPopularTags int, unsupportedCountries []string, elasticClient *elastic.Client) error {
	log.Infof("Generating %s", filename)
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

	// UnsupportedCountries
	uc := make([]interface{}, len(unsupportedCountries))
	for i, v := range unsupportedCountries {
		uc[i] = v
	}
	// Extract all the softwares.
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTypeQuery("software"))
	query = query.MustNot(elastic.NewTermsQuery("intended-audience-unsupported-countries", uc...))

	searchResult, err := elasticClient.Search().
		Index(viper.GetString("ELASTIC_ALIAS")).     // search in index "publiccode"
		Query(query).            // specify the query
		Pretty(true).            // pretty print request and response JSON
		From(0).Size(10000).     // get first 10k elements. The limit can be changed in ES.
		Do(context.Background()) // execute
	if err != nil {
		log.Error(err)
	}

	var pctype crawler.PublicCodeES
	for _, item := range searchResult.Each(reflect.TypeOf(pctype)) {
		i := item.(crawler.PublicCodeES)

		softwareExtracted := Software{
			AmministrazioneLabel: i.ItRiusoCodiceIPALabel,
			ID:                   i.ID,
			CrawlTime:            i.CrawlTime,
			Name:                 i.Name,
			ApplicationSuite:     i.ApplicationSuite,
			URL:                  i.URL,
			LandingURL:           i.LandingURL,
			IsBasedOn:            i.IsBasedOn,
			SoftwareVersion:      i.SoftwareVersion,
			ReleaseDate:          i.ReleaseDate,
			Logo:                 i.Logo,
			MonochromeLogo:       i.MonochromeLogo,
			Platforms:            i.Platforms,
			Tags:                 i.Tags,
			FreeTags:             populateFreeTags(i.Description),
			PopularTags:          populatePopularTags(i.Tags, numberOfPopularTags, elasticClient), // PopularTags are the first n tags that are more popular.
			ShareTags:            i.Tags,                                                          // ShareTags are tags.
			UsedBy:               i.UsedBy,
			Roadmap:              i.Roadmap,
			DevelopmentStatus:    i.DevelopmentStatus,
			VitalityScore:        i.VitalityScore,
			VitalityDataChart:    i.VitalityDataChart,
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
			OldFeatureList: map[string][]string{},
			TagsRelate:     i.Tags,
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
				Open:        i.DependsOnOpen,
				Proprietary: i.DependsOnHardware,
				Hardware:    i.DependsOnProprietary,
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
		similarSoftware := findSimilarSoftwares(i.Tags, numberOfSimilarSoftware, unsupportedCountries, elasticClient)
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

		// Search softwares basedOn this one.
		isBasedOnSoftware := findIsBasedOnSoftwares(i, unsupportedCountries, elasticClient)
		softwareExtracted.OldVariant = findOldVariants(isBasedOnSoftware, softwareExtracted)

		// Diff features.
		softwareExtracted.OldFeatureList = findDiffFeatures(softwareExtracted)

		// Append only supported countries.
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

// findOldVariants returns a list of variants for the given software.
func findOldVariants(isBasedOnSoftware []crawler.PublicCodeES, softwareExtracted Software) []OldVariantData {
	var oldVariantData []OldVariantData

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

			oldVariantData = append(oldVariantData, basedOn)

		}
	}
	return oldVariantData
}

// findDiffFeatures returns a list of features from a software.
func findDiffFeatures(softwareExtracted Software) map[string][]string {
	diffFeatures := map[string][]string{}
	for _, variant := range softwareExtracted.OldVariant {
		// Diff for eng.
		for _, oldFeature := range variant.Eng.Features {
			if !contains(softwareExtracted.Description["eng"].FeatureList, oldFeature) {
				diffFeatures["eng"] = append(diffFeatures["eng"], oldFeature)
			}
		}
		//Diff for ita.
		for _, oldFeature := range variant.Ita.Features {
			if !contains(softwareExtracted.Description["ita"].FeatureList, oldFeature) {
				diffFeatures["ita"] = append(diffFeatures["ita"], oldFeature)
			}
		}
	}

	return diffFeatures
}

// findSimilarSoftwares returns a list of similar software based on tags.
func findSimilarSoftwares(tags []string, numberOfSimilarSoftware int, unsupportedCountries []string, elasticClient *elastic.Client) []crawler.PublicCodeES {
	var pcs []crawler.PublicCodeES

	// UnsupportedCountries.
	uc := make([]interface{}, len(unsupportedCountries))
	for i, v := range unsupportedCountries {
		uc[i] = v
	}

	// Generate query.
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTypeQuery("software"))
	for _, tag := range tags {
		query = query.Should(elastic.NewTermQuery("tags", tag))
	}
	query = query.MustNot(elastic.NewTermsQuery("intended-audience-unsupported-countries", uc...))

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

func findIsBasedOnSoftwares(document crawler.PublicCodeES, unsupportedCountries []string, elasticClient *elastic.Client) []crawler.PublicCodeES {
	var pcs []crawler.PublicCodeES

	// UnsupportedCountries.
	uc := make([]interface{}, len(unsupportedCountries))
	for i, v := range unsupportedCountries {
		uc[i] = v
	}

	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTypeQuery("software"))
	query = query.MustNot(elastic.NewTermsQuery("intended-audience-unsupported-countries", uc...))
	// Extract all the documents. It should filter only the ones with isBaseOn=url.
	searchResult, err := elasticClient.Search().
		Index("publiccode").     // search in index "publiccode"
		Query(query).            // specify the query
		Pretty(true).            // pretty print request and response JSON
		From(0).Size(10000).     // get first 10k elements. The limit can be changed in ES.
		Do(context.Background()) // execute
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

func populateFreeTags(description map[string]crawler.Desc) map[string][]string {
	freeTags := make(map[string][]string)
	for lang, desc := range description {
		freeTags[lang] = desc.FreeTags
	}
	return freeTags
}

func populatePopularTags(tags []string, number int, elasticClient *elastic.Client) []string {
	if len(tags) < number {
		return tags
	}

	var popularTags []string

	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTypeQuery("software"))
	query = query.MustNot(elastic.NewTermsQuery("intended-audience-unsupported-countries", "it", "us", "de"))
	// Extract all the documents. It should filter only the ones with isBaseOn=url.
	searchResult, err := elasticClient.Search().
		Index("publiccode").     // search in index "publiccode"
		Query(query).            // specify the query
		Pretty(true).            // pretty print request and response JSON
		From(0).Size(10000).     // get first 10k elements. The limit can be changed in ES.
		Do(context.Background()) // execute
	if err != nil {
		log.Error(err)
	}

	results := map[string]int{}

	// Range over the publiccodes in ES.
	var pctype crawler.PublicCodeES
	for _, item := range searchResult.Each(reflect.TypeOf(pctype)) {
		i := item.(crawler.PublicCodeES)
		for _, v := range i.Tags {
			results[v]++
		}
	}

	// Order the map into a slice.
	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range results {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	// Populate the popularTags slice with most popular tags.
	for n, kv := range ss {
		if n < number {
			break
		}
		popularTags = append(popularTags, kv.Key)
	}

	return popularTags
}
