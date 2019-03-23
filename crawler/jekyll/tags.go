package jekyll

import (
	"context"
	"encoding/json"
	"os"

	"github.com/ghodss/yaml"
	"github.com/icza/dyno"
	"github.com/italia/developers-italia-backend/crawler/elastic"
	es "github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
)

// CategoriesYML generates a YAML file containing all the categories in ES.
func CategoriesYML(categoriesDestFile string, elasticClient *es.Client) error {
	log.Infof("Generating %s", categoriesDestFile)

	// Extract all the softwares.
	query := elastic.NewBoolQuery("software")
	searchResult, err := elasticClient.Search().
		Index(viper.GetString("ELASTIC_PUBLICCODE_INDEX")). // search in index "publiccode"
		Query(query).                                       // specify the query
		From(0).Size(10000).                                // get first 10k elements. The limit can be changed in ES.
		Do(context.Background())                            // execute
	if err != nil {
		log.Error(err)
	}

	// Result tag list.
	var categories []string

	for _, hit := range searchResult.Hits.Hits {
		var v interface{}
		if err := json.Unmarshal(*hit.Source, &v); err != nil {
			log.Error(err)
		}

		// TODO: we should just ask Elasticsearch for the unique values
		// instead of computing them ourselves.

		// Range over categories.
		if swTags, err := dyno.GetSlice(v, "publiccode", "categories"); err == nil {
			for _, tag := range swTags {
				categories = append(categories, tag.(string))
			}
		}
	}

	categories = funk.Uniq(categories).([]string)

	// Debug note if file will be empty.
	if len(categories) == 0 {
		log.Warn("No categories found")
	}

	return writeYAMLList(&categories, categoriesDestFile)
}

// ScopesYML exports a YAML file containing the list of the distinct scopes mentioned in the catalog.
func ScopesYML(destFile string, elasticClient *es.Client) error {
	log.Infof("Generating %s", destFile)

	// Extract all the softwares.
	query := elastic.NewBoolQuery("software")
	agg := es.NewTermsAggregation().Field("publiccode.intendedAudience.scope").Size(10000).OrderByTermAsc()
	searchResult, err := elasticClient.Search().
		Index(viper.GetString("ELASTIC_PUBLICCODE_INDEX")). // search in index "publiccode"
		Query(query).                                       // specify the query
		Aggregation("publiccode.intendedAudience.scope", agg).
		From(0).Size(10000).     // get first 10k elements. The limit can be changed in ES.
		Do(context.Background()) // execute
	if err != nil {
		log.Error(err)
	}

	aggRes, ok := searchResult.Aggregations.Terms("publiccode.intendedAudience.scope")
	if !ok {
		log.Error("did not find publiccode.intendedAudience.scope in Elasticsearch response")
	}

	var scopes []string
	for _, bucket := range aggRes.Buckets {
		scopes = append(scopes, bucket.Key.(string))
	}

	return writeYAMLList(&scopes, destFile)
}

func writeYAMLList(list *[]string, destFile string) error {
	// Create file if not exists.
	if _, err := os.Stat(destFile); os.IsExist(err) {
		err := os.Remove(destFile)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(destFile)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	// Open file.
	f, err := os.OpenFile(destFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close() // nolint: errcheck

	// Marshal yml.
	d, err := yaml.Marshal(list)
	if err != nil {
		return err
	}
	// Append data to file.
	if _, err = f.WriteString(string(d)); err != nil {
		return err
	}

	return err
}
