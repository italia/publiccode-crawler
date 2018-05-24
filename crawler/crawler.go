package crawler

import (
	"bytes"
	"context"
	"crypto/rand"
	"html/template"
	"math/big"
	"sync"

	"net/http"

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
}

// ProcessRepositories process the repositories channel and check the availability of the file.
func ProcessRepositories(repositories chan Repository, index string, wg *sync.WaitGroup, elasticClient *elastic.Client) {
	log.Debug("Repositories are going to be processed...")

	for repository := range repositories {
		wg.Add(1)
		go checkAvailability(repository, index, wg, elasticClient)
	}
}

// checkAvailability looks for the fileRawUrl and, if found, save it.
func checkAvailability(repository Repository, index string, wg *sync.WaitGroup, elasticClient *elastic.Client) {
	name := repository.Name
	fileRawUrl := repository.FileRawURL
	domain := repository.Domain
	headers := repository.Headers

	// Increment counter for the number of repositories processed.
	metrics.GetCounter("repository_processed", index).Inc()

	resp, err := httpclient.GetURL(fileRawUrl, headers)
	// If it's available and no error returned.
	if resp.Status.Code == http.StatusOK && err == nil {

		// Save to file.
		SaveToFile(domain, name, resp.Body, index)

		// Save to ES.
		SaveToES(domain, name, resp.Body, index, elasticClient)

		// Validate file.
		// TODO: uncomment these lines when mapping and File structure are ready for publiccode.
		// TODO: now validation is useless because we test on .gitignore file.
		// err := validateRemoteFile(resp.Body, fileRawUrl, index)
		// if err != nil {
		// 	log.Warn("Validator fails for: " + fileRawUrl)
		// 	log.Warn("Validator errors:" + err.Error())
		// }
	}

	// Defer waiting group close.
	wg.Done()
}

// ProcessPA delegates the work to single PA crawlers.
func ProcessPA(pa PA, domains []Domain, repositories chan Repository, index string, wg *sync.WaitGroup) {
	log.Debugf("ProcessPA: %s", pa.CodiceIPA)
	// range over repositories.
	for _, repository := range pa.Repositories {
		for _, domain := range domains {
			// if repository API is the domain
			if domain.ClientApi == repository.API {
				for _, org := range repository.Organizations {
					log.Debugf("ProcessPADomain: %s - API: %s", org, repository.API)
					wg.Add(1)
					ProcessPADomain(org, domain, repositories, index, wg)
				}
			}
		}
	}
	wg.Done()

}

// ProcessPADomain starts from the org page and process all the next.
func ProcessPADomain(org string, domain Domain, repositories chan Repository, index string, wg *sync.WaitGroup) {
	// Generate url using go templates. It will replace {{.OrgName}} with "org".
	url := domain.ApiOrgURL
	data := struct{ OrgName string }{OrgName: org}
	// Create a new template and parse the Url into it.
	t := template.Must(template.New("url").Parse(url))
	buf := new(bytes.Buffer)
	// Execute the template: add "data" data in "url".
	t.Execute(buf, data)
	url = buf.String()

	// Process the pages until the end is reached.
	for {
		log.Debugf("processAndGetNextURL handler:%s", url)
		nextURL, err := domain.processAndGetNextURL(url, wg, repositories)
		if err != nil {
			log.Errorf("error reading %s repository list: %v. NextUrl: %v", url, err, nextURL)
			log.Errorf("Retry: %s", nextURL)
			nextURL = url
		}

		// If end is reached, nextUrl is empty.
		if nextURL == "" {
			log.Infof("Url: %s - is the last one for %s.", url, org)
			wg.Done()
			return
		}
		// Update url to nextURL.
		url = nextURL
	}
}

// WaitingLoop waits until all the goroutines counter is zero and close the repositories channel.
// It also switch the alias for elasticsearch index.
func WaitingLoop(repositories chan Repository, index string, wg *sync.WaitGroup, elasticClient *elastic.Client) {
	wg.Wait()

	// Remove old aliases.
	res, err := elasticClient.Aliases().Index("_all").Do(context.Background())
	if err != nil {
		panic(err)
	}
	aliasService := elasticClient.Alias()
	indices := res.IndicesByAlias("publiccode")
	for _, name := range indices {
		log.Debugf("Remove alias from %s to %s", "publiccode", name)
		aliasService.Remove(name, "publiccode").Do(context.Background())
	}

	// Add an alias to the new index.
	log.Debugf("Add alias from %s to %s", index, "publiccode")
	aliasService.Add(index, "publiccode").Do(context.Background())

	close(repositories)
}

// ProcessSingleRepository process a single repository given his url and domain.
func ProcessSingleRepository(url string, domain Domain, repositories chan Repository) error {

	err := domain.processSingleRepo(url, repositories)
	if err != nil {
		return err
	}

	return nil
}

// generateRandomInt returns an integer between 0 and max parameter.
// "Max" must be less than math.MaxInt32
func generateRandomInt(max int) (int, error) {
	result, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(result.Int64()), err
}
