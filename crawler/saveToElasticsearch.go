package crawler

import (
	"context"

	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"github.com/italia/developers-italia-backend/metrics"
)

// File is a generic structure for saveToES() function.
// TODO: Will be replaced with a parsed publiccode.PublicCode whit proper mapping.
type File struct {
	Source string `json:"source"`
	Name   string `json:"name"`
	Data   string `json:"data"`
}

// SaveES save the chosen <file_name> in elasticsearch
func SaveToES(domain Domain, name string, data []byte, index string, elasticClient *elastic.Client) {
	const (
		// Elasticsearch mapping for publiccode. Checkout elasticsearch/mappings/software.json
		// TODO: Mapping must reflect the publiccode.PublicCode structure.
		mapping = ""
	)

	// Starting with elastic.v5, you must pass a context to execute each service.
	ctx := context.Background()

	// Add a document to the index.
	file := File{Source: domain.Id, Name: name, Data: string(data)}

	// Put publiccode data in ES.
	put, err := elasticClient.Index().
		Index(index).
		Type("doc").
		Id(domain.Id + "/" + name + "_" + index).
		BodyJson(file).
		Do(ctx)
	if err != nil {
		log.Error(err)
	}

	metrics.GetCounter("repository_file_indexed", index).Inc()

	log.Debugf("Indexed file %s to index %s, type %s", put.Id, put.Index, put.Type)
}
