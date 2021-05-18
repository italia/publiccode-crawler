package crawler

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/alranel/go-vcsurl/v2"
	"github.com/italia/developers-italia-backend/crawler/elastic"
	"github.com/italia/developers-italia-backend/crawler/ipa"
	"github.com/italia/developers-italia-backend/crawler/jekyll"
	"github.com/italia/developers-italia-backend/crawler/metrics"
	httpclient "github.com/italia/httpclient-lib-go"
	publiccode "github.com/italia/publiccode-parser-go/v2"
	es "github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Crawler is a helper class representing a crawler.
type Crawler struct {
	DryRun         bool

	// Sync mutex guard.
	es             *es.Client
	index          string
	domains        []Domain
	repositories   chan Repository
	publishersWg   sync.WaitGroup
	repositoriesWg sync.WaitGroup
}

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

// NewCrawler initializes a new Crawler object, updates the IPA list and connects to Elasticsearch (if dryRun == false).
func NewCrawler(dryRun bool) *Crawler {
	var c Crawler
	var err error

	c.DryRun = dryRun

	// Make sure the data directory exists or spit an error
	if stat, err := os.Stat(viper.GetString("CRAWLER_DATADIR")); err != nil || !stat.IsDir() {
		log.Fatalf("The configured data directory (%v) does not exist: %v", viper.GetString("CRAWLER_DATADIR"), err)
	}

	// Read and parse list of domains.
	c.domains, err = ReadAndParseDomains("domains.yml")
	if err != nil {
		log.Fatal(err)
	}

	// Initiate a channel of repositories.
	c.repositories = make(chan Repository, 1000)

	// Register Prometheus metrics.
	metrics.RegisterPrometheusCounter("repository_processed", "Number of repository processed.", c.index)
	metrics.RegisterPrometheusCounter("repository_file_saved", "Number of file saved.", c.index)
	metrics.RegisterPrometheusCounter("repository_file_indexed", "Number of file indexed.", c.index)
	metrics.RegisterPrometheusCounter("repository_cloned", "Number of repository cloned", c.index)
	//metrics.RegisterPrometheusCounter("repository_file_saved_valid", "Number of valid file saved.", c.index)

	if c.DryRun {
		log.Info("Skipping ElasticSearch update (--dry-run)")
		return &c
	}

	log.Debug("Connecting to ElasticSearch...")
	c.es, err = elastic.ClientFactory(
		viper.GetString("ELASTIC_URL"),
		viper.GetString("ELASTIC_USER"),
		viper.GetString("ELASTIC_PWD"))
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Successfully connected to ElasticSearch")

	// Update ipa to lastest data.
	err = ipa.UpdateFromIndicePAIfNeeded(c.es)
	if err != nil {
		log.Error(err)
	}

	// Initialize ES index mapping
	c.index = viper.GetString("ELASTIC_PUBLICCODE_INDEX")
	err = elastic.CreateIndexMapping(c.index, elastic.PubliccodeMapping, c.es)
	if err != nil {
		log.Fatal(err)
	}

	// Create ES index with mapping "administration-codiceIPA".
	err = elastic.CreateIndexMapping(viper.GetString("ELASTIC_PUBLISHERS_INDEX"), elastic.AdministrationsMapping, c.es)
	if err != nil {
		log.Fatal(err)
	}

	return &c
}

// CrawlRepo crawls a single repository.
func (c *Crawler) CrawlRepo(repoURL string, pa PA) error {
	log.Infof("Processing repository: %s", repoURL)

	// Check if current host is in known in domains.yml hosts.
	domain, err := c.KnownHost(repoURL)
	if err != nil {
		return err
	}

	// Process repository.
	err = domain.processSingleRepo(repoURL, c.repositories, pa)
	if err != nil {
		return err
	}
	close(c.repositories)
	return c.crawl()
}

