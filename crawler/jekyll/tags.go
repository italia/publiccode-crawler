package jekyll

import (
	"github.com/spf13/viper"
	"context"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/italia/developers-italia-backend/crawler/crawler"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// Tag represent a single tag translated in two languages: english (en) and italian (it).
type Tag struct {
	En string `yaml:"en"`
	It string `yaml:"it"`
}

// TagsYML generate the software-tags.yml that will contain all the tags in ES.
func TagsYML(tagsDestFile, tagsSrcFile string, elasticClient *elastic.Client) error {
	log.Infof("Generating %s from %s and ES", tagsDestFile, tagsSrcFile)

	// Create file if not exists.
	if _, err := os.Stat(tagsDestFile); os.IsExist(err) {
		err := os.Remove(tagsDestFile)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(tagsDestFile)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	// Open file.
	f, err := os.OpenFile(tagsDestFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close() // nolint: errcheck

	// Tags data.
	var tags map[string]Tag
	// Read tags data.
	data, err := ioutil.ReadFile(tagsSrcFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &tags)
	if err != nil {
		return err
	}

	// Extract all the softwares.
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTypeQuery("software"))

	searchResult, err := elasticClient.Search().
		Index(viper.GetString("ELASTIC_PUBLICCODE_INDEX")).     // search in index "publiccode"
		Query(query).            // specify the query
		Pretty(true).            // pretty print request and response JSON
		From(0).Size(10000).     // get first 10k elements. The limit can be changed in ES.
		Do(context.Background()) // execute
	if err != nil {
		log.Error(err)
	}

	// Result tag list.
	result := make(map[string]Tag)

	// Foreach search result check if codiceIPA is not empty.
	var pctype crawler.PublicCodeES
	for _, item := range searchResult.Each(reflect.TypeOf(pctype)) {
		i := item.(crawler.PublicCodeES)

		// TODO: add unsupported countries.
		unsupportedCountries := []string{"it"}
		// Append only supported countries.
		unsupported := checkUnsupportedCountries(i.IntendedAudienceUnsupportedCountries, unsupportedCountries)
		if !unsupported {
			// Range over tags.
			for tag, value := range tags {
				if contains(i.Tags, tag) && !containsTags(result, tag) {
					result[tag] = value
				}
			}
		}
	}

	// Debug note if file will be empty.
	if len(result) == 0 {
		log.Warnf("%s is empty.", tagsDestFile)
	}

	// Marshal yml.
	d, err := yaml.Marshal(&result)
	if err != nil {
		return err
	}
	// Append data to file.
	if _, err = f.WriteString(string(d)); err != nil {
		return err
	}

	return err
}

// containsTags returns true if the map key contains the given string.
func containsTags(m map[string]Tag, name string) bool {
	for k := range m {
		if k == name {
			return true
		}
	}
	return false
}

// checkUnsupportedCountries returns true if an unsupported country is in a list of countries.
func checkUnsupportedCountries(listCountries, unsupportedCountries []string) bool {
	for _, unsupportedCountry := range unsupportedCountries {
		if contains(listCountries, unsupportedCountry) {
			return true
		}
	}
	return false
}
