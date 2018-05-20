package crawler

import (
	"context"
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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Repository is a single code repository.
type Repository struct {
	Name       string
	FileRawURL string
	Domain     Domain
	Headers    map[string]string
}

type Handler func(domain Domain, url string, repositories chan Repository, wg *sync.WaitGroup) (string, error)

func AddIndex(index string, elasticClient *elastic.Client) {
	// Use the IndexExists service to check if a specified index exists.
	exists, err := elasticClient.IndexExists(index).Do(context.Background())
	if err != nil {
		log.Error(err)
	}
	if !exists {
		// Create a new index.
		// TODO: When mapping will be available: client.CreateIndex(index).BodyString(mapping).Do(ctx).
		_, err = elasticClient.CreateIndex(index).Do(context.Background())
		if err != nil {
			log.Error(err)
		}
	}
}

// UpdateIndex return an old index if found on redis. A new index if no one is found.
func UpdateIndex(domains []Domain, redisClient *redis.Client, elasticClient *elastic.Client) (string, error) {
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
				AddIndex(Index, elasticClient)
				return Index, nil
			}
		}
	}

	// If reached there will be a new index.
	Index = strconv.FormatInt(time.Now().Unix(), 10)
	AddIndex(Index, elasticClient)
	return Index, nil
}

// Process delegates the work to single domain crawlers.
func ProcessDomain(domain Domain, repositories chan Repository, wg *sync.WaitGroup) {
	// Redis connection.
	redisClient, err := RedisClientFactory(viper.GetString("REDIS_URL"))
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
	metrics.RegisterPrometheusCounter("repository_processed", "Number of repository processed.")
	metrics.RegisterPrometheusCounter("repository_file_saved", "Number of file saved.")
	metrics.RegisterPrometheusCounter("repository_file_saved_valid", "Number of valid file saved.")

	for repository := range repositories {
		wg.Add(1)
		go checkAvailability(repository, index, wg, elasticClient)
	}

}

func checkAvailability(repository Repository, index string, wg *sync.WaitGroup, elasticClient *elastic.Client) {
	name := repository.Name
	fileRawUrl := repository.FileRawURL
	domain := repository.Domain
	headers := repository.Headers

	metrics.GetCounter("repository_processed").Inc()
	metrics.GetCounter(repository.Domain.Id).Inc()

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
	fileName := viper.GetString("CRAWLED_FILENAME")
	// Parse data into pc struct and validate.
	baseURL := strings.TrimSuffix(url, fileName)
	// Set remore URL for remote validation (it will check files availability).
	publiccode.BaseDir = baseURL
	var pc publiccode.PublicCode

	err := publiccode.Parse(data, &pc)

	if err != nil {
		return err
	}

	metrics.GetCounter("repository_file_saved_valid").Inc()
	return err

}

// WaitingLoop waits until all the goroutines counter is zero and close the repositories channel.
func WaitingLoop(repositories chan Repository, index string, wg *sync.WaitGroup, elasticClient *elastic.Client) {
	wg.Wait()

	// Remove old aliases.
	res, err := elasticClient.Aliases().Index("_all").Do(context.Background())
	if err != nil {
		panic(err)
	}
	aliasService := elasticClient.Alias()
	indices := res.IndicesByAlias("publiccode")
	for _, name := range indices {
		log.Debugf("Remove alias from %s to %s", "publiccode", name)
		aliasService.Remove(name, "publiccode").Do(context.Background())
	}

	// Add an alias to the new index.
	log.Debugf("Add alias from %s to %s", index, "publiccode")
	aliasService.Add(index, "publiccode").Do(context.Background())

	close(repositories)
}
