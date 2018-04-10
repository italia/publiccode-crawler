package crawler

import log "github.com/sirupsen/logrus"

// Crawler is the interface for every specific crawler instances.
type Crawler interface {
	GetRepositories(repositories chan Repository) error
}

// Process delegates the work to single hosting crawlers.
func Process(hosting Hosting, repositories chan Repository) {
	if hosting.ServiceInstance == nil {
		return
	}

	err := hosting.ServiceInstance.GetRepositories(repositories)
	if err != nil {
		log.Errorf("error reading %s repository list: %v", hosting.ServiceName, err)
		return
	}
}
