package jekyll

import (
	"github.com/spf13/viper"
	"context"
	"os"
	"encoding/json"
	"github.com/icza/dyno"

	"github.com/italia/developers-italia-backend/crawler/elastic"
	"github.com/italia/developers-italia-backend/crawler/ipa"
	log "github.com/sirupsen/logrus"
	"github.com/ghodss/yaml"
	es "github.com/olivere/elastic"
)

// AmministrazioniYML generate a yml file with all the amministrazioni in es.
func AmministrazioniYML(filename string, elasticClient *es.Client) error {
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

	query := elastic.NewBoolQuery("software")
	query = query.Must(es.NewExistsQuery("publiccode.it.riuso.codiceIPA"))

	searchResult, err := elasticClient.Search().
		Index(viper.GetString("ELASTIC_PUBLICCODE_INDEX")).     // search in index "publiccode"
		Query(query).            // specify the query
		Pretty(true).            // pretty print request and response JSON
		From(0).Size(10000).     // get first 10k elements. It can be changed.
		Do(context.Background()) // execute
	if err != nil {
		log.Error(err)
	}

	// Administrations data.
	type administrationType struct{
		CodiceIPA string `json:"ipa"`
		Name      string `json:"name"`
	}
	var administrations []administrationType

	seen := make(map[string]struct{})
	for _, hit := range searchResult.Hits.Hits {
		var v interface{}
		if err := json.Unmarshal(*hit.Source, &v); err != nil {
			log.Error(err)
		}

		// TODO: we should just ask Elasticsearch for the unique values
		// instead of computing them ourselves.

		codiceIPA, _ := dyno.GetString(v, "publiccode", "it", "riuso", "codiceIPA")
		if _, ok := seen[codiceIPA]; !ok {
			seen[codiceIPA] = struct{}{}
			administrations = append(administrations, administrationType{
				codiceIPA,
				ipa.GetAdministrationName(codiceIPA),
			})
		}
	}

	// Debug note if file will be empty.
	if len(administrations) == 0 {
		log.Warnf("%s is empty.", filename)
	}

	// Marshal yml.
	d, err := yaml.Marshal(&administrations)
	if err != nil {
		return err
	}

	// Append data to file.
	if _, err = f.WriteString(string(d)); err != nil {
		return err
	}

	return err
}

