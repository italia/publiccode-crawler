package crawler

import (
	"gopkg.in/yaml.v2"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"errors"
	"fmt"
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
	Name       string
	FileRawURL string
	Domain     string
	Headers    map[string]string
}

func ReadAndParseHosting() ([]Hosting, error) {
	// Open and read hosting file list.
	hostingFile := "hosting.yml"
	data, err := ioutil.ReadFile(hostingFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in reading %s file: %v", hostingFile, err))
	}
	// Parse hosting file list.
	hostings, err := parseHostingFile(data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in parsing %s file: %v", hostingFile, err))
	}
	log.Info("Loaded and parsed hosting.yml")

	return hostings, nil
}


// ParseHostingFile parses the hosting file to build a slice of Hosting.
func parseHostingFile(data []byte) ([]Hosting, error) {
	hostings := []Hosting{}

	// Unmarshal the yml in hostings list.
	err := yaml.Unmarshal(data, &hostings)
	if err != nil {
		return nil, err
	}
	// Redis connection.
	redisClient, err := redisClientFactory("localhost:6379")
	if err != nil {
		return hostings, err
	}

	// Manage every host
	for i, hosting := range hostings {
		// Switch over hostings.
		switch hosting.ServiceName {

		case "bitbucket":
			// Check if there is some failed URL in redis.
			data, err := checkFailedBitbucket(hostings[i], redisClient)
			if err != nil {
				log.Warn(err)
			}

			hostings[i].ServiceInstance = data
			hostings[i].URL = data.URL

		case "github":
			// Check if there is some failed URL in redis.
			data, err := checkFailedGithub(hostings[i], redisClient)
			if err != nil {
				log.Warn(err)
			}

			hostings[i].ServiceInstance = data
			hostings[i].URL = data.URL

		case "gitlab":
			// Check if there is some failed URL in redis.
			data, err := checkFailedGitlab(hostings[i], redisClient)
			if err != nil {
				log.Warn(err)
			}

			hostings[i].ServiceInstance = data
			hostings[i].URL = data.URL

		default:
			log.Warningf("implementation not found for service %s, skipping", hosting.ServiceName)
		}
	}

	return hostings, nil
}

// checkFailedBitbucket checks if a repository list previously failed to be retrieved in bitbucket.
func checkFailedBitbucket(hosting Hosting, redisClient *redis.Client) (Bitbucket, error) {

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
			log.Debugf("Found one interrupted URL. Starts from here: %s", key)
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

// checkFailedGithub checks if a repository list previously failed to be retrieved in github.
func checkFailedGithub(hosting Hosting, redisClient *redis.Client) (Github, error) {

	// Check if there is an URL that wasn't correctly retrieved.
	// URL.value="false" => set hosting.URL to the one that one ("false")
	keys, _ := redisClient.HKeys(hosting.ServiceName).Result()

	// First launch.
	if len(keys) == 0 {
		return Github{
			URL:       hosting.URL,
			RateLimit: hosting.RateLimit,
			BasicAuth: hosting.BasicAuth,
		}, nil
	}

	// N launch. Check if some repo list was interrupted.
	for _, key := range keys {
		if redisClient.HGet(hosting.ServiceName, key).Val() == "failed" {
			log.Debugf("Found one interrupted URL. Starts from here: %s", key)
			return Github{
				URL:       key,
				RateLimit: hosting.RateLimit,
				BasicAuth: hosting.BasicAuth,
			}, nil
		}
	}

	return Github{
		URL:       hosting.URL,
		RateLimit: hosting.RateLimit,
		BasicAuth: hosting.BasicAuth,
	}, nil
}

// checkFailedGitlab checks if a repository list previously failed to be retrieved in gitlab.
func checkFailedGitlab(hosting Hosting, redisClient *redis.Client) (Gitlab, error) {

	// Check if there is an URL that wasn't correctly retrieved.
	// URL.value="false" => set hosting.URL to the one that one ("false")
	keys, _ := redisClient.HKeys(hosting.ServiceName).Result()

	// First launch.
	if len(keys) == 0 {
		return Gitlab{
			URL:       hosting.URL,
			RateLimit: hosting.RateLimit,
			BasicAuth: hosting.BasicAuth,
		}, nil
	}

	// N launch. Check if some repo list was interrupted.
	for _, key := range keys {
		if redisClient.HGet(hosting.ServiceName, key).Val() == "failed" {
			log.Debugf("Found one interrupted URL. Starts from here: %s", key)
			return Gitlab{
				URL:       key,
				RateLimit: hosting.RateLimit,
				BasicAuth: hosting.BasicAuth,
			}, nil
		}
	}

	return Gitlab{
		URL:       hosting.URL,
		RateLimit: hosting.RateLimit,
		BasicAuth: hosting.BasicAuth,
	}, nil
}