// CrawlPublishers processes a list of publishers.
func (c *Crawler) CrawlPublishers(publishers []PA) ([]string, error) {
	// Count configured orgs
	orgCount := 0
	for _, pa := range publishers {
		orgCount += len(pa.Organizations)
	}
	log.Infof("%v organizations belonging to %v publishers are going to be scanned",
		orgCount, len(publishers))

	// Process every item in publishers.
	for _, pa := range publishers {
		c.publishersWg.Add(1)
		go c.CrawlPublisher(pa)
	}

	// Close the repositories channel when all the publisher goroutines are done
	go func() {
		c.publishersWg.Wait()
		close(c.repositories)
	}()

	// here we got all repos to be scanned
	// it's time to check blacklist and wheter one of them
	// is listed.
	// we should return the ones listed to crawl command
	// and call deleteFromES if present
	toBeRemoved := c.removeBlackListedFromRepositories(GetAllBlackListedRepos())

	return toBeRemoved, c.crawl()
}

// removeBlackListedFromRepositories this function is in charge
// to discard repositories in blacklists.
// It returns a slice of them, ready to be removed
// from elasticsearch.
func (c *Crawler) removeBlackListedFromRepositories(listedRepos map[string]string) (toBeRemoved []string) {
	temp := make(chan Repository, 1000)
	for repo := range c.repositories {
		if val, ok := listedRepos[repo.GitCloneURL]; ok {
			// add repository that should be processed but
			// they are marked as blacklisted
			// and then ready to be removed from ES if they exist
			toBeRemoved = append(toBeRemoved, val)
			log.Warnf("marked as blacklisted %s", val)
		} else {
			temp <- repo
		}
	}
	close(temp)
	c.repositories = temp
	return
}

func (c *Crawler) crawl() error {
	reposChan := make(chan Repository)

	// Start the metrics server.
	go metrics.StartPrometheusMetricsServer()

	defer c.publishersWg.Wait()

	// Get cpus number
	numCPUs := runtime.NumCPU()

	// Process the repositories in order to retrieve the files.
	for i := 0; i < numCPUs; i++ {
		c.repositoriesWg.Add(1)
		go c.ProcessRepositories(reposChan)
	}

	for repo := range c.repositories {
		reposChan <- repo
	}
	close(reposChan)
	c.repositoriesWg.Wait()

	if c.DryRun {
		log.Info("Skipping ElasticSearch indexes update (--dry-run)")

		return nil
	}

	// ElasticFlush to flush all the operations on ES.
	err := elastic.Flush(c.index, c.es)
	if err != nil {
		log.Errorf("Error flushing ElasticSearch: %v", err)
	}

	// Update Elastic alias.
	err = elastic.AliasUpdate(viper.GetString("ELASTIC_PUBLISHERS_INDEX"), viper.GetString("ELASTIC_ALIAS"), c.es)
	if err != nil {
		return fmt.Errorf("Error updating Elastic Alias: %v", err)
	}
	err = elastic.AliasUpdate(c.index, viper.GetString("ELASTIC_ALIAS"), c.es)
	if err != nil {
		return fmt.Errorf("Error updating Elastic Alias: %v", err)
	}

	return nil
}

// ExportForJekyll exports YAML data files for the Jekyll website.
func (c *Crawler) ExportForJekyll() error {
	if c.DryRun {
		log.Info("Skipping YAML output (--dry-run)")
		return nil;
	}

	return jekyll.GenerateJekyllYML(c.es)
}

// CrawlPublisher delegates the work to single PA crawlers.
func (c *Crawler) CrawlPublisher(pa PA) {
	log.Infof("Processing publisher: %s", pa.Name)
	defer c.publishersWg.Done()

	for _, orgURL := range pa.Organizations {
		// Check if host is in list of known code hosting domains
		domain, err := c.KnownHost(orgURL)
		if err != nil {
			log.Error(err)
		}

		// Process the organization
		c.CrawlOrg(orgURL, domain, pa)
	}

	for _, repoURL := range pa.Repositories {
		// Check if host is in list of known code hosting domains
		domain, err := c.KnownHost(repoURL)
		if err != nil {
			log.Error(err)
		}

		domain.processSingleRepo(repoURL, c.repositories, pa)
	}
}

