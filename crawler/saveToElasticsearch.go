package crawler

import (
	"context"
	"os"

	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

// File is a generic structure for saveToES() function.
// TODO: Will be replaced with a parsed publiccode.PublicCode whith proper mapping.
type File struct {
	Source string `json:"source"`
	Name   string `json:"name"`
	Data   string `json:"data"`
}

// SaveES save the chosen <file_name> in elasticsearch
func SaveToES(domain Domain, name string, data []byte, elasticClient *elastic.Client) {
	const (
		// Elasticsearch mapping for publiccode. Checkout elasticsearch/mappings/software.json
		// TODO: Mapping must reflect the publiccode.PublicCode structure.
		mapping = ""
	)
	index := domain.Index

	// Starting with elastic.v5, you must pass a context to execute each service.
	ctx := context.Background()

	client, err := ElasticClientFactory(
		os.Getenv("ELASTIC_URL"),
		os.Getenv("ELASTIC_USER"),
		os.Getenv("ELASTIC_PWD"))
	if err != nil {
		log.Error(err)
	}
	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		log.Error(err)
	}

	if !exists {
		// Create a new index.
		// TODO: When mapping will be available: client.CreateIndex(index).BodyString(mapping).Do(ctx).
		_, err = client.CreateIndex(index).Do(ctx)
		if err != nil {
			log.Error(err)
		}
	}
	// Add a document to the index.
	file := File{Source: domain.Id, Name: name, Data: string(data)}

	// Put publiccode data in ES.
	put, err := client.Index().
		Index(index).
		Type("doc").
		Id(domain.Id + "/" + name + "_" + domain.Index).
		BodyJson(file).
		Do(ctx)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("Indexed file %s to index %s, type %s", put.Id, put.Index, put.Type)

}
