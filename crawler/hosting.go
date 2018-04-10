package crawler

import (
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

// Hosting is a single hosting service.
type Hosting struct {
	ServiceName string `yaml:"name"`
	URL         string `yaml:"url"`
	RateLimit   struct {
		ReqH int `yaml:"req/h"`
		ReqM int `yaml:"req/m"`
	} `yaml:"rate-limit"`

	ServiceInstance Crawler
}

// Repository is a single code repository.
type Repository struct {
	Name string
	URL  string
}

// ParseHostingFile parses the hosting file to build a slice of Hosting.
func ParseHostingFile(data []byte) ([]Hosting, error) {
	hostings := []Hosting{}

	// Unmarshal the yml in hostings list
	err := yaml.Unmarshal(data, &hostings)
	if err != nil {
		return nil, err
	}

	for i, hosting := range hostings {
		switch hosting.ServiceName {
		case "bitbucket":
			hostings[i].ServiceInstance = Bitbucket{
				URL:       hosting.URL,
				RateLimit: hostings[i].RateLimit,
			}
			break
		default:
			log.Warningf("implementation not found for service %s, skipping", hosting.ServiceName)
			break
		}
	}

	return hostings, nil
}
