package persistency

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

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
func SaveToES(source, name string, data []byte, fileTimestamp int64) {
	const (
		index = "publiccode" // Elasticsearch index.
		// Elasticsearch mapping for publiccode. Checkout elasticsearch/mappings/software.json
		// TODO: Mapping must reflect the publiccode.PublicCode structure.
		mapping = ""
	)
	fileTime := strconv.FormatInt(fileTimestamp, 10)

	// Starting with elastic.v5, you must pass a context to execute each service.
	ctx := context.Background()

	// Create a client.
	client, err := elastic.NewClient(
		elastic.SetURL(os.Getenv("ELASTIC_URL")),
		elastic.SetRetrier(NewESRetrier()),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(os.Getenv("ELASTIC_USER"), os.Getenv("ELASTIC_PWD")))
	if err != nil {
		panic(err)
	}
	if elastic.IsConnErr(err) {
		log.Error("Elasticsearch connection problem: %v", err)
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
	file := File{Source: source, Name: name, Data: string(data)}

	// Put publiccode data in ES.
	put, err := client.Index().
		Index(index).
		Type("doc").
		Id(source + "/" + name + "_" + fileTime).
		BodyJson(file).
		Do(ctx)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("Indexed file %s to index %s, type %s", put.Id, put.Index, put.Type)

}

type ElasticRetrier struct {
	backoff elastic.Backoff
}

func NewESRetrier() *ElasticRetrier {
	return &ElasticRetrier{
		backoff: elastic.NewExponentialBackoff(10*time.Millisecond, 8*time.Second),
	}
}

func (r *ElasticRetrier) Retry(ctx context.Context, retry int, req *http.Request, resp *http.Response, err error) (time.Duration, bool, error) {
	log.Warn("Elasticsearch connection problem. Retry.")

	// Stop after 8 retries: 2m
	if retry >= 8 {
		return 0, false, errors.New("Elasticsearch or network down")
	}

	// Let the backoff strategy decide how long to wait and whether to stop
	wait, stop := r.backoff.Next(retry)
	return wait, stop, nil
}
