package crawler

import (
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

// Crawler is the interface for every specific crawler instances.
type Crawler interface {
	GetRepositories(repositories chan Repository) error
}

// Hosting is a single hosting service.
type Hosting struct {
	URL             string `yaml:"url"`
	ServiceName     string `yaml:"service"`
	ServiceInstance Crawler
}

// Repository is a single code repository.
type Repository struct {
	Name string
	URL  string
}

// ParseHostingFile parses the hosting file to build a slice of Hosting.
func ParseHostingFile(data []byte) ([]Hosting, error) {
	var rows = make([]Hosting, 3)

	err := yaml.Unmarshal(data, &rows)
	if err != nil {
		return nil, err
	}

	for i, row := range rows {
		switch row.ServiceName {
		case "bitbucket":
			rows[i].ServiceInstance = Bitbucket{
				URL: row.URL,
			}
			break
		default:
			log.Warningf("implementation not found for service %s, skipping", row.ServiceName)
			break
		}
	}

	return rows, nil
}

// Process delegates the work to single hosting crawlers.
func Process(hosting Hosting, repositories chan Repository) {
	if hosting.ServiceInstance == nil {
		return
	}

	err := hosting.ServiceInstance.GetRepositories(repositories)
	if err != nil {
		log.Errorf("error reading %s repository list: %v", hosting.ServiceName, err)
	}
}
