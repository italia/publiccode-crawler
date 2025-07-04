package crawler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/alranel/go-vcsurl/v2"
	httpclient "github.com/italia/httpclient-lib-go"
	"github.com/italia/publiccode-crawler/v4/apiclient"
	"github.com/italia/publiccode-crawler/v4/common"
	"github.com/italia/publiccode-crawler/v4/git"
	"github.com/italia/publiccode-crawler/v4/metrics"
	"github.com/italia/publiccode-crawler/v4/scanner"
	publiccode "github.com/italia/publiccode-parser-go/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Crawler is a helper class representing a crawler.
type Crawler struct {
	DryRun bool

	Index        string
	repositories chan common.Repository
	// Sync mutex guard.
	publishersWg   sync.WaitGroup
	repositoriesWg sync.WaitGroup

	gitHubScanner    scanner.Scanner
	gitLabScanner    scanner.Scanner
	bitBucketScanner scanner.Scanner

	apiClient apiclient.APIClient
}

// NewCrawler initializes a new Crawler object and connects to Elasticsearch (if dryRun == false).
func NewCrawler(dryRun bool) *Crawler {
	var c Crawler

	const channelSize = 1000

	c.DryRun = dryRun

	datadir := viper.GetString("DATADIR")
	if err := os.MkdirAll(datadir, 0o744); err != nil {
		log.Fatalf("can't create data directory (%s): %s", datadir, err.Error())
	}

	// Initiate a channel of repositories.
	c.repositories = make(chan common.Repository, channelSize)

	// Register Prometheus metrics.
	metrics.RegisterPrometheusCounter("repository_processed", "Number of repository processed.", c.Index)
	metrics.RegisterPrometheusCounter(
		"repository_good_publiccodeyml", "Number of valid publiccode.yml files in the processed repos.",
		c.Index,
	)
	metrics.RegisterPrometheusCounter(
		"repository_bad_publiccodeyml", "Number of invalid publiccode.yml files in the processed repos.",
		c.Index,
	)
	metrics.RegisterPrometheusCounter("repository_cloned", "Number of repository cloned", c.Index)
	metrics.RegisterPrometheusCounter("repository_new", "Number of new repositories", c.Index)
	metrics.RegisterPrometheusCounter("repository_known", "Number of already known repositories", c.Index)
	metrics.RegisterPrometheusCounter(
		"repository_upsert_failures", "Number of failures in creating or updating software in the API",
		c.Index,
	)

	c.gitHubScanner = scanner.NewGitHubScanner()
	c.gitLabScanner = scanner.NewGitLabScanner()
	c.bitBucketScanner = scanner.NewBitBucketScanner()

	c.apiClient = apiclient.NewClient()

	return &c
}

// CrawlSoftwareByAPIURL crawls a single software.
func (c *Crawler) CrawlSoftwareByID(software string, publisher common.Publisher) error {
	var id string

	softwareURL, err := url.Parse(software)
	if err != nil {
		id = software
	} else {
		id = path.Base(softwareURL.Path)
	}

	s, err := c.apiClient.GetSoftware(id)
	if err != nil {
		return err
	}

	s.URL = strings.TrimSuffix(s.URL, ".git")

	repoURL, err := url.Parse(s.URL)
	if err != nil {
		return err
	}

	log.Infof("Processing repository: %s", softwareURL.String())

	switch {
	case vcsurl.IsGitHub(repoURL):
		err = c.gitHubScanner.ScanRepo(*repoURL, publisher, c.repositories)
	case vcsurl.IsBitBucket(repoURL):
		err = c.bitBucketScanner.ScanRepo(*repoURL, publisher, c.repositories)
	case vcsurl.IsGitLab(repoURL):
		err = c.gitLabScanner.ScanRepo(*repoURL, publisher, c.repositories)
	default:
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
		c.ScanPublisher(publisher)
	}

	// Close the repositories channel when all the publisher goroutines are done
	go func() {
		c.publishersWg.Wait()
		close(c.repositories)
	}()

	return c.crawl()
}

