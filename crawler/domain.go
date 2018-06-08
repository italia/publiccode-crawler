package crawler

import (
	"errors"
	"sync"

	"gopkg.in/yaml.v2"

	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

// Domain is a single code hosting service.
type Domain struct {
	// Domains.yml data
	Host      string   `yaml:"host"`
	BasicAuth []string `yaml:"basic-auth"`
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

	return domains, err
}

// parseDomainsFile parses the domains file to build a slice of Domain.
func parseDomainsFile(data []byte) ([]Domain, error) {
	domains := []Domain{}

	// Unmarshal the yml in domains list.
	err := yaml.Unmarshal(data, &domains)
	if err != nil {
		return nil, err
	}
	return domains, err
}

func (domain Domain) processAndGetNextURL(url string, wg *sync.WaitGroup, repositories chan Repository) (string, error) {
	crawler, err := GetClientAPICrawler(domain.Host)
	if err != nil {
		return "", err
	}
	return crawler(domain, url, repositories, wg)
}

func (domain Domain) processSingleRepo(url string, repositories chan Repository) error {
	crawler, err := GetSingleClientAPICrawler(domain.Host)
	if err != nil {
		return err
	}
	return crawler(domain, url, repositories)
}

func (domain Domain) generateAPIURL(u string) (string, error) {
	crawler, err := GetAPIURL(domain.Host)
	if err != nil {
		return u, err
	}
	return crawler(u)
}

// KnownHost detect the the right Domain API from the given URL and returns it.
// If no API is recognized will return an empty domain and an error.
func KnownHost(link, host string, domains []Domain) (Domain, error) {
	for _, domain := range domains {
		if host == domain.Host {
			// Host is found in the host list.
			return domain, nil
		}
	}

	// host unknown, needs to be inferred.
	if IsGithub(link) {
		log.Infof("%s - API inferred:%s", link, "github.com")
		return Domain{Host: "github.com"}, nil
	} else if IsBitbucket(link) {
		log.Infof("%s - API inferred:%s", link, "bitbucket.org")
		return Domain{Host: "bitbucket.org"}, nil
	} else if IsGitlab(link) {
		log.Infof("%s - API inferred:%s", link, "gitlab.com")
		return Domain{Host: "gitlab.com"}, nil
	}

	return Domain{}, errors.New("this host is not registered: " + host)
}
