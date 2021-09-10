package jekyll

import (
	"context"
	"os"

	"github.com/ghodss/yaml"
	"github.com/italia/developers-italia-backend/crawler/elastic"
	es "github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// CategoriesYML generates a YAML file containing all the categories in ES.
func CategoriesYML(destFile string, elasticClient *es.Client) error {
	return exportDistinctValuesToYAML("publiccode.categories", destFile, elasticClient)
}

// ScopesYML exports a YAML file containing the list of the distinct scopes mentioned in the catalog.
func ScopesYML(destFile string, elasticClient *es.Client) error {
	return exportDistinctValuesToYAML("publiccode.intendedAudience.scope", destFile, elasticClient)
}

func exportDistinctValuesToYAML(key, destFile string, elasticClient *es.Client) error {
	log.Infof("Generating %s", destFile)

	// Extract all the softwares.
	query := elastic.NewBoolQuery("software")
	agg := es.NewTermsAggregation().Field(key).Size(10000).OrderByTermAsc()
	searchResult, err := elasticClient.Search().
		Index(viper.GetString("ELASTIC_PUBLICCODE_INDEX")). // search in index "publiccode"
		Query(query).                                       // specify the query
		Aggregation(key, agg).
		From(0).Size(10000).     // get first 10k elements. The limit can be changed in ES.
		Do(context.Background()) // execute
	if err != nil {
		log.Error(err)
	}

	aggRes, ok := searchResult.Aggregations.Terms(key)
	if !ok {
		log.Errorf("did not find %v in Elasticsearch response", key)
	}

	var values []string
	for _, bucket := range aggRes.Buckets {
		values = append(values, bucket.Key.(string))
	}

	return writeYAMLList(&values, destFile)
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
