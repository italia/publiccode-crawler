package crawler

import (
	"sync"

	"gopkg.in/yaml.v2"

	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

// Domain is a single code hosting service.
type Domain struct {
	// Domains.yml data
	ID          string   `yaml:"id"`
	Description string   `yaml:"description"`
	ClientAPI   string   `yaml:"client-api"`
	APIOrgURL   string   `yaml:"apiOrgUrl"`
	APIRepoURL  string   `yaml:"apiRepoUrl"`
	RawBaseURL  string   `yaml:"rawBaseUrl"`
	BasicAuth   []string `yaml:"basic-auth"`
}

// ReadAndParseDomains read domainsFile and return the parsed content in a Domain slice.
func ReadAndParseDomains(domainsFile string) ([]Domain, error) {
	// Open and read domains file list.
	data, err := ioutil.ReadFile(domainsFile)
	if err != nil {
		return nil, fmt.Errorf("error in reading %s file: %v", domainsFile, err)
	}
	// Parse domains file list.
	domains, err := parseDomainsFile(data)
	if err != nil {
		return nil, fmt.Errorf("error in parsing %s file: %v", domainsFile, err)
	}
	log.Infof("Loaded and parsed %s", domainsFile)

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
	crawler, err := GetClientAPICrawler(domain.ClientAPI)
	if err != nil {
		return "", err
	}

	return crawler(domain, url, repositories, wg)
}

func (domain Domain) processSingleRepo(url string, repositories chan Repository) error {
	crawler, err := GetSingleClientAPICrawler(domain.ClientAPI)
	if err != nil {
		return err
	}

	return crawler(domain, url, repositories)
}
