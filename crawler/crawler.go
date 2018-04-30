package crawler

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"time"
	"github.com/italia/developers-italia-backend/httpclient"
	"net/http"
	"path/filepath"
	"io/ioutil"
	"strings"
)

// Crawler is the interface for every specific crawler instances.
type Crawler interface {
	GetRepositories(url string, repositories chan Repository) (string, error)
}

// Process delegates the work to single hosting crawlers.
func ProcessHosting(hosting Hosting, repositories chan Repository) {
	if hosting.ServiceInstance == nil {
		log.Warnf("Hosting %s is not available.", hosting.ServiceName)
		return
	}

	// Redis connection.
	redisClient, err := redisClientFactory("localhost:6379")
	if err != nil {
		log.Error(err)
	}

	// Base starting URL.
	url := hosting.URL

	for {
		// Set the value of nextURL on redis to "failed".
		err = redisClient.HSet(hosting.ServiceName, url, "failed").Err()
		if err != nil {
			log.Error(err)
		}

		nextURL, err := hosting.ServiceInstance.GetRepositories(url, repositories)
		if err != nil {
			log.Errorf("error reading %s repository list: %v. NextUrl: %v", url, err, nextURL)
			log.Errorf("Retry:", nextURL)
			nextURL = url
			//close(repositories): ok if only one repo. If more parallel it generates panics.
			//return
		}
		// If reached, the repository list was successfully retrieved.
		// Delete the repository url from redis.
		err = redisClient.HDel(hosting.ServiceName, url).Err()
		if err != nil {
			log.Error(err)
		}
		// Update url to nextURL.
		url = nextURL
	}
}

func ProcessRepositories(repositories chan Repository, processedCounter prometheus.Counter) {
	log.Debug("Repositories are going to be processed...")
	// Throttle requests.
	// Time limits should be calibrated on more tests in order to avoid errors and bans.
	throttleRate := time.Second / 1000
	throttle := time.Tick(throttleRate)

	for repository := range repositories {
		// Throttle down the calls.
		<-throttle
		go checkAvailability(repository.Name, repository.URL, repository.Source, repository.Headers, processedCounter)

	}
}

func checkAvailability(fullName, url, source string, headers map[string]string, processedCounter prometheus.Counter) {
	processedCounter.Inc()

	body, status, _, err := httpclient.GetURL(url, headers)
	// If it's available and no error returned.
	if status.StatusCode == http.StatusOK && err == nil {
		// Save the file.
		vendor, repo := splitFullName(fullName)
		fileName := os.Getenv("CRAWLED_FILENAME")
		saveFile(source, vendor, repo, fileName, body)
	}
}

// saveFile save the choosen <file_name> in ./data/<vendor>/<repo>/<file_name>
func saveFile(source, vendor, repo, fileName string, data []byte) {
	path := filepath.Join("./data", source, vendor, repo)

	// MkdirAll will create all the folder path, if not exists.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	err := ioutil.WriteFile(path+"/"+fileName, data, 0644)
	if err != nil {
		log.Error(err)
	}
}

// splitFullName split a git FullName format to vendor and repo strings.
func splitFullName(fullName string) (string, string) {
	s := strings.Split(fullName, "/")
	return s[0], s[1]
}

