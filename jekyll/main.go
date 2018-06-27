package jekyll

import (
	"context"
	"os"
	"reflect"

	"github.com/davecgh/go-spew/spew"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v1"
)

// Administration is a simple description of an Administration.
type Administration struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	CodiceIPA string `json:"ipa"`
}

// AmministrazioniYML generate a yml file with all the amministrazioni in es.
func AmministrazioniYML(filename string, elasticClient *elastic.Client) error {
	// Uncommnet when the publiccodes is ready to be written on file.
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

	// Administrations data.
	var administrations []Administration

	// Extract all the documents.
	searchResult, err := elasticClient.Search().
		Index("publiccode").               // search in index "publiccode"
		Query(elastic.NewMatchAllQuery()). // specify the query
		Pretty(true).                      // pretty print request and response JSON
		From(0).Size(10000).               // get first 10k elements. It can be changed.
		Do(context.Background())           // execute
	if err != nil {
		log.Error(err)
	}

	// Foreach search result check if codiceIPA is not empty.
	var pctype PublicCode
	for _, item := range searchResult.Each(reflect.TypeOf(pctype)) {
		i := item.(PublicCode)
		// Debug.
		spew.Dump(i)

		if i.ItRiusoCodiceIPA != "" {
			administrations = append(administrations, Administration{
				Name:      i.Name,
				URL:       i.URL,
				CodiceIPA: i.ItRiusoCodiceIPA,
			})

		}
	}

	// Marshal yml.
	d, err := yaml.Marshal(&administrations)
	if err != nil {
		log.Error(err)
	}

	//Append data to file.
	if _, err = f.WriteString(string(d)); err != nil {
		return err
	}

	return err
}
