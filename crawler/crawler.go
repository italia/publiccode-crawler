package crawler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/alranel/go-vcsurl/v2"
	"github.com/italia/developers-italia-backend/apiclient"
	"github.com/italia/developers-italia-backend/common"
	"github.com/italia/developers-italia-backend/git"
	"github.com/italia/developers-italia-backend/metrics"
	"github.com/italia/developers-italia-backend/scanner"
	httpclient "github.com/italia/httpclient-lib-go"
	publiccode "github.com/italia/publiccode-parser-go/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
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

	datadir := viper.GetString("CRAWLER_DATADIR")
	if err := os.MkdirAll(datadir, 0744); err != nil {
		log.Fatalf("can't create data directory (%s): %s", datadir, err.Error())
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
		err = fmt.Errorf(
			"publisher %s: unsupported code hosting platform for %s",
			publisher.Name,
			repoURL.String(),
		)
	}

	if err != nil {
		return err
	}

	close(c.repositories)
	return c.crawl()
}

// CrawlPublishers processes a list of publishers.
func (c *Crawler) CrawlPublishers(publishers []common.Publisher) error {
	groupsNum := 0
	for _, publisher := range publishers {
		groupsNum += len(publisher.Organizations)
	}

	reposNum := 0
	for _, publisher := range publishers {
		reposNum += len(publisher.Repositories)
	}

	log.Infof("Scanning %d publishers (%d orgs + %d repositories)", len(publishers), groupsNum, reposNum)

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
				err = fmt.Errorf(
				"publisher %s: unsupported code hosting platform for %s",
				publisher.Name,
				u.String(),
			)
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
			err = fmt.Errorf(
				"publisher %s: unsupported code hosting platform for %s",
				publisher.Name,
				u.String(),
			)
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


// ProcessRepo looks for a publiccode.yml file in a repository, and if found it processes it.
func (c *Crawler) ProcessRepo(repository common.Repository) {
	var logEntries []string

	var software *apiclient.Software

	defer func() {
		for _, e := range logEntries {
			log.Info(e)
		}

		if !c.DryRun {
			entries := strings.Join(logEntries, "\n")

			var err error
			if software != nil {
				_, err = c.apiClient.PostSoftwareLog(software.ID, entries)
			} else {
				_, err = c.apiClient.PostLog(entries)
			}

			if err != nil {
				log.Errorf("[%s]: %s", repository.Name, err.Error())
			}
		}
	}()

	// Increment counter for the number of repositories processed.
	metrics.GetCounter("repository_processed", c.Index).Inc()

	software, err := c.apiClient.GetSoftwareByURL(repository.URL.String())
	if err != nil {
		logEntries = append(logEntries, "[%s] failed to GET software from API: %s\n", repository.Name, err.Error())

		return
	}

	if software != nil && !software.Active {
		logEntries = append(logEntries, "[%s] software has active = false, skipping update")

		return
	}

	resp, err := httpclient.GetURL(repository.FileRawURL, repository.Headers)
	if resp.Status.Code != http.StatusOK || err != nil {
		logEntries = append(logEntries, fmt.Sprintf("[%s] Failed to GET publiccode.yml", repository.Name))

		return
	}

	logEntries = append(
		logEntries,
		fmt.Sprintf(
			"[%s] publiccode.yml found at %s\n",
				repository.CanonicalURL.String(),
				repository.FileRawURL,
		),
	)

	var parser *publiccode.Parser
	parser, err = publiccode.NewParser(repository.FileRawURL)
	if err != nil {
		logEntries = append(logEntries,fmt.Sprintf("[%s] BAD publiccode.yml: %+v\n", repository.Name, err))

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
			logEntries = append(logEntries, fmt.Sprintf("[%s] BAD publiccode.yml: %+v\n", repository.Name, err))

			return
		}
	}

	// HACK: Publishers named "_"" are special and get to skip the additional checks.
	// This can be used to add repositories and organizations, under the crawler's admins control,
	// that describe arbitrary repos (eg. metarepos for other entities)
	if repository.Publisher.Name != "_" {
		err = validateFile(repository.Publisher, *parser, repository.FileRawURL)
		if err != nil {
			logEntries = append(logEntries, fmt.Sprintf("[%s] BAD publiccode.yml: %+v\n", repository.Name, err))

			return
		}
	}

	logEntries = append(logEntries, fmt.Sprintf("[%s] GOOD publiccode.yml\n", repository.Name))

	if c.DryRun {
		log.Infof("[%s]: Skipping other steps (--dry-run)", repository.Name)
		return
	}

	if !viper.GetBool("SKIP_VITALITY") {
		// Clone repository.
		err = git.CloneRepository(repository.URL.Host, repository.Name, parser.PublicCode.URL.String(), c.Index)
		if err != nil {
			logEntries = append(logEntries, fmt.Sprintf("[%s] error while cloning: %v\n", repository.Name, err))
		}

		// Calculate Repository activity index and vitality. Defaults to 60 days.
		var activityDays int = 60
		if viper.IsSet("ACTIVITY_DAYS") {
			activityDays = viper.GetInt("ACTIVITY_DAYS")
		}
		activityIndex, _, err := git.CalculateRepoActivity(repository, activityDays)
		if err != nil {
			logEntries = append(logEntries, fmt.Sprintf("[%s] error calculating activity index: %v\n", repository.Name, err))
		} else {
			logEntries = append(logEntries, fmt.Sprintf("[%s] activity index in the last %d days: %f\n", repository.Name, activityDays, activityIndex))
		}
	}

	var aliases []string
	url := repository.CanonicalURL.String()

	// If the URL of the repo we have is different from the canonical URL
	// we got from the code hosting API, it means the repo got renamed, so we
	// add it to the slice of aliases for this software.
	if repository.URL.String() != repository.CanonicalURL.String() {
		aliases = append(aliases, repository.URL.String())
	}

	publiccodeYml, err := parser.ToYAML()
	if err != nil {
		logEntries = append(logEntries, fmt.Sprintf("[%s] parsing error: %s", repository.Name, err.Error()))

		return
	}

	if software == nil {
		_, err = c.apiClient.PostSoftware(url, aliases, string(publiccodeYml))
	} else {
		for _, alias := range software.Aliases {
			if !slices.Contains(aliases, alias) {
				aliases = append(aliases, alias)
			}
		}

		_, err = c.apiClient.PatchSoftware(software.ID, url, aliases, string(publiccodeYml))
	}
	if err != nil {
		logEntries = append(logEntries, fmt.Sprintf("[%s]: %s", repository.Name, err.Error()))
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

	// When the publisher id is a UUID, it means that the Publisher didn't originally
	// have an explicit AlternativeId, which in turn means that the Publisher
	// is not an Italian Public Administration, since those are registered in
	// the API with an alternativeId set to their iPA code (Italian PA code).
	//
	// TODO: This is not ideal and also an Italian-specific check
	// (https://github.com/italia/developers-italia-backend/issues/298)
	idIsUUID, _ := regexp.MatchString("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", publisher.Id)

	if !idIsUUID && !strings.EqualFold(
		strings.TrimSpace(publisher.Id),
		strings.TrimSpace(parser.PublicCode.It.Riuso.CodiceIPA),
	) {
		return fmt.Errorf(
			"codiceIPA is '%s', but '%s' was expected for '%s' in %s",
			parser.PublicCode.It.Riuso.CodiceIPA,
			publisher.Id,
			publisher.Name,
			fileRawURL,
		)
	}

	return nil
}
