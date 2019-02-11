package jekyll

import (
	"context"
	"encoding/json"
	"os"

	"github.com/ghodss/yaml"
	"github.com/icza/dyno"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
)

// TagsYML generate the software-tags.yml that will contain all the tags in ES.
func TagsYML(tagsDestFile string, elasticClient *elastic.Client) error {
	log.Infof("Generating %s", tagsDestFile)

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

	/*
		// UnsupportedCountries.
		uc := make([]interface{}, len(unsupportedCountries))
		for i, v := range unsupportedCountries {
			uc[i] = v
		}
	*/

	// Extract all the softwares.
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTypeQuery("software"))
	//query = query.MustNot(elastic.NewTermsQuery("publiccode.intendedAudience.unsupportedCountries", uc...))

	searchResult, err := elasticClient.Search().
		Index(viper.GetString("ELASTIC_PUBLICCODE_INDEX")). // search in index "publiccode"
		Query(query).                                       // specify the query
		Pretty(true).                                       // pretty print request and response JSON
		From(0).Size(10000).                                // get first 10k elements. The limit can be changed in ES.
		Do(context.Background())                            // execute
	if err != nil {
		log.Error(err)
	}

	// Result tag list.
	var tags []string

	for _, hit := range searchResult.Hits.Hits {
		var v interface{}
		if err := json.Unmarshal(*hit.Source, &v); err != nil {
			log.Error(err)
		}
		
		// TODO: we should just ask Elasticsearch for the unique values
		// instead of computing them ourselves.

		// Range over tags.
		if swTags, err := dyno.GetSlice(v, "publiccode", "tags"); err == nil {
			for _, tag := range swTags {
				tags = append(tags, tag.(string))
			}
		}
	}

	tags = funk.Uniq(tags).([]string)

	// Debug note if file will be empty.
	if len(tags) == 0 {
		log.Warnf("%s is empty.", tagsDestFile)
	}

	// Marshal yml.
	d, err := yaml.Marshal(&tags)
	if err != nil {
		return err
	}
	// Append data to file.
	if _, err = f.WriteString(string(d)); err != nil {
		return err
	}

	return err
}
