package crawler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Domain is a single code hosting service.
type Domain struct {
	// Domains.yml data
	Host        string   `yaml:"host"`
	UseTokenFor []string `yaml:"use-token-for"`
	BasicAuth   []string `yaml:"basic-auth"`
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

func (domain Domain) processAndGetNextURL(url url.URL, repositories chan Repository, publisher Publisher) (*url.URL, error) {
	crawler, err := GetClientAPICrawler(domain.API())
	if err != nil {
		return nil, err
	}
	return crawler(domain, url, repositories, publisher)
}

func (domain Domain) processSingleRepo(url url.URL, repositories chan Repository, publisher Publisher) error {
	crawler, err := GetSingleClientAPICrawler(domain.API())
	if err != nil {
		return err
	}
	return crawler(domain, url, repositories, publisher)
}

func (domain Domain) generateAPIURLs(u url.URL) ([]url.URL, error) {
	crawler, err := GetAPIURL(domain.API())
	if err != nil {
		return []url.URL{u}, err
	}
	return crawler(u)
}

// KnownHost detect the the right Domain API from the given URL and returns it.
// If no API is recognized will return an empty domain and an error.
func (c *Crawler) KnownHost(link url.URL) (*Domain, error) {
	for _, domain := range c.domains {
		if link.Hostname() == domain.Host {
			// Host is found in the host list.
			return &domain, nil
		}
	}

	// host unknown, needs to be inferred.
	if IsGithub(link.String()) {
		log.Infof("%s - API inferred: %s", link.String(), "github")
		return &Domain{Host: "github"}, nil
	} else if IsBitbucket(link.String()) {
		log.Infof("%s - API inferred: %s", link.String(), "bitbucket")
		return &Domain{Host: "bitbucket"}, nil
	} else if IsGitlab(link.String()) {
		log.Infof("%s - API inferred: %s", link.String(), "gitlab")
		return &Domain{Host: "gitlab"}, nil
	}

	return &Domain{}, errors.New("unable to detect code hosting platform: " + link.Hostname())
}
