package jekyll

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

func AllSoftwareYML(filename string, numberOfSimilarSoftware int, elasticClient *elastic.Client) error {
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
	var publiccodes []PublicCode

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

	var pctype PublicCode
	for _, item := range searchResult.Each(reflect.TypeOf(pctype)) {
		i := item.(PublicCode)
		publiccodes = append(publiccodes, i)

		// Search similar softwares for this software and add them to olSoftwares.
		similarSoftware := findSimilarSoftwares(i.Tags, numberOfSimilarSoftware, elasticClient)
		// TODO: add similarSoftwaresNames to similar software list.
		for _, _ = range similarSoftware {
			// TODO: append to relatedSoftwares for this item
			// https://github.com/italia/developers.italia.it/blob/new-version-master/_data/softwares.yml#L200
		}

		// Search softwares basedOn this one.
		isBasedOnSoftware := findisBasedOnSoftwares(i.URL, elasticClient)

		// otherVariantsFeaturesList will be populated with every feature not present into this software.
		var otherVariantsFeaturesList map[string][]string
		var oldVariant []OldVariant

		for _, pc := range isBasedOnSoftware {
			variant := OldVariant{
				Name: i.Name,
				URL:  i.URL,
			}
			// for every language in Description
			for language, desc := range pc.Description {
				// Prepare the list of otherVariantsFeaturesList.
				for _, feature := range desc.FeatureList {
					if !contains(i.Description[language].FeatureList, feature) {
						otherVariantsFeaturesList[language] = append(otherVariantsFeaturesList[language], feature)
					}
				}
			}
			oldVariant = append(oldVariant, variant)

		}

		fmt.Printf("similarSoftware for %s: %+v\n", i.Name, similarSoftware)
		fmt.Printf("isBasedOnSoftware for %s: %+v\n\n\n", i.Name, isBasedOnSoftware)
		fmt.Printf("Complete data: %+v\n", i)
	}

	// Marshal yml.
	d, err := yaml.Marshal(&publiccodes)
	if err != nil {
		log.Error(err)
	}

	//Append data to file.
	if _, err = f.WriteString(string(d)); err != nil {
		return err
	}

	return err
}

func findSimilarSoftwares(tags []string, numberOfSimilarSoftware int, elasticClient *elastic.Client) []PublicCode {
	var pcs []PublicCode

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
	var pctype PublicCode
	for _, item := range searchResult.Each(reflect.TypeOf(pctype)) {
		i := item.(PublicCode)
		pcs = append(pcs, i)
	}

	return pcs

}

func findisBasedOnSoftwares(url string, elasticClient *elastic.Client) []PublicCode {
	var pcs []PublicCode
	query := elastic.NewBoolQuery()
	query = query.Must(elastic.NewTermQuery("is-based-on", url))

	searchResult, err := elasticClient.Search().
		Index("publiccode").     // search in index "publiccode"
		Query(query).            // specify the query
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Error(err)
	}
	var pctype PublicCode
	for _, item := range searchResult.Each(reflect.TypeOf(pctype)) {
		i := item.(PublicCode)
		pcs = append(pcs, i)
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
