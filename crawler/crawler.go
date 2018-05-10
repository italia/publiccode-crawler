package crawler

import (
	"context"
	"os"

	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/italia/developers-italia-backend/httpclient"
	"github.com/italia/developers-italia-backend/metrics"
	elastic "github.com/olivere/elastic"

	"github.com/italia/developers-italia-backend/publiccode"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Repository is a single code repository.
type Repository struct {
	Name       string
	FileRawURL string
	Domain     string
	Headers    map[string]string
}

type File struct {
	Source string `json:"source"`
	Name   string `json:"name"`
	Data   string `json:"data"`
}

type Handler func(domain Domain, url string, repositories chan Repository) (string, error)

// Process delegates the work to single domain crawlers.
func ProcessDomain(domain Domain, repositories chan Repository) {
	// Redis connection.
	redisClient, err := RedisClientFactory(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Error(err)
	}

	// Base starting URL.
	url := domain.URL
	for {
		// Set the value of nextURL on redis to "failed".
		err = redisClient.HSet(domain.Id, url, "failed").Err()
		if err != nil {
			log.Error(err)
		}

		nextURL, err := domain.processAndGetNextURL(url, repositories)
		if err != nil {
			log.Errorf("error reading %s repository list: %v. NextUrl: %v", url, err, nextURL)
			log.Errorf("Retry: %s", nextURL)
			nextURL = url
			//close(repositories): ok if only one repo. If more parallel it generates panics.
			//return
		}
		// If reached, the repository list was successfully retrieved.
		// Delete the repository url from redis.
		err = redisClient.HDel(domain.Id, url).Err()
		if err != nil {
			log.Error(err)
		}

		// If end is reached, url and nextURL contains the same value.
		if nextURL == "" {
			log.Infof("Url: %s - is the last one.", url)
			return
		}
		// Update url to nextURL.
		url = nextURL
	}
}

func ProcessRepositories(repositories chan Repository) {
	log.Debug("Repositories are going to be processed...")

	// Init Prometheus for metrics.
	processedCounter := metrics.PrometheusCounter("repository_processed", "Number of repository processed.")

	// Throttle requests.
	// Time limits should be calibrated on more tests in order to avoid errors and bans.
	throttleRate := time.Second / 1000
	throttle := time.Tick(throttleRate)

	for repository := range repositories {
		// Throttle down the calls.
		<-throttle
		go checkAvailability(repository, processedCounter)
	}
}

func checkAvailability(repository Repository, processedCounter prometheus.Counter) {
	name := repository.Name
	fileRawUrl := repository.FileRawURL
	domain := repository.Domain
	headers := repository.Headers

	processedCounter.Inc()

	resp, err := httpclient.GetURL(fileRawUrl, headers)
	// If it's available and no error returned.
	if resp.Status.Code == http.StatusOK && err == nil {
		// Save the file.
		saveFile(domain, name, resp.Body)

		// Save on ES
		saveES(domain, name, resp.Body)

		// Validate file.
		// err := validateRemoteFile(resp.Body, fileRawUrl)
		// if err != nil {
		// 	log.Warn("Validator fails for: " + fileRawUrl)
		// 	log.Warn("Validator errors:" + err.Error())
		// }
	}
}

// saveFile save the chosen <file_name> in ./data/<source>/<vendor>/<repo>/<file_name>
func saveFile(source, name string, data []byte) {
	fileName := os.Getenv("CRAWLED_FILENAME")
	vendor, repo := splitFullName(name)

	path := filepath.Join("./data", source, vendor, repo)

	// MkdirAll will create all the folder path, if not exists.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	err := ioutil.WriteFile(filepath.Join(path, fileName), data, 0644)
	if err != nil {
		log.Error(err)
	}
}

// splitFullName split a git FullName format to vendor and repo strings.
func splitFullName(fullName string) (string, string) {
	s := strings.Split(fullName, "/")
	return s[0], s[1]
}

// validateRemoteFile save the chosen <file_name> in ./data/<source>/<vendor>/<repo>/<file_name>
func validateRemoteFile(data []byte, url string) error {
	fileName := os.Getenv("CRAWLED_FILENAME")
	// Parse data into pc struct and validate.
	baseURL := strings.TrimSuffix(url, fileName)
	// Set remore URL for remote validation (it will check files availability).
	publiccode.BaseDir = baseURL
	var pc publiccode.PublicCode

	return publiccode.Parse(data, &pc)
}

// saveES save the chosen <file_name> in elasticsearch
func saveES(source, name string, data []byte) {
	index := os.Getenv("CRAWLED_FILENAME")
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()
	// Create a client
	client, err := elastic.NewClient(
		elastic.SetURL(os.Getenv("ELASTIC_URL")),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(os.Getenv("ELASTIC_USER"), os.Getenv("ELASTIC_PWD")))
	if err != nil {
		log.Error(err)
	}
	if elastic.IsConnErr(err) {
		log.Error("Elasticsearch connection problem: %v", err)
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		// Handle error
		log.Error(err)
	}

	if !exists {
		// Create a new index
		_, err = client.CreateIndex(index).Do(ctx)
		if err != nil {
			log.Error(err)
		}
	}
	// Add a document to the index
	file := File{Source: source, Name: name, Data: string(data)}

	put, err := client.Index().
		Index(index).
		Type("doc").
		Id(source + "/" + name).
		BodyJson(file).
		Do(ctx)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("Indexed file %s to index %s, type %s", put.Id, put.Index, put.Type)

}