// CrawlOrg fetches all the repositories belonging to an org and crawls them.
func (c *Crawler) CrawlOrg(orgURL string, domain *Domain, pa PA) {
	orgURLs, err := domain.generateAPIURLs(orgURL)
	if err != nil {
		log.Errorf("generateAPIURLs error: %v", err)
	}

ORG:
	for _, orgURL := range orgURLs {
		// Process the pages until the end is reached.
		for {
			nextURL, err := domain.processAndGetNextURL(orgURL, c.repositories, pa)
			if err != nil {
				log.Errorf("error reading %s repository list: %v; nextURL: %v", orgURL, err, nextURL)
				continue ORG
			}

			// If end is reached or fails, nextURL is empty.
			if nextURL == "" {
				return
			}
			// Update url to nextURL.
			orgURL = nextURL
		}
	}
}

// generateRandomInt returns an integer between 0 and max parameter.
// "Max" must be less than math.MaxInt32
func generateRandomInt(max int) (int, error) {
	result, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(result.Int64()), err
}

// ProcessRepositories process the repositories channel and check the availability of the file.
func (c *Crawler) ProcessRepositories(repos chan Repository) {
	defer c.repositoriesWg.Done()

	for repository := range repos {
		c.ProcessRepo(repository)
	}
}

type logEntry struct {
	Datetime string `json:"datetime"`
	Message  string `json:"message"`
}

func addLogEntry(logEntries *[]logEntry, message string) {
	*logEntries = append(
		*logEntries,
		logEntry{Datetime: time.Now().UTC().Format(time.RFC3339), Message: message},
	)
}

