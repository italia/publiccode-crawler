package crawler

import (
	"gopkg.in/yaml.v2"

	"github.com/go-redis/redis"
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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Unmarshal the yml in hostings list
	err := yaml.Unmarshal(data, &hostings)
	if err != nil {
		return nil, err
	}

	for i, hosting := range hostings {
		switch hosting.ServiceName {
		case "bitbucket":
			// AAA: se in redis c'è un url che è "false" => set hosting.URL to the one that is "false"
			keys, _ := redisClient.Keys("*").Result()
			for _, key := range keys {
				if redisClient.Get(key).Val() == "false" {
					log.Debug("Found one false URL! start from here: " + key)
					hostings[i].ServiceInstance = Bitbucket{
						URL:       key,
						RateLimit: hostings[i].RateLimit,
					}
					break
					// Altrimenti usa ciò che legge dal file
				} else {
					hostings[i].ServiceInstance = Bitbucket{
						URL:       hosting.URL,
						RateLimit: hostings[i].RateLimit,
					}
				}
			}
			break
		default:
			log.Warningf("implementation not found for service %s, skipping", hosting.ServiceName)
			break
		}
	}

	return hostings, nil
}
