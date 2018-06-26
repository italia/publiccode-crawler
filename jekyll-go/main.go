package main

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v1"
)

type Administration struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	CodiceIPA string `json:"ipa"`
}

func main() {
	file := "generated/amministrazioni.yml"
	err := AmministrazioniYML(file)
	if err != nil {
		log.Error(err)
	}
}

func AmministrazioniYML(filename string) error {
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
	file.Close()
	// Open file.
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	// Administrations data.
	var administrations []Administration

	// Elastic connection.
	elasticClient, err := crawler.ElasticClientFactory(
		"http://localhost:9200",
		"",
		"")
	if err != nil {
		fmt.Println("error connecting es")
		log.Error(err)
	}

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