// ScanPublisher scans all the publisher' repositories and sends the ones
// with a valid publiccode.yml to the repositories channel.
func (c *Crawler) ScanPublisher(publisher common.Publisher) {
	log.Infof("Processing publisher: %s", publisher.Name)

	defer c.publishersWg.Done()

	var err error

	for _, u := range publisher.Organizations {
		orgURL := (url.URL)(u)

		switch {
		case vcsurl.IsGitHub(&orgURL):
			err = c.gitHubScanner.ScanGroupOfRepos(orgURL, publisher, c.repositories)
		case vcsurl.IsBitBucket(&orgURL):
			err = c.bitBucketScanner.ScanGroupOfRepos(orgURL, publisher, c.repositories)
		case vcsurl.IsGitLab(&orgURL):
			err = c.gitLabScanner.ScanGroupOfRepos(orgURL, publisher, c.repositories)
		default:
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

		switch {
		case vcsurl.IsGitHub(&repoURL):
			err = c.gitHubScanner.ScanRepo(repoURL, publisher, c.repositories)
		case vcsurl.IsBitBucket(&repoURL):
			err = c.bitBucketScanner.ScanRepo(repoURL, publisher, c.repositories)
		case vcsurl.IsGitLab(&repoURL):
			err = c.gitLabScanner.ScanRepo(repoURL, publisher, c.repositories)
		default:
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
func (c *Crawler) ProcessRepo(repository common.Repository) { //nolint:maintidx
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
		logEntries = append(
			logEntries,
			fmt.Sprintf("[%s] failed to GET software from API: %s\n", repository.Name, err.Error()),
		)

		return
	}

	// We don't want to re-activate software that was de-activated manually, but just
	// if it was previously added automatically for the first time as inactive (to keep
	// track of errors in new software - https://github.com/italia/publiccode-crawler/issues/325).
	//
	// When CreatedAt != UpdatedAt it is most likely a manual deactivation for a good reason and
	// we don't wan't to re-enable in that case.
	if software != nil && !software.Active && software.CreatedAt != software.UpdatedAt {
		logEntries = append(
			logEntries,
			fmt.Sprintf(
				`[%s] software %s has "active": false and "created_at" is different than "updated_at. `+
					`This means the software was deactivated manually, skipping update.`, repository.Name, software.ID,
			),
		)

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

	//nolint:godox
	// FIXME: this is hardcoded for now, because it requires changes to publiccode-parser-go.
	domain := publiccode.Domain{
		Host:        "github.com",
		UseTokenFor: []string{"github.com", "api.github.com", "raw.githubusercontent.com"},
		BasicAuth:   []string{os.Getenv("GITHUB_TOKEN")},
	}

	var parser *publiccode.Parser

	parser, err = publiccode.NewParser(publiccode.ParserConfig{Domain: domain})
	if err != nil {
		logEntries = append(
			logEntries,
			fmt.Sprintf("[%s] can't create a Parser: %s\n", repository.Name, err.Error()),
		)

		return
	}

	var parsed publiccode.PublicCode
	parsed, err = parser.Parse(repository.FileRawURL)

	valid := true

	if err != nil {
		var validationResults publiccode.ValidationResults
		if errors.As(err, &validationResults) {
			var validationError publiccode.ValidationError
			for _, res := range validationResults {
				if errors.As(res, &validationError) {
					valid = false

					break
				}
			}
		}
	}

	publisherID := viper.GetString("MAIN_PUBLISHER_ID")
	if valid && repository.Publisher.ID != publisherID {
		//nolint:forcetypeassert // we'd want to panic here anyway if the library returns a non v0
		err = validateFile(repository.Publisher, parsed.(publiccode.PublicCodeV0), repository.FileRawURL)
		if err != nil {
			valid = false
		}
	}

	if !valid {
		logEntries = append(logEntries, fmt.Sprintf("[%s] BAD publiccode.yml: %+v\n", repository.Name, err))

		metrics.GetCounter("repository_bad_publiccodeyml", c.Index).Inc()

		return
	}

	logEntries = append(logEntries, fmt.Sprintf("[%s] GOOD publiccode.yml\n", repository.Name))

	metrics.GetCounter("repository_good_publiccodeyml", c.Index).Inc()

	if c.DryRun {
		log.Infof("[%s]: Skipping other steps (--dry-run)", repository.Name)
	}

	var aliases []string

	url := repository.CanonicalURL.String()

	// If the URL of the repo we have is different from the canonical URL
	// we got from the code hosting API, it means the repo got renamed, so we
	// add it to the slice of aliases for this software.
	if repository.URL.String() != repository.CanonicalURL.String() {
		aliases = append(aliases, repository.URL.String())
	}

	publiccodeYml, err := parsed.ToYAML()
	if err != nil {
		logEntries = append(logEntries, fmt.Sprintf("[%s] parsing error: %s", repository.Name, err.Error()))

		return
	}

	if software == nil {
		// New software to add
		metrics.GetCounter("repository_new", c.Index).Inc()

		if !c.DryRun {
			// Add the software even if publiccode.yml is invalid, setting active to
			// false so that we know about the new software and for example
			// [publiccode-issueopener](https://github.com/italia/publiccode-issueopener) can
			// notify maintainers about the errors.
			active := valid

			software, err = c.apiClient.PostSoftware(url, aliases, string(publiccodeYml), active)
			if err != nil {
				return
			}
		}
	} else {
		// Known software
		for _, alias := range software.Aliases {
			if !slices.Contains(aliases, alias) {
				aliases = append(aliases, alias)
			}
		}

		metrics.GetCounter("repository_known", c.Index).Inc()

		if !c.DryRun {
			_, err = c.apiClient.PatchSoftware(software.ID, url, aliases, string(publiccodeYml))
		}
	}

	if err != nil {
		logEntries = append(logEntries, fmt.Sprintf("[%s]: %s", repository.Name, err.Error()))

		metrics.GetCounter("repository_upsert_failures", c.Index).Inc()
	}

	if !viper.GetBool("SKIP_VITALITY") && !c.DryRun {
		// Clone repository.
		err = git.CloneRepository(repository.URL.Host, repository.Name, parsed.Url().String(), c.Index)
		if err != nil {
			logEntries = append(logEntries, fmt.Sprintf("[%s] error while cloning: %v\n", repository.Name, err))
		}

		// Calculate Repository activity index and vitality. Defaults to 60 days.
		activityDays := 60
		if viper.IsSet("ACTIVITY_DAYS") {
			activityDays = viper.GetInt("ACTIVITY_DAYS")
		}

		activityIndex, _, err := git.CalculateRepoActivity(repository, activityDays)
		if err != nil {
			logEntries = append(
				logEntries, fmt.Sprintf("[%s] error calculating activity index: %v\n", repository.Name, err),
			)
		} else {
			logEntries = append(
				logEntries,
				fmt.Sprintf("[%s] activity index in the last %d days: %f\n", repository.Name, activityDays, activityIndex),
			)
		}
	}
}

func (c *Crawler) crawl() error {
	reposChan := make(chan common.Repository)

	// Start the metrics server.
	go metrics.StartPrometheusMetricsServer()

	defer c.publishersWg.Wait()

	// Get cpus number
	numCPUs := runtime.NumCPU()
	log.Debugf("CPUs #: %d", numCPUs)

	// Process the repositories in order to retrieve the files.
	for i := range numCPUs {
		c.repositoriesWg.Add(1)

		go func(id int) {
			log.Debugf("Starting ProcessRepositories() goroutine (#%d)", id)
			c.ProcessRepositories(reposChan)
		}(i)
	}

	for repo := range c.repositories {
		reposChan <- repo
	}

	close(reposChan)
	c.repositoriesWg.Wait()

	log.Infof(
		"Summary: Total repos scanned: %v. With good publiccode.yml file: %v. With bad publiccode.yml file: %v\n"+
			"Repos with good publiccode.yml file: New repos: %v, Known repos: %v, Failures saving to API: %v",
		metrics.GetCounterValue("repository_processed", c.Index),
		metrics.GetCounterValue("repository_good_publiccodeyml", c.Index),
		metrics.GetCounterValue("repository_bad_publiccodeyml", c.Index),
		metrics.GetCounterValue("repository_new", c.Index),
		metrics.GetCounterValue("repository_known", c.Index),
		metrics.GetCounterValue("repository_upsert_failures", c.Index),
	)

	return nil
}

// validateFile performs additional validations that are not strictly mandated
// by the publiccode.yml Standard.
// Using `one` command this check will be skipped.
func validateFile(publisher common.Publisher, parsed publiccode.PublicCodeV0, fileRawURL string) error {
	u, _ := url.Parse(fileRawURL)
	repo1 := vcsurl.GetRepo(u)

	repo2 := vcsurl.GetRepo((*url.URL)(parsed.Url()))

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
				parsed.Url(), fileRawURL, repo2, repo1,
			)
		}
	}

	// When the publisher id is a UUID, it means that the Publisher didn't originally
	// have an explicit AlternativeId, which in turn means that the Publisher
	// is not an Italian Public Administration, since those are registered in
	// the API with an alternativeId set to their iPA code (Italian PA code).
	//
	// //nolint:godox
	// TODO: This is not ideal and also an Italian-specific check
	// (https://github.com/italia/publiccode-crawler/issues/298)
	idIsUUID, _ := regexp.MatchString("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", publisher.ID)

	if !idIsUUID && !strings.EqualFold(
		strings.TrimSpace(publisher.ID),
		strings.TrimSpace(parsed.It.Riuso.CodiceIPA),
	) {
		return fmt.Errorf(
			"codiceIPA is '%s', but '%s' was expected for '%s' in %s",
			parsed.It.Riuso.CodiceIPA,
			publisher.ID,
			publisher.Name,
			fileRawURL,
		)
	}

	return nil
}
