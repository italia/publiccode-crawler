package crawler

import (
	"gopkg.in/yaml.v2"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"errors"
	"fmt"
	"os"
)

// Domain is a single code hosting service.
type Domain struct {
	Id          string `yaml:"id"`
	Description string `yaml:"description"`
	ClientApi   string `yaml:"client-api"`
	URL         string `yaml:"url"`
	RateLimit struct {
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

func ReadAndParseDomains(domainsFile string) ([]Domain, error) {
	// Open and read domains file list.
	data, err := ioutil.ReadFile(domainsFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in reading %s file: %v", domainsFile, err))
	}
	// Parse domains file list.
	domains, err := parseDomainsFile(data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in parsing %s file: %v", domainsFile, err))
	}
	log.Info("Loaded and parsed domains.yml")

	return domains, nil
}

// parseDomainsFile parses the domains file to build a slice of Domain.
func parseDomainsFile(data []byte) ([]Domain, error) {
	domains := []Domain{}

	// Unmarshal the yml in domains list.
	err := yaml.Unmarshal(data, &domains)
	if err != nil {
		return nil, err
	}
	// Redis connection.
	redisClient, err := redisClientFactory(os.Getenv("REDIS_URL"))
	if err != nil {
		return domains, err
	}

	// Manage every host
	for i, domain := range domains {
		// Switch over domains.
		switch domain.Id {

		case "bitbucket":
			// Check if there is some failed URL in redis.
			data, err := checkFailedBitbucket(domains[i], redisClient)
			if err != nil {
				log.Warn(err)
			}

			domains[i].ServiceInstance = data
			domains[i].URL = data.URL

		case "github":
			// Check if there is some failed URL in redis.
			data, err := checkFailedGithub(domains[i], redisClient)
			if err != nil {
				log.Warn(err)
			}

			domains[i].ServiceInstance = data
			domains[i].URL = data.URL

		case "gitlab":
			// Check if there is some failed URL in redis.
			data, err := checkFailedGitlab(domains[i], redisClient)
			if err != nil {
				log.Warn(err)
			}

			domains[i].ServiceInstance = data
			domains[i].URL = data.URL

		default:
			log.Warningf("implementation not found for service %s, skipping", domain.Id)
		}
	}

	return domains, nil
}

// checkFailedBitbucket checks if a repository list previously failed to be retrieved in bitbucket.
func checkFailedBitbucket(domain Domain, redisClient *redis.Client) (Bitbucket, error) {

	// Check if there is an URL that wasn't correctly retrieved.
	// URL.value="false" => set domain.URL to the one that one ("false")
	keys, _ := redisClient.HKeys(domain.Id).Result()

	// First launch.
	if len(keys) == 0 {
		return Bitbucket{
			URL:       domain.URL,
			RateLimit: domain.RateLimit,
			BasicAuth: domain.BasicAuth,
		}, nil
	}

	// N launch. Check if some repo list was interrupted.
	for _, key := range keys {
		if redisClient.HGet(domain.Id, key).Val() == "failed" {
			log.Debugf("Found one interrupted URL. Starts from here: %s", key)
			return Bitbucket{
				URL:       key,
				RateLimit: domain.RateLimit,
				BasicAuth: domain.BasicAuth,
			}, nil
		}
	}

	return Bitbucket{
		URL:       domain.URL,
		RateLimit: domain.RateLimit,
		BasicAuth: domain.BasicAuth,
	}, nil
}

// checkFailedGithub checks if a repository list previously failed to be retrieved in github.
func checkFailedGithub(domain Domain, redisClient *redis.Client) (Github, error) {

	// Check if there is an URL that wasn't correctly retrieved.
	// URL.value="false" => set domain.URL to the one that one ("false")
	keys, _ := redisClient.HKeys(domain.Id).Result()

	// First launch.
	if len(keys) == 0 {
		return Github{
			URL:       domain.URL,
			RateLimit: domain.RateLimit,
			BasicAuth: domain.BasicAuth,
		}, nil
	}

	// N launch. Check if some repo list was interrupted.
	for _, key := range keys {
		if redisClient.HGet(domain.Id, key).Val() == "failed" {
			log.Debugf("Found one interrupted URL. Starts from here: %s", key)
			return Github{
				URL:       key,
				RateLimit: domain.RateLimit,
				BasicAuth: domain.BasicAuth,
			}, nil
		}
	}

	return Github{
		URL:       domain.URL,
		RateLimit: domain.RateLimit,
		BasicAuth: domain.BasicAuth,
	}, nil
}

// checkFailedGitlab checks if a repository list previously failed to be retrieved in gitlab.
func checkFailedGitlab(domain Domain, redisClient *redis.Client) (Gitlab, error) {

	// Check if there is an URL that wasn't correctly retrieved.
	// URL.value="false" => set domain.URL to the one that one ("false")
	keys, _ := redisClient.HKeys(domain.Id).Result()

	// First launch.
	if len(keys) == 0 {
		return Gitlab{
			URL:       domain.URL,
			RateLimit: domain.RateLimit,
			BasicAuth: domain.BasicAuth,
		}, nil
	}

	// N launch. Check if some repo list was interrupted.
	for _, key := range keys {
		if redisClient.HGet(domain.Id, key).Val() == "failed" {
			log.Debugf("Found one interrupted URL. Starts from here: %s", key)
			return Gitlab{
				URL:       key,
				RateLimit: domain.RateLimit,
				BasicAuth: domain.BasicAuth,
			}, nil
		}
	}

	return Gitlab{
		URL:       domain.URL,
		RateLimit: domain.RateLimit,
		BasicAuth: domain.BasicAuth,
	}, nil
}
