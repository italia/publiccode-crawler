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
	BasicAuth string `yaml:"basic-auth"`

	ServiceInstance Crawler
}

// Repository is a single code repository.
type Repository struct {
	Name    string
	URL     string
	Source  string
	Headers map[string]string
}

// ParseHostingFile parses the hosting file to build a slice of Hosting.
func ParseHostingFile(data []byte) ([]Hosting, error) {
	hostings := []Hosting{}

	// Unmarshal the yml in hostings list.
	err := yaml.Unmarshal(data, &hostings)
	if err != nil {
		return nil, err
	}

	// Manage every host
	for i, hosting := range hostings {
		// Switch over hostings.
		switch hosting.ServiceName {

		case "bitbucket":
			// Check if there is some work failed.
			data, err := checkFailed(hostings[i])
			if err != nil {
				log.Warn(err)
			}

			hostings[i].ServiceInstance = data
			hostings[i].URL = data.URL
			break
		default:
			log.Warningf("implementation not found for service %s, skipping", hosting.ServiceName)
			break
		}
	}

	return hostings, nil
}

func checkFailed(hosting Hosting) (Bitbucket, error) {
	// Redis connection.
	redisClient, err := redisClientFactory("redis:6379")
	if err != nil {
		return Bitbucket{
			URL:       hosting.URL,
			RateLimit: hosting.RateLimit,
			BasicAuth: hosting.BasicAuth,
		}, err
	}

	// Check if there is an URL that wasn't correctly retrieved.
	// URL.value="false" => set hosting.URL to the one that one ("false")
	keys, _ := redisClient.HKeys(hosting.ServiceName).Result()

	// First launch.
	if len(keys) == 0 {
		return Bitbucket{
			URL:       hosting.URL,
			RateLimit: hosting.RateLimit,
			BasicAuth: hosting.BasicAuth,
		}, nil

	}

	// N launch. Check if some repo list was interrupted.
	for _, key := range keys {

		if redisClient.HGet(hosting.ServiceName, key).Val() == "failed" {
			log.Debug("Found one interrupted URL. Starts from here: " + key)
			return Bitbucket{
				URL:       key,
				RateLimit: hosting.RateLimit,
				BasicAuth: hosting.BasicAuth,
			}, nil

		}
	}

	return Bitbucket{
		URL:       hosting.URL,
		RateLimit: hosting.RateLimit,
		BasicAuth: hosting.BasicAuth,
	}, nil
}
