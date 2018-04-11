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
	Name   string
	URL    string
	Source string
}

// ParseHostingFile parses the hosting file to build a slice of Hosting.
func ParseHostingFile(data []byte) ([]Hosting, error) {
	hostings := []Hosting{}

	// Redis connection.
	redisClient, err := redisClientFactory("redis:6379")
	if err != nil {
		return nil, err
	}

	// Unmarshal the yml in hostings list.
	err = yaml.Unmarshal(data, &hostings)
	if err != nil {
		return nil, err
	}

	// Manage every host
	for i, hosting := range hostings {
		switch hosting.ServiceName {
		case "bitbucket":
			// Check if there is an URL that wasn't correctly retrieved.
			// URL.value="false" => set hosting.URL to the one that one ("false")
			keys, _ := redisClient.Keys("*").Result()
			// First launch. TODO: refactory. This break is terrible.
			if len(keys) == 0 {
				hostings[i].ServiceInstance = Bitbucket{
					URL:       hosting.URL,
					RateLimit: hostings[i].RateLimit,
				}
				break
			}
			for _, key := range keys {
				if redisClient.Get(key).Val() == "false" {
					log.Debug("Found one false URL! start from here: " + key)
					hostings[i].ServiceInstance = Bitbucket{
						URL:       key,
						RateLimit: hostings[i].RateLimit,
					}
					break
					// Or read from file.
				} else {
					hostings[i].ServiceInstance = Bitbucket{
						URL:       hosting.URL,
						RateLimit: hostings[i].RateLimit,
					}
				}
			}
			break
		case "github":
			log.Warningf("implementation not found for service %s, skipping", hosting.ServiceName)
			break
		case "gitlab":
			log.Warningf("implementation not found for service %s, skipping", hosting.ServiceName)
			break
		default:
			log.Warningf("implementation not found for service %s, skipping", hosting.ServiceName)
			break
		}
	}

	return hostings, nil
}
