package crawler

import (
	"crypto/rand"
	"math/big"
	"sync"

	"net/http"
	"net/url"

	"github.com/italia/developers-italia-backend/httpclient"
	"github.com/italia/developers-italia-backend/metrics"
	"github.com/olivere/elastic"

	log "github.com/sirupsen/logrus"
)

// Repository is a single code repository. FileRawURL contains the direct url to the raw file.
type Repository struct {
	Name       string
	FileRawURL string
	Domain     Domain
	Headers    map[string]string
	Metadata   []byte
}

// ProcessPA delegates the work to single PA crawlers.
func ProcessPA(pa PA, domains []Domain, repositories chan Repository, wg *sync.WaitGroup) {
	log.Debugf("ProcessPA: %s", pa.CodiceIPA)

	// range over organizations..
	for _, org := range pa.Organizations {
		knownHost := false
		domain := Domain{}
		// Parse as url.URL.
		u, err := url.Parse(org)
		if err != nil {
			log.Errorf("invalid host: %v", err)
		}

		// Check if host is in list of "famous" hosts.
		for _, d := range domains {
			if u.Hostname() == d.Host {
				// Process this host
				knownHost = true
				domain = d
			}

		}

		if knownHost {
			log.Infof("%s - API known:%s", org, u.Hostname())
			// Host is detected.
			ProcessPADomain(org, domain, repositories, wg)
		} else {
			// host unknown, needs to be inferred.
			if isGithub(org) {
				log.Infof("%s - API inferred:%s", org, "github")
				ProcessPADomain(org, Domain{Host: u.Hostname()}, repositories, wg)
			} else if isBitbucket(org) {
				log.Infof("%s - API inferred:%s", org, "bitbucket")
				ProcessPADomain(org, Domain{Host: u.Hostname()}, repositories, wg)
			} else if isGitlab(org) {
				log.Infof("%s - API inferred:%s", org, "gitlab")
				ProcessPADomain(org, Domain{Host: u.Hostname()}, repositories, wg)
			}
		}
	}

	wg.Done()

}

// ProcessPADomain starts from the org page and process all the next.
func ProcessPADomain(orgURL string, domain Domain, repositories chan Repository, wg *sync.WaitGroup) {
	orgURL, err := domain.generateAPIURL(orgURL)
	if err != nil {
		log.Errorf("generateAPIURL error: %v", err)
	}
	// Process the pages until the end is reached.
	for {
		log.Debugf("processAndGetNextURL handler: %s", orgURL)
		nextURL, err := domain.processAndGetNextURL(orgURL, wg, repositories)
		if err != nil {
			log.Errorf("error reading %s repository list: %v. NextUrl: %v", orgURL, err, nextURL)
			log.Errorf("Retry: %s", nextURL)
			nextURL = orgURL
		}

		// If end is reached, nextUrl is empty.
		if nextURL == "" {
			log.Infof("Url: %s - is the last one for %s.", orgURL, domain.Host)
			return
		}
		// Update url to nextURL.
		orgURL = nextURL
	}
}

// WaitingLoop waits until all the goroutines counter is zero and close the repositories channel.
// It also switch the alias for elasticsearch index.
func WaitingLoop(repositories chan Repository, wg *sync.WaitGroup) {
	wg.Wait()

	// Close repositories channel.
	log.Debugf("closing repositories chan: len=%d", len(repositories))
	close(repositories)
}

// ProcessSingleRepository process a single repository given his url and domain.
func ProcessSingleRepository(url string, domain Domain, repositories chan Repository) error {
	err := domain.processSingleRepo(url, repositories)

	return err

}

// generateRandomInt returns an integer between 0 and max parameter.
// "Max" must be less than math.MaxInt32
func generateRandomInt(max int) (int, error) {
	result, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(result.Int64()), err
}

// ProcessRepositories process the repositories channel and check the availability of the file.
func ProcessRepositories(repositories chan Repository, index string, wg *sync.WaitGroup, elasticClient *elastic.Client) {
	log.Debug("Repositories are going to be processed...")

	for repository := range repositories {
		wg.Add(1)
		go checkAvailability(repository, index, wg, elasticClient)
	}
}

// checkAvailability looks for the FileRawURL and, if found, save it.
func checkAvailability(repository Repository, index string, wg *sync.WaitGroup, elasticClient *elastic.Client) {
	name := repository.Name
	FileRawURL := repository.FileRawURL
	domain := repository.Domain
	headers := repository.Headers
	metadata := repository.Metadata

	// Increment counter for the number of repositories processed.
	metrics.GetCounter("repository_processed", index).Inc()

	resp, err := httpclient.GetURL(FileRawURL, headers)
	// If it's available and no error returned.
	if resp.Status.Code == http.StatusOK && err == nil {

		// Save Metadata.
		SaveToFile(domain, name, metadata, index+"_metadata")

		// Save to file.
		SaveToFile(domain, name, resp.Body, index)

		// Save to ES.
		SaveToES(domain, name, resp.Body, index, elasticClient)

		// TODO: save "metadata" on ES. When mapping is ready.

		// Validate file.
		// TODO: uncomment these lines when mapping and File structure are ready for publiccode.
		// TODO: now validation is useless because we test on .gitignore file.
		// err := validateRemoteFile(resp.Body, FileRawURL, index)
		// if err != nil {
		// 	log.Warn("Validator fails for: " + FileRawURL)
		// 	log.Warn("Validator errors:" + err.Error())
		// }
	}

	// Defer waiting group close.
	wg.Done()
}
