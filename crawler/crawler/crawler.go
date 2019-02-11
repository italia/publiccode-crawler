package crawler

import (
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/italia/developers-italia-backend/crawler/httpclient"
	"github.com/italia/developers-italia-backend/crawler/metrics"
	"github.com/olivere/elastic"
	publiccode "github.com/italia/publiccode-parser-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Sync mutex guard.
var mu sync.Mutex

// Repository is a single code repository. FileRawURL contains the direct url to the raw file.
type Repository struct {
	Name        string
	Hostname    string
	FileRawURL  string
	GitCloneURL string
	GitBranch   string
	Domain      Domain
	Pa          PA
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
		ProcessPADomain(org, domain, pa, repositories, wg)
	}

	wg.Done()
	log.Infof("End ProcessPA on '%s'", pa.ID)
}

// ProcessPADomain starts from the org page and process all the next.
func ProcessPADomain(orgURL string, domain Domain, pa PA, repositories chan Repository, wg *sync.WaitGroup) {
	// generateAPIURL
	orgURL, err := domain.generateAPIURL(orgURL)
	if err != nil {
		log.Errorf("generateAPIURL error: %v", err)
	}
	// Process the pages until the end is reached.
	for {
		log.Debugf("processAndGetNextURL handler: %s", orgURL)
		nextURL, err := domain.processAndGetNextURL(orgURL, wg, repositories, pa)
		if err != nil {
			log.Errorf("error reading %s repository list: %v. NextUrl: %v", orgURL, err, nextURL)
			nextURL = ""
		}

		// If end is reached or fails, nextUrl is empty.
		if nextURL == "" {
			log.Infof("Url: %s - is the last one.", orgURL)
			return
		}
		// Update url to nextURL.
		orgURL = nextURL
	}
}

// WaitingLoop waits until all the goroutines counter is zero and close the repositories channel.
func WaitingLoop(repositories chan Repository, wg *sync.WaitGroup) {
	wg.Wait()

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
	for repository := range repositories {
		wg.Add(1)
		go CheckAvailability(repository, index, wg, elasticClient)
	}
	wg.Wait()
}

// CheckAvailability looks for the FileRawURL and, if found, save it.
func CheckAvailability(repository Repository, index string, wg *sync.WaitGroup, elasticClient *elastic.Client) {
	name := repository.Name
	hostname := repository.Hostname
	fileRawURL := repository.FileRawURL
	gitURL := repository.GitCloneURL
	gitBranch := repository.GitBranch
	domain := repository.Domain
	headers := repository.Headers
	metadata := repository.Metadata
	pa := repository.Pa

	// Hash based on unique git repo URL.
	hash := sha1.New()
	_, err := hash.Write([]byte(gitURL))
	if err != nil {
		log.Errorf("Error generating the repository hash: %+v", err)
		wg.Done()
		return
	}
	hashedRepoURL := fmt.Sprintf("%x", hash.Sum(nil))

	// Increment counter for the number of repositories processed.
	metrics.GetCounter("repository_processed", index).Inc()

	resp, err := httpclient.GetURL(fileRawURL, headers)
	log.Debugf("repository checkAvailability: %s", name)

	// If it's available and no error returned.
	if resp.Status.Code == http.StatusOK && err == nil {
		mu.Lock()
		// Validate file. If invalid, terminate the check.
		err = validateRemoteFile(resp.Body, fileRawURL, pa)
		mu.Unlock()
		if err != nil {
			log.Errorf("%s is an invalid publiccode.", fileRawURL)
			log.Errorf("Errors: %+v", err)
			logBadYamlToFile(fileRawURL)
			wg.Done()
			return
		}

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

		// Calculate Repository activity index and vitality.
		days := 60 // to add in configs.
		activityIndex, vitality, err := CalculateRepoActivity(domain, hostname, name, days)
		if err != nil {
			log.Errorf("error calculating repository Activity to file: %v", err)
		}
		log.Infof("Activity Index for %s: %f", name, activityIndex)
		var vitalitySlice []int
		for i := 0; i < len(vitality); i++ {
			vitalitySlice = append(vitalitySlice, int(vitality[i]))
		}

		// Save to ES.
		err = SaveToES(fileRawURL, hashedRepoURL, name, activityIndex, vitalitySlice, resp.Body, index, elasticClient)
		if err != nil {
			log.Errorf("error saving to ElastcSearch: %v", err)
		}
	}

	// Defer waiting group close.
	wg.Done()
}

func validateRemoteFile(data []byte, fileRawURL string, pa PA) error {
	parser := publiccode.NewParser() 
	parser.RemoteBaseURL = strings.TrimRight(fileRawURL, viper.GetString("CRAWLED_FILENAME"))

	err := parser.Parse(data)
	if err != nil {
		log.Errorf("Error parsing publiccode.yml for %s.", fileRawURL)
		return err
	}

	if pa.CodiceIPA != "" && parser.PublicCode.It.Riuso.CodiceIPA != "" && pa.CodiceIPA != parser.PublicCode.It.Riuso.CodiceIPA {
		return errors.New("codiceIPA for: " + fileRawURL + " is " + parser.PublicCode.It.Riuso.CodiceIPA + ", which differs from the one assigned to the org in the whitelist: " + pa.CodiceIPA)
	}

	return err
}
