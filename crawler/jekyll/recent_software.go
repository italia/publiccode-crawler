package jekyll

import (
	"context"
	"encoding/json"
	"os"

	"github.com/ghodss/yaml"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// shortSoftware is the subset of a software document that we want to output
type shortSoftware struct {
	ID         string `json:"id"`
	CrawlTime  string `json:"crawltime"`
	PublicCode struct {
		Name      string `json:"name"`
		Logo      string `json:"logo"`
		URL       string `json:"url"`
		CodiceIPA string `json:"it-riuso-codice-ipa"`
	} `json:"publiccode"`
}

// FirstSoftwareRiuso generates a YAML file with simplified info about software, ordered by releaseDate.
func FirstSoftwareRiuso(filename string, results int, unsupportedCountries []string, elasticClient *elastic.Client) error {
	query := elastic.NewBoolQuery()
	query = query.Must(elastic.NewExistsQuery("publiccode.it-riuso-codice-ipa"))

	return exportSoftwareList(query, filename, results, unsupportedCountries, elasticClient)
}

// FirstSoftwareOpenSource generates a YAML file with simplified info about software, ordered by releaseDate.
func FirstSoftwareOpenSource(filename string, results int, unsupportedCountries []string, elasticClient *elastic.Client) error {
	query := elastic.NewBoolQuery()
	query = query.MustNot(elastic.NewExistsQuery("publiccode.it-riuso-codice-ipa"))

	return exportSoftwareList(query, filename, results, unsupportedCountries, elasticClient)
}

// exportSoftwareList generates a yml file with simplified info about software, ordered by releaseDate.
func exportSoftwareList(query *elastic.BoolQuery, filename string, results int, unsupportedCountries []string, elasticClient *elastic.Client) error {
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

	// UnsupportedCountries.
	uc := make([]interface{}, len(unsupportedCountries))
	for i, v := range unsupportedCountries {
		uc[i] = v
	}
	query = query.MustNot(elastic.NewTermsQuery("publiccode.intended-audience-unsupported-countries", uc...))
	query = query.Filter(elastic.NewTypeQuery("software"))

	// Extract all the documents.
	searchResult, err := elasticClient.Search().
		Index(viper.GetString("ELASTIC_PUBLICCODE_INDEX")). // search in index "publiccode"
		Query(query).                                       // specify the query
		Sort("publiccode.releaseDate", false).              // sort by releaseDate, from newest to oldest.
		Pretty(true).                                       // pretty print request and response JSON
		From(0).Size(results).                              // get first 10k elements. It can be changed.
		Do(context.Background())                            // execute
	if err != nil {
		log.Error(err)
	}

	var items []shortSoftware
	for _, hit := range searchResult.Hits.Hits {
		var sw shortSoftware
		if err := json.Unmarshal(*hit.Source, &sw); err != nil {
			log.Error(err)
		}
		items = append(items, sw)
	}

	// Debug note if file will be empty.
	if len(items) == 0 {
		log.Warnf("%s is empty.", filename)
	}

	// Marshal yml.
	d, err := yaml.Marshal(&items)
	if err != nil {
		return err
	}
	// Append data to file.
	if _, err = f.WriteString(string(d)); err != nil {
		return err
	}

	return err
}
