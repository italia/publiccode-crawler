package crawler

import (
	"crypto/rand"
	"math/big"
	"sync"
	"time"

	"net/http"
	"net/url"

	"github.com/italia/developers-italia-backend/httpclient"
	"github.com/italia/developers-italia-backend/metrics"
	"github.com/olivere/elastic"

	log "github.com/sirupsen/logrus"
)

// Repository is a single code repository. FileRawURL contains the direct url to the raw file.
type Repository struct {
	Name        string
	Hostname    string
	FileRawURL  string
	GitCloneURL string
	GitBranch   string
	Domain      Domain
	Headers     map[string]string
	Metadata    []byte
}

// ProcessPA delegates the work to single PA crawlers.
func ProcessPA(pa PA, domains []Domain, repositories chan Repository, wg *sync.WaitGroup) {
	log.Infof("Start ProcessPA on '%s'", pa.ID)

	// range over organizations..
	for _, org := range pa.Organizations {
		// Parse as url.URL.
		u, err := url.Parse(org)
		if err != nil {
			log.Errorf("invalid host: %v", err)
		}

		// Check if host is in list of "famous" hosts.
		domain, err := KnownHost(org, u.Hostname(), domains)
		if err != nil {
			log.Error(err)
		}

		// Process the PA domain
		ProcessPADomain(org, domain, repositories, wg)
	}

	wg.Done()

	log.Infof("End ProcessPA on '%s'", pa.ID)
}

// ProcessPADomain starts from the org page and process all the next.
func ProcessPADomain(orgURL string, domain Domain, repositories chan Repository, wg *sync.WaitGroup) {
	// generateAPIURL
	orgURL, err := domain.generateAPIURL(orgURL)
	if err != nil {
		log.Errorf("generateAPIURL error: %v", err)
	}
	// Process the pages until the end is reached.
	for {
		wg.Add(1)
		log.Debugf("processAndGetNextURL handler: %s", orgURL)
		nextURL, err := domain.processAndGetNextURL(orgURL, wg, repositories)
		if err != nil {
			log.Errorf("error reading %s repository list: %v. NextUrl: %v", orgURL, err, nextURL)
			log.Errorf("Retry: %s", nextURL)
			nextURL = orgURL
		}

		// If end is reached, nextUrl is empty.
		if nextURL == "" {
			log.Infof("Url: %s - is the last one.", orgURL)
			wg.Done()
			return
		}
		// Update url to nextURL.
		orgURL = nextURL
		wg.Done()
	}
}

// WaitingLoop waits until all the goroutines counter is zero and close the repositories channel.
func WaitingLoop(repositories chan Repository, wg *sync.WaitGroup) {
	// Waiting initial timer.
	time.Sleep(5 * time.Second)

	wg.Wait()

	// Waiting final timer.
	time.Sleep(5 * time.Second)

	// Close repositories channel.
	log.Debugf("closing repositories chan: len=%d", len(repositories))
	close(repositories)
}

// ProcessSingleRepository process a single repository given his url and domain.
func ProcessSingleRepository(url string, domain Domain, repositories chan Repository) error {
	return domain.processSingleRepo(url, repositories)
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
	hostname := repository.Hostname
	fileRawURL := repository.FileRawURL
	gitURL := repository.GitCloneURL
	gitBranch := repository.GitBranch
	domain := repository.Domain
	headers := repository.Headers
	metadata := repository.Metadata

	// Increment counter for the number of repositories processed.
	metrics.GetCounter("repository_processed", index).Inc()

	resp, err := httpclient.GetURL(fileRawURL, headers)
	log.Debugf("repository checkAvailability: %s", name)

	// If it's available and no error returned.
	if resp.Status.Code == http.StatusOK && err == nil {
		// Save Metadata.
		err = SaveToFile(domain, hostname, name, metadata, index+"_metadata")
		if err != nil {
			log.Errorf("error saving to file: %v", err)
		}

		// Save to file.
		err = SaveToFile(domain, hostname, name, resp.Body, index)
		if err != nil {
			log.Errorf("error saving to file: %v", err)
		}

		// Clone repository.
		err = CloneRepository(domain, hostname, name, gitURL, gitBranch, index)
		if err != nil {
			log.Errorf("error cloning repository %s: %v", gitURL, err)
		}

		// Calculate Repository activity index.
		activityIndex, err := CalculateRepoActivity(domain, hostname, name)
		if err != nil {
			log.Errorf("error calculating repository Activity to file: %v", err)
		}
		log.Debugf("Activity Index for %s: %f", name, activityIndex)

		// Save to ES.
		err = SaveToES(domain, name, resp.Body, index, elasticClient)
		if err != nil {
			log.Errorf("error saving to ElastcSearch: %v", err)
		}
		// TODO: save "metadata" on ES. When mapping is ready.

		// Validate file.
		// TODO: uncomment these lines when mapping and File structure are ready for publiccode.
		// TODO: now validation is useless because we test on .gitignore file.
		// err := validateRemoteFile(resp.Body, fileRawURL, index)
		// if err != nil {
		// 	log.Errorf("Validator fails for: " + fileRawURL)
		// 	log.Errorf("Validator errors:" + err.Error())
		//	wg.Done()
		//	return
		// }
	}

	// Defer waiting group close.
	wg.Done()
}
