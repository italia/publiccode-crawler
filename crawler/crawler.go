package crawler

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"net/http"
	"strings"

	"github.com/go-redis/redis"
	"github.com/italia/developers-italia-backend/httpclient"
	"github.com/italia/developers-italia-backend/metrics"
	"github.com/olivere/elastic"

	"github.com/italia/developers-italia-backend/publiccode"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Repository is a single code repository.
type Repository struct {
	Name       string
	FileRawURL string
	Domain     Domain
	Headers    map[string]string
}

type Handler func(domain Domain, url string, repositories chan Repository, wg *sync.WaitGroup) (string, error)

// UpdateIndex return an old index if found on redis. A new index if no one is found.
func UpdateIndex(domains []Domain, redisClient *redis.Client) (string, error) {
	for i, _ := range domains {
		// Check if there is an URL that wasn't correctly retrieved.
		keys, err := redisClient.HKeys(domains[i].Id).Result()
		if err != nil {
			return "", err
		}
		// If some repository is left the Idex sould remain the last one saved.
		for _, key := range keys {
			if redisClient.HGet(domains[i].Id, key).Val() != "" {
				Index = redisClient.HGet(domains[i].Id, key).Val()
				return Index, nil
			}
		}
	}

	// If reached there will be a new index.
	Index = strconv.FormatInt(time.Now().Unix(), 10)
	return Index, nil
}

// Process delegates the work to single domain crawlers.
func ProcessDomain(domain Domain, repositories chan Repository, wg *sync.WaitGroup) {
	// Redis connection.
	redisClient, err := RedisClientFactory(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Error(err)
	}

	// Base starting URL.
	url := domain.URL
	for {
		// Set the value of nextURL on redis to domain.Index that describe the current execution.
		err = redisClient.HSet(domain.Id, url, Index).Err()
		if err != nil {
			log.Error(err)
		}

		nextURL, err := domain.processAndGetNextURL(url, wg, repositories)
		if err != nil {
			log.Errorf("error reading %s repository list: %v. NextUrl: %v", url, err, nextURL)
			log.Errorf("Retry: %s", nextURL)
			nextURL = url
		}
		// If reached, the repository list was successfully retrieved.
		// Delete the repository url from redis.
		err = redisClient.HDel(domain.Id, url).Err()
		if err != nil {
			log.Error(err)
		}

		// If end is reached, nextUrl is empty.
		if nextURL == "" {
			log.Infof("Url: %s - is the last one.", url)

			// WaitingGroupd
			wg.Done()
			return
		}
		// Update url to nextURL.
		url = nextURL
	}
}

func ProcessRepositories(repositories chan Repository, index string, wg *sync.WaitGroup, elasticClient *elastic.Client) {
	log.Debug("Repositories are going to be processed...")
	// Init Prometheus for metrics.
	processedCounter := metrics.PrometheusCounter("repository_processed", "Number of repository processed.")

	for repository := range repositories {
		wg.Add(1)
		go checkAvailability(repository, index, wg, processedCounter, elasticClient)
	}

}

func checkAvailability(repository Repository, index string, wg *sync.WaitGroup, processedCounter prometheus.Counter, elasticClient *elastic.Client) {
	name := repository.Name
	fileRawUrl := repository.FileRawURL
	domain := repository.Domain
	headers := repository.Headers

	processedCounter.Inc()

	resp, err := httpclient.GetURL(fileRawUrl, headers)
	// If it's available and no error returned.
	if resp.Status.Code == http.StatusOK && err == nil {

		// Save to file.
		SaveToFile(domain, name, resp.Body, index)

		// Save to ES.
		SaveToES(domain, name, resp.Body, index, elasticClient)

		// Validate file.
		// TODO: uncomment these lines when mapping and File structure are ready for publiccode.
		// TODO: now validation is ulesess because we test on .gitignore file.
		// err := validateRemoteFile(resp.Body, fileRawUrl)
		// if err != nil {
		// 	log.Warn("Validator fails for: " + fileRawUrl)
		// 	log.Warn("Validator errors:" + err.Error())
		// }
	}

	// Defer waiting group close.
	wg.Done()
}

// validateRemoteFile validate the remote file
func validateRemoteFile(data []byte, url string) error {
	fileName := os.Getenv("CRAWLED_FILENAME")
	// Parse data into pc struct and validate.
	baseURL := strings.TrimSuffix(url, fileName)
	// Set remore URL for remote validation (it will check files availability).
	publiccode.BaseDir = baseURL
	var pc publiccode.PublicCode

	return publiccode.Parse(data, &pc)
}

// WaitingLoop waits until all the goroutines counter is zero and close the repositories channel.
func WaitingLoop(repositories chan Repository, index string, wg *sync.WaitGroup, elasticClient *elastic.Client) {
	wg.Wait()
	log.Debugf("Index END! %s", index)
	// WIP: if reached, all the domains are ended, so I should do the switch on Alias
	elasticClient.Alias().Remove("publiccode", index).Do(context.Background())

	close(repositories)
}
