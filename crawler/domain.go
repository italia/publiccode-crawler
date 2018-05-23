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
	ApiBaseUrl  string   `yaml:"apiBaseUrl"`
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
