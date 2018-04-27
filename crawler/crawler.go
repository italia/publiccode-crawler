package crawler

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// Crawler is the interface for every specific crawler instances.
type Crawler interface {
	GetRepositories(url string, repositories chan Repository) (string, error)
}

// Process delegates the work to single hosting crawlers.
func Process(hosting Hosting, repositories chan Repository) {
	if hosting.ServiceInstance == nil {
		log.Warnf("Hosting %s is not available.", hosting.ServiceName)
		return
	}

	// Redis connection.
	redisClient, err := redisClientFactory(os.Getenv("REDIS_URL"))
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
