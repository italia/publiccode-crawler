package crawler

import (
	"errors"
	"strings"
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

// API returns a Domain without tld.
func (domain Domain) API() string {
	truncateIndex := strings.LastIndexAny(domain.Host, ".")
	// It is already an API without tld.
	if truncateIndex == -1 {
		return domain.Host
	}

	return domain.Host[:truncateIndex]
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

func (domain Domain) processAndGetNextURL(url string, wg *sync.WaitGroup, repositories chan Repository, pa PA) (string, error) {
	crawler, err := GetClientAPICrawler(domain.API())
	if err != nil {
		return "", err
	}
	return crawler(domain, url, repositories, pa, wg)
}

func (domain Domain) processSingleRepo(url string, repositories chan Repository) error {
	crawler, err := GetSingleClientAPICrawler(domain.API())
	if err != nil {
		return err
	}
	return crawler(domain, url, repositories)
}

func (domain Domain) generateAPIURL(u string) (string, error) {
	crawler, err := GetAPIURL(domain.API())
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
		log.Infof("%s - API inferred:%s", link, "github")
		return Domain{Host: "github"}, nil
	} else if IsBitbucket(link) {
		log.Infof("%s - API inferred:%s", link, "bitbucket")
		return Domain{Host: "bitbucket"}, nil
	} else if IsGitlab(link) {
		log.Infof("%s - API inferred:%s", link, "gitlab")
		return Domain{Host: "gitlab"}, nil
	}

	return Domain{}, errors.New("this host is not registered: " + host)
}