// ProcessRepo looks for a publiccode.yml file in a repository, and if found it processes it.
func (c *Crawler) ProcessRepo(repository Repository) {
	var logEntries []logEntry

	var message string = ""

	// Write the log to a file, so it can be accessed from outside at
	// http://crawler-host/$codehosting/$org/$reponame/log.txt
	defer func() {
		fname := path.Join(
			viper.GetString("OUTPUT_DIR"),
			repository.Hostname,
			path.Clean(repository.Name),
			"log.json",
		)

		if err := os.MkdirAll(filepath.Dir(fname), 0775); err != nil {
			log.Errorf("[%s]: %s", repository.Name, err.Error())

			return
		}

		jsonOut, _ := json.Marshal(logEntries)
		if err := ioutil.WriteFile(fname, jsonOut, 0644); err != nil {
			log.Errorf("[%s]: %s", repository.Name, err.Error())

			return
		}
	}()

	// Increment counter for the number of repositories processed.
	metrics.GetCounter("repository_processed", c.index).Inc()

	resp, err := httpclient.GetURL(repository.FileRawURL, repository.Headers)

	if resp.Status.Code != http.StatusOK || err != nil {
		message = fmt.Sprintf("[%s] Failed to GET publiccode.yml\n", repository.Name)
		log.Errorf(message)

		addLogEntry(&logEntries, message)
		return
	}

	message = fmt.Sprintf("[%s] publiccode.yml found at %s\n", repository.Name, repository.FileRawURL)
	log.Infof(message)
	addLogEntry(&logEntries, message)

	var parser *publiccode.Parser
	// Validate the publiccode.yml
	if repository.Pa.UnknownIPA {
		message = fmt.Sprintf(
			"[%s] When UnknownIPA is set to true IPA match with whitelists will be skipped\n",
			repository.Name,
		)

		log.Warn(message)
		addLogEntry(&logEntries, message)
	} else {
		parser, err = publiccode.NewParser(repository.FileRawURL)
		if err != nil {
			message = fmt.Sprintf("[%s] BAD publiccode.yml: %+v\n", repository.Name, err)
			log.Errorf(message)
			addLogEntry(&logEntries, message)

			return
		}
		err = parser.ParseInDomain(resp.Body, repository.Domain.Host, repository.Domain.UseTokenFor, repository.Domain.BasicAuth)
		if err != nil {
			message = fmt.Sprintf("[%s] BAD publiccode.yml: %+v\n", repository.Name, err)
			log.Errorf(message)
			addLogEntry(&logEntries, message)

			return
		}

		err = validateFile(repository.Pa, *parser, repository.FileRawURL)
		if err != nil {
			message = fmt.Sprintf("[%s] BAD publiccode.yml: %+v\n", repository.Name, err)
			log.Errorf(message)
			addLogEntry(&logEntries, message)

			if ! c.DryRun {
				logBadYamlToFile(repository.FileRawURL)
			}

			return
		}
	}

	message = fmt.Sprintf("[%s] GOOD publiccode.yml\n", repository.Name)
	log.Infof(message)
	addLogEntry(&logEntries, message)

	if c.DryRun {
		log.Infof("[%s]: Skipping repository clone and save to ElasticSearch (--dry-run)", repository.Name)
		return;
	}

	// Clone repository.
	err = CloneRepository(repository.Domain, repository.Hostname, repository.Name, repository.GitCloneURL, repository.GitBranch, c.index)
	if err != nil {
		message = fmt.Sprintf("[%s] error while cloning: %v\n", repository.Name, err)
		log.Errorf(message)

		addLogEntry(&logEntries, message)
	}

	// Calculate Repository activity index and vitality. Defaults to 60 days.
	var activityDays int = 60
	if viper.IsSet("ACTIVITY_DAYS") {
		activityDays = viper.GetInt("ACTIVITY_DAYS")
	}
	activityIndex, vitality, err := repository.CalculateRepoActivity(activityDays)
	if err != nil {
		message = fmt.Sprintf("[%s] error calculating activity index: %v\n", repository.Name, err)

		log.Errorf(message)
		addLogEntry(&logEntries, message)
	}
	message = fmt.Sprintf("[%s] activity index in the last %d days: %f\n", repository.Name, activityDays, activityIndex)
	log.Infof(message)
	addLogEntry(&logEntries, message)

	var vitalitySlice []int
	for i := 0; i < len(vitality); i++ {
		vitalitySlice = append(vitalitySlice, int(vitality[i]))
	}

	// Save to ES.
	err = c.saveToES(repository, activityIndex, vitalitySlice, *parser)
	if err != nil {
		message = fmt.Sprintf("[%s] error saving to ElasticSearch: %v\n", repository.Name, err)
		log.Errorf(message)

		addLogEntry(&logEntries, message)
	}
}

// validateFile will check if codiceIPA match
// with relative entry in whitelist.
// Using `one` command this check will be skipped.
func validateFile(pa PA, parser publiccode.Parser, fileRawURL string) error {
	u, _ := url.Parse(fileRawURL)
	repo1 := vcsurl.GetRepo((*url.URL)(u))

	repo2 := vcsurl.GetRepo((*url.URL)(parser.PublicCode.URL))

	if repo1 != nil && repo2 != nil {
		// Let's ignore the schema when checking for equality.
		//
		// This is mainly to match repos regardless of whether they are served
		// through HTTPS or HTTP.
		repo1.Scheme, repo2.Scheme = "", ""

		if !strings.EqualFold(repo1.String(), repo2.String()) {
			return errors.New(
				fmt.Sprintf(
					"declared url (%s) and actual publiccode.yml location URL (%s) "+
					"are not in the same repo: '%s' vs '%s'",
					parser.PublicCode.URL, fileRawURL, repo2, repo1,
				),
			)
		}
	}

	if !strings.EqualFold(
		strings.TrimSpace(pa.CodiceIPA),
		strings.TrimSpace(parser.PublicCode.It.Riuso.CodiceIPA),
	) {
		return errors.New("codiceIPA for: " + fileRawURL + " is " + parser.PublicCode.It.Riuso.CodiceIPA + ", which differs from the one assigned to the org in the whitelist: " + pa.CodiceIPA)
	}

	return nil
}
