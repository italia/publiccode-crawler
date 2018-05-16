package crawler

import (
	"os"
	"time"

	"net/http"
	"strings"

	"github.com/italia/developers-italia-backend/httpclient"
	"github.com/italia/developers-italia-backend/metrics"

	"github.com/italia/developers-italia-backend/persistency"
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
	fileTimestamp := time.Now().Unix() // Add unique timestamp to repo retrieve.

	// Init Prometheus for metrics.
	processedCounter := metrics.PrometheusCounter("repository_processed", "Number of repository processed.")

	for repository := range repositories {
		go checkAvailability(repository, fileTimestamp, processedCounter)
	}
}

func checkAvailability(repository Repository, fileTimestamp int64, processedCounter prometheus.Counter) {
	name := repository.Name
	fileRawUrl := repository.FileRawURL
	domain := repository.Domain
	headers := repository.Headers

	processedCounter.Inc()

	resp, err := httpclient.GetURL(fileRawUrl, headers)
	// If it's available and no error returned.
	if resp.Status.Code == http.StatusOK && err == nil {

		// Save to file.
		persistency.SaveToFile(domain, name, resp.Body, fileTimestamp)

		// Save to ES.
		persistency.SaveToES(domain, name, resp.Body, fileTimestamp)

		// Validate file.
		// TODO: uncomment these lines when mapping and File structure are ready for publiccode.
		// TODO: now validation is ulesess because we test on .gitignore file.
		// err := validateRemoteFile(resp.Body, fileRawUrl)
		// if err != nil {
		// 	log.Warn("Validator fails for: " + fileRawUrl)
		// 	log.Warn("Validator errors:" + err.Error())
		// }
	}
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
