package crawler

import (
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// Crawler is the interface for every specific crawler instances.
type Crawler interface {
	GetRepositories(repositories chan Repository) (string, error)
}

// Process delegates the work to single hosting crawlers.
func Process(hosting Hosting, repositories chan Repository) {
	if hosting.ServiceInstance == nil {
		log.Warnf("Hosting %s is not available.", hosting.ServiceName)
		return
	}

	sourceURL, err := hosting.ServiceInstance.GetRepositories(repositories)
	if err != nil {
		// Redis connection.
		redisClient, err := redisClientFactory("redis:6379")
		if err != nil {
			log.Error(err)
		}

		// If reached, the repository list was not retrieved.
		// Set the value of sourceURL on redis to false.
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		err = redisClient.HSet(hosting.ServiceName, sourceURL, false).Err()
		if err != nil {
			log.Error(err)
		}
		log.Errorf("error reading %s repository list: %v", hosting.ServiceName, err)
		log.Info("source url saved on REDIS at time: " + timestamp)
		close(repositories)
		return
	}
}
