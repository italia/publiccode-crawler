package jekyll

import (
	"context"
	"encoding/json"
	"os"

	"github.com/ghodss/yaml"
	"github.com/italia/developers-italia-backend/elastic"
	es "github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// shortSoftware is the subset of a software document that we want to output
type shortSoftware struct {
	ID         string `json:"id"`
	CrawlTime  string `json:"crawltime"`
	PublicCode struct {
		Name string `json:"name"`
		Logo string `json:"logo"`
		URL  string `json:"url"`
		It   struct {
			Riuso struct {
				CodiceIPA string `json:"codiceIPA"`
			} `json:"riuso"`
		} `json:"it"`
	} `json:"publiccode"`
}

// FirstSoftwareRiuso generates a YAML file with simplified info about software, ordered by releaseDate.
func FirstSoftwareRiuso(filename string, results int, elasticClient *es.Client) error {
	query := elastic.NewBoolQuery("software")
	query = query.Must(es.NewExistsQuery("publiccode.it.riuso.codiceIPA"))

	return exportSoftwareList(query, filename, results, elasticClient)
}

// FirstSoftwareOpenSource generates a YAML file with simplified info about software, ordered by releaseDate.
func FirstSoftwareOpenSource(filename string, results int, elasticClient *es.Client) error {
	query := elastic.NewBoolQuery("software")
	query = query.MustNot(es.NewExistsQuery("publiccode.it.riuso.codiceIPA"))

	return exportSoftwareList(query, filename, results, elasticClient)
}

// exportSoftwareList generates a yml file with simplified info about software, ordered by releaseDate.
func exportSoftwareList(query *es.BoolQuery, filename string, results int, elasticClient *es.Client) error {
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
		if err := json.Unmarshal(hit.Source, &sw); err != nil {
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
