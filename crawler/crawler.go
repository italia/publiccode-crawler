package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	"github.com/italia/developers-italia-backend/apiclient"
	"github.com/italia/developers-italia-backend/common"
	"github.com/italia/developers-italia-backend/git"
	"github.com/italia/developers-italia-backend/jekyll"
	"github.com/italia/developers-italia-backend/metrics"
	"github.com/italia/developers-italia-backend/scanner"
	httpclient "github.com/italia/httpclient-lib-go"
	publiccode "github.com/italia/publiccode-parser-go/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Crawler is a helper class representing a crawler.
type Crawler struct {
	DryRun bool

	Index          string
	domains        []common.Domain
	repositories   chan common.Repository
	// Sync mutex guard.
	publishersWg   sync.WaitGroup
	repositoriesWg sync.WaitGroup

	gitHubScanner    scanner.Scanner
	gitLabScanner    scanner.Scanner
	bitBucketScanner scanner.Scanner

	apiClient        apiclient.ApiClient
}

// NewCrawler initializes a new Crawler object and connects to Elasticsearch (if dryRun == false).
func NewCrawler(dryRun bool) *Crawler {
	var c Crawler
	var err error

	c.DryRun = dryRun

	// Make sure the data directory exists or spit an error
	if stat, err := os.Stat(viper.GetString("CRAWLER_DATADIR")); err != nil || !stat.IsDir() {
		log.Fatalf("The configured data directory (%v) does not exist: %v", viper.GetString("CRAWLER_DATADIR"), err)
	}

	// Read and parse list of domains.
	c.domains, err = common.ReadAndParseDomains("domains.yml")
	if err != nil {
		log.Fatal(err)
	}

	// Initiate a channel of repositories.
	c.repositories = make(chan common.Repository, 1000)

	// Register Prometheus metrics.
	metrics.RegisterPrometheusCounter("repository_processed", "Number of repository processed.", c.Index)
	metrics.RegisterPrometheusCounter("repository_file_saved", "Number of file saved.", c.Index)
	metrics.RegisterPrometheusCounter("repository_file_indexed", "Number of file indexed.", c.Index)
	metrics.RegisterPrometheusCounter("repository_cloned", "Number of repository cloned", c.Index)
	//metrics.RegisterPrometheusCounter("repository_file_saved_valid", "Number of valid file saved.", c.Index)

	c.gitHubScanner = scanner.NewGitHubScanner()
	c.gitLabScanner = scanner.NewGitLabScanner()
	c.bitBucketScanner = scanner.NewBitBucketScanner()

	c.apiClient = apiclient.NewClient()

	return &c
}

// CrawlRepo crawls a single repository (only used by the 'one' command).
func (c *Crawler) CrawlRepo(repoURL url.URL, publisher common.Publisher) error {
	log.Infof("Processing repository: %s", repoURL.String())

	var err error
	if vcsurl.IsGitHub(&repoURL) {
		err = c.gitHubScanner.ScanRepo(repoURL, publisher, c.repositories)
	} else if vcsurl.IsBitBucket(&repoURL) {
		err = c.bitBucketScanner.ScanRepo(repoURL, publisher, c.repositories)
	} else if vcsurl.IsGitLab(&repoURL) {
		err = c.gitLabScanner.ScanRepo(repoURL, publisher, c.repositories)
	} else {
		err = fmt.Errorf("unsupported code hosting platform for %s", repoURL.String())
	}

	if err != nil {
		return err
	}

	close(c.repositories)
	return c.crawl()
}

// CrawlPublishers processes a list of publishers.
func (c *Crawler) CrawlPublishers(publishers []common.Publisher) error {
	// Count configured orgs
	orgCount := 0
	for _, publisher := range publishers {
		orgCount += len(publisher.Organizations)
	}
	log.Infof("%v organizations belonging to %v publishers are going to be scanned",
		orgCount, len(publishers))

	// Process every item in publishers.
	for _, publisher := range publishers {
		c.publishersWg.Add(1)
		go c.ScanPublisher(publisher)
	}

	// Close the repositories channel when all the publisher goroutines are done
	go func() {
		c.publishersWg.Wait()
		close(c.repositories)
	}()

	return c.crawl()
}

