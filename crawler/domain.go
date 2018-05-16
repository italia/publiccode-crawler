package crawler

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"gopkg.in/yaml.v2"

	"errors"
	"fmt"
	"io/ioutil"

	"github.com/go-redis/redis"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

// Domain is a single code hosting service.
type Domain struct {
	// Domains.yml data
	Id          string `yaml:"id"`
	Description string `yaml:"description"`
	ClientApi   string `yaml:"client-api"`
	URL         string `yaml:"url"`
	RateLimit   struct {
		ReqH int `yaml:"req/h"`
		ReqM int `yaml:"req/m"`
	} `yaml:"rate-limit"`
	BasicAuth []string `yaml:"basic-auth"`
	// Specific domain options.
	Index string // Index define the current crawler execution.
}

func ReadAndParseDomains(domainsFile string, redisClient *redis.Client, elasticClient *elastic.Client) ([]Domain, error) {
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

	// Update the start URL if a failed one found in Redis.
	for i, _ := range domains {
		domains[i].updateDomainState(redisClient)
		domains[i].addAlias(elasticClient)
	}

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

	return domains, nil
}

// updateStartURL checks if a repository list previously failed to be retrieved.
func (domain *Domain) updateDomainState(redisClient *redis.Client) error {
	// Set initial domain Index.
	domain.Index = strconv.FormatInt(time.Now().Unix(), 10)

	// Check if there is an URL that wasn't correctly retrieved.
	// URL.value="failed" => set domain.URL to that one
	keys, err := redisClient.HKeys(domain.Id).Result()
	if err != nil {
		return err
	}

	// N launch. Check if some repo list was interrupted.
	for _, key := range keys {
		if redisClient.HGet(domain.Id, key).Val() != "" {
			log.Debugf("Found one interrupted URL. Starts from here: %s with Index: %s", key, redisClient.HGet(domain.Id, key).Val())
			domain.URL = key
			domain.Index = redisClient.HGet(domain.Id, key).Val()
		}
	}

	return nil
}

// addAlias adds alias to elastic.
func (domain *Domain) addAlias(elasticClient *elastic.Client) error {
	// Create a client.
	client, err := ElasticClientFactory(
		os.Getenv("ELASTIC_URL"),
		os.Getenv("ELASTIC_USER"),
		os.Getenv("ELASTIC_PWD"))
	if err != nil {
		return err
	}

	// Add alias to publiccode.
	client.Alias().Add("publiccode", domain.Index).Do(context.Background())

	return err
}

func (domain Domain) processAndGetNextURL(url string, wg *sync.WaitGroup, repositories chan Repository) (string, error) {
	crawler, err := GetClientApiCrawler(domain.ClientApi)
	if err != nil {
		return "", err
	}

	return crawler(domain, url, repositories, wg)
}
