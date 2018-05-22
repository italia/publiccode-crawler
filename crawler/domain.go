package crawler

import (
	"sync"

	"gopkg.in/yaml.v2"

	"errors"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

// Domain is a single code hosting service.
type Domain struct {
	// Domains.yml data
	Id          string   `yaml:"id"`
	Description string   `yaml:"description"`
	ClientApi   string   `yaml:"client-api"`
	RawBaseUrl  string   `yaml:"rawBaseUrl"`
	BasicAuth   []string `yaml:"basic-auth"`
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

	return domains, nil
}

// // updateStartURL checks if a repository list previously failed to be retrieved.
// func (domain *Domain) updateDomainState(redisClient *redis.Client) error {
// 	// Check if there is an URL that wasn't correctly retrieved.
// 	// URL.value="failed" => set domain.URL to that one
// 	keys, err := redisClient.HKeys(domain.Id).Result()
// 	if err != nil {
// 		return err
// 	}
//
// 	// N launch. Check if some repo list was interrupted.
// 	for _, key := range keys {
// 		if redisClient.HGet(domain.Id, key).Val() != "" {
// 			log.Debugf("Found one interrupted URL. Starts from here: %s with Index: %s", key, redisClient.HGet(domain.Id, key).Val())
// 			domain.URL = key
// 		}
// 	}
//
// 	return nil
// }

func (domain Domain) processAndGetNextURL(url string, wg *sync.WaitGroup, repositories chan Repository) (string, error) {
	crawler, err := GetClientApiCrawler(domain.ClientApi)
	if err != nil {
		return "", err
	}

	return crawler(domain, url, repositories, wg)
}

// func (domain Domain) processSingleRepo(url string, repositories chan Repository) error {
// 	crawler, err := GetSingleClientApiCrawler(domain.ClientApi)
// 	if err != nil {
// 		return err
// 	}
//
// 	return crawler(domain, url, repositories)
// }