// removeBlackListedFromRepositories this function is in charge
// to discard repositories in blacklists.
// It returns a slice of them, ready to be removed
// from elasticsearch.
func (c *Crawler) removeBlackListedFromRepositories(listedRepos map[string]string) (toBeRemoved []string) {
	temp := make(chan common.Repository, 1000)
	for repo := range c.repositories {
		if val, ok := listedRepos[repo.URL.String()]; ok {
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
	reposChan := make(chan common.Repository)

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

	return nil
}

// ExportForJekyll exports YAML data files for the Jekyll website.
func (c *Crawler) ExportForJekyll() error {
	if c.DryRun {
		log.Info("Skipping YAML output (--dry-run)")
		return nil
	}

	return jekyll.GenerateJekyllYML(c.Es)
}

// ScanPublisher scans all the publisher' repositories and sends the ones
// with a valid publiccode.yml to the repositories channel.
func (c *Crawler) ScanPublisher(publisher common.Publisher) {
	log.Infof("Processing publisher: %s", publisher.Name)
	defer c.publishersWg.Done()

	var err error
	for _, u := range publisher.Organizations {
		orgURL := (url.URL)(u)

		if vcsurl.IsGitHub(&orgURL) {
			err = c.gitHubScanner.ScanGroupOfRepos(orgURL, publisher, c.repositories)
		} else if vcsurl.IsBitBucket(&orgURL) {
			err = c.bitBucketScanner.ScanGroupOfRepos(orgURL, publisher, c.repositories)
		} else if vcsurl.IsGitLab(&orgURL) {
			err = c.gitLabScanner.ScanGroupOfRepos(orgURL, publisher, c.repositories)
		} else {
			err = fmt.Errorf("unsupported code hosting platform for %s", u.String())
		}
		if err != nil {
			if errors.Is(err, scanner.ErrPubliccodeNotFound) {
				log.Warnf("[%s] %s", orgURL.String(), err.Error())
			} else {
				log.Error(err)
			}
		}
	}

	for _, u := range publisher.Repositories {
		repoURL := (url.URL)(u)

		if vcsurl.IsGitHub(&repoURL) {
			err = c.gitHubScanner.ScanRepo(repoURL, publisher, c.repositories)
		} else if vcsurl.IsBitBucket(&repoURL) {
			err = c.bitBucketScanner.ScanRepo(repoURL, publisher, c.repositories)
		} else if vcsurl.IsGitLab(&repoURL) {
			err = c.gitLabScanner.ScanRepo(repoURL, publisher, c.repositories)
		} else {
			err = fmt.Errorf("unsupported code hosting platform for %s", u.String())
		}

		if err != nil {
			if errors.Is(err, scanner.ErrPubliccodeNotFound) {
				log.Warnf("[%s] %s", repoURL.String(), err.Error())
			} else {
				log.Error(err)
			}
		}
	}
}

// ProcessRepositories process the repositories channel, check the repo's publiccode.yml
// and send new data to the API if the publiccode.yml file is valid.
func (c *Crawler) ProcessRepositories(repos chan common.Repository) {
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
func (c *Crawler) ProcessRepo(repository common.Repository) {
	var logEntries []logEntry

	var message string

	// Write the log to a file, so it can be accessed from outside at
	// http://crawler-host/$codehosting/$org/$reponame/log.txt
	defer func() {
		fname := path.Join(
			viper.GetString("OUTPUT_DIR"),
			repository.URL.String(),
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
	metrics.GetCounter("repository_processed", c.Index).Inc()

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
	parser, err = publiccode.NewParser(repository.FileRawURL)
	if err != nil {
		message = fmt.Sprintf("[%s] BAD publiccode.yml: %+v\n", repository.Name, err)
		log.Errorf(message)
		addLogEntry(&logEntries, message)

		return
	}

	domain := common.GetDomain(c.domains, repository.URL.Host)
	err = parser.ParseInDomain(resp.Body, domain.Host, domain.UseTokenFor, domain.BasicAuth)
    if err != nil {
		valid := true
	out:
		for _, res := range err.(publiccode.ValidationResults) {
			switch res.(type) {
			case publiccode.ValidationError:
				valid = false
				break out
			}
		}

		if !valid {
			message = fmt.Sprintf("[%s] BAD publiccode.yml: %+v\n", repository.Name, err)
			log.Errorf(message)
			addLogEntry(&logEntries, message)

			return
		}
	}

	// HACK: Publishers named "_"" are special and get to skip the additional checks.
	// This can be used to add repositories and organizations, under the crawler's admins control,
	// that describe arbitrary repos (eg. metarepos for other entities)
	if repository.Publisher.Name != "_" {
		err = validateFile(repository.Publisher, *parser, repository.FileRawURL)
		if err != nil {
			message = fmt.Sprintf("[%s] BAD publiccode.yml: %+v\n", repository.Name, err)
			log.Errorf(message)
			addLogEntry(&logEntries, message)

			if !c.DryRun {
				common.LogBadYamlToFile(repository.FileRawURL)
			}

			return
		}
	}

	message = fmt.Sprintf("[%s] GOOD publiccode.yml\n", repository.Name)
	log.Infof(message)
	addLogEntry(&logEntries, message)

	if c.DryRun {
		log.Infof("[%s]: Skipping repository clone and save to ElasticSearch (--dry-run)", repository.Name)
		return
	}

	// Clone repository.
	err = git.CloneRepository(repository.URL.Host, repository.Name, parser.PublicCode.URL.String(), c.Index)
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
	activityIndex, vitality, err := git.CalculateRepoActivity(repository, activityDays)
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

	// XXX doc first is current url
	urls := []url.URL{repository.CanonicalURL,}
	if repository.URL != repository.CanonicalURL {
		urls = append(urls, repository.URL)
	}

	publiccodeYml, err := parser.ToYAML()
	if err != nil {
		log.Errorf("XXX: %w", err)
	}

	_, err = c.apiClient.PutSoftware(urls, string(publiccodeYml))
	if err != nil {
		log.Errorf("XXX: %w", err)
	}
}

// validateFile checks if it.riuso.codiceIPA in the publiccode.yml matches with the
// Publisher's Id
// Using `one` command this check will be skipped.
func validateFile(publisher common.Publisher, parser publiccode.Parser, fileRawURL string) error {
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
			return fmt.Errorf(
				"declared url (%s) and actual publiccode.yml location URL (%s) "+
					"are not in the same repo: '%s' vs '%s'",
				parser.PublicCode.URL, fileRawURL, repo2, repo1,
			)
		}
	}

	if !strings.EqualFold(
		strings.TrimSpace(publisher.Id),
		strings.TrimSpace(parser.PublicCode.It.Riuso.CodiceIPA),
	) {
		return errors.New("id for: " + fileRawURL + " is " + parser.PublicCode.It.Riuso.CodiceIPA + ", which differs from the one assigned to the org in the publishers file: " + publisher.Id)
	}

	return nil
}
