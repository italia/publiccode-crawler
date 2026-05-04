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
	"time"

	"github.com/alranel/go-vcsurl/v2"
	httpclient "github.com/italia/httpclient-lib-go"
	"github.com/italia/publiccode-crawler/v4/apiclient"
	"github.com/italia/publiccode-crawler/v4/common"
	"github.com/italia/publiccode-crawler/v4/git"
	"github.com/italia/publiccode-crawler/v4/metrics"
	"github.com/italia/publiccode-crawler/v4/scanner"
	publiccode "github.com/italia/publiccode-parser-go/v5"
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
	giteaScanner     scanner.Scanner

	apiClient apiclient.APIClient
}

// NewCrawler initializes a new Crawler object and connects to Elasticsearch (if dryRun == false).
func NewCrawler(dryRun bool) *Crawler {
	var crwlr Crawler

	const channelSize = 1000

	crwlr.DryRun = dryRun

	datadir := viper.GetString("DATADIR")
	if err := os.MkdirAll(datadir, 0o744); err != nil {
		log.Fatalf("can't create data directory (%s): %s", datadir, err.Error())
	}

	// Initiate a channel of repositories.
	crwlr.repositories = make(chan common.Repository, channelSize)

	// Register Prometheus metrics.
	metrics.RegisterPrometheusCounter("repository_processed", "Number of repository processed.", crwlr.Index)
	metrics.RegisterPrometheusCounter(
		"repository_good_publiccodeyml", "Number of valid publiccode.yml files in the processed repos.",
		crwlr.Index,
	)
	metrics.RegisterPrometheusCounter(
		"repository_bad_publiccodeyml", "Number of invalid publiccode.yml files in the processed repos.",
		crwlr.Index,
	)
	metrics.RegisterPrometheusCounter("repository_cloned", "Number of repository cloned", crwlr.Index)
	metrics.RegisterPrometheusCounter("repository_new", "Number of new repositories", crwlr.Index)
	metrics.RegisterPrometheusCounter("repository_known", "Number of already known repositories", crwlr.Index)
	metrics.RegisterPrometheusCounter(
		"repository_upsert_failures", "Number of failures in creating or updating software in the API",
		crwlr.Index,
	)
	metrics.RegisterPrometheusCounter(
		"repository_fetch_failed", "Number of repositories where fetching publiccode.yml failed (non-404)",
		crwlr.Index,
	)

	crwlr.gitHubScanner = scanner.NewGitHubScanner()
	crwlr.gitLabScanner = scanner.NewGitLabScanner()
	crwlr.bitBucketScanner = scanner.NewBitBucketScanner()
	crwlr.giteaScanner = scanner.NewGiteaScanner()

	crwlr.apiClient = apiclient.NewClient()

	return &crwlr
}

// CrawlSoftwareByID crawls a single software.
func (c *Crawler) CrawlSoftwareByID(software string, publisher common.Publisher) error {
	var softwareID string

	softwareURL, err := url.Parse(software)
	if err != nil {
		softwareID = software
	} else {
		softwareID = path.Base(softwareURL.Path)
	}

	softwareData, err := c.apiClient.GetSoftware(softwareID)
	if err != nil {
		return err
	}

	softwareData.URL = strings.TrimSuffix(softwareData.URL, ".git")

	repoURL, err := url.Parse(softwareData.URL)
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
	case vcsurl.IsGitea(repoURL) || vcsurl.IsForgeJo(repoURL):
		err = c.giteaScanner.ScanRepo(*repoURL, publisher, c.repositories)
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

	for _, orgEntry := range publisher.Organizations { //nolint:dupl
		orgURL := (url.URL)(orgEntry)

		switch {
		case vcsurl.IsGitHub(&orgURL):
			err = c.gitHubScanner.ScanGroupOfRepos(orgURL, publisher, c.repositories)
		case vcsurl.IsBitBucket(&orgURL):
			err = c.bitBucketScanner.ScanGroupOfRepos(orgURL, publisher, c.repositories)
		case vcsurl.IsGitLab(&orgURL):
			err = c.gitLabScanner.ScanGroupOfRepos(orgURL, publisher, c.repositories)
		case vcsurl.IsGitea(&orgURL) || vcsurl.IsForgeJo(&orgURL):
			err = c.giteaScanner.ScanGroupOfRepos(orgURL, publisher, c.repositories)
		default:
			err = fmt.Errorf(
				"publisher %s: unsupported code hosting platform for %s",
				publisher.Name,
				orgEntry.String(),
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

	for _, repoEntry := range publisher.Repositories { //nolint:dupl
		repoURL := (url.URL)(repoEntry)

		switch {
		case vcsurl.IsGitHub(&repoURL):
			err = c.gitHubScanner.ScanRepo(repoURL, publisher, c.repositories)
		case vcsurl.IsBitBucket(&repoURL):
			err = c.bitBucketScanner.ScanRepo(repoURL, publisher, c.repositories)
		case vcsurl.IsGitLab(&repoURL):
			err = c.gitLabScanner.ScanRepo(repoURL, publisher, c.repositories)
		case vcsurl.IsGitea(&repoURL) || vcsurl.IsForgeJo(&repoURL):
			err = c.giteaScanner.ScanRepo(repoURL, publisher, c.repositories)
		default:
			err = fmt.Errorf(
				"publisher %s: unsupported code hosting platform for %s",
				publisher.Name,
				repoEntry.String(),
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
func (c *Crawler) ProcessRepo(repository common.Repository) { //nolint:funlen,gocyclo
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
				err = c.apiClient.PostSoftwareLog(software.ID, entries)
			} else {
				err = c.apiClient.PostLog(entries)
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
		if resp.Status.Code == http.StatusNotFound {
			logEntries = append(logEntries, fmt.Sprintf("[%s] publiccode.yml not found (404)", repository.Name))
		} else {
			// Code -1 means all backoff retries were exhausted (usually sustained rate limiting).
			logEntries = append(logEntries, fmt.Sprintf(
				"[%s] failed to fetch publiccode.yml (HTTP %d): %v",
				repository.Name, resp.Status.Code, err,
			))
			metrics.GetCounter("repository_fetch_failed", c.Index).Inc()
		}

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

	if valid {
		logEntries = append(logEntries, fmt.Sprintf("[%s] GOOD publiccode.yml\n", repository.Name))
		metrics.GetCounter("repository_good_publiccodeyml", c.Index).Inc()
	} else {
		logEntries = append(logEntries, fmt.Sprintf("[%s] BAD publiccode.yml: %+v\n", repository.Name, err))
		metrics.GetCounter("repository_bad_publiccodeyml", c.Index).Inc()
	}

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

	var publiccodeYml []byte

	if parsed != nil {
		publiccodeYml, err = parsed.ToYAML()
		if err != nil {
			logEntries = append(logEntries, fmt.Sprintf("[%s] parsing error: %s", repository.Name, err.Error()))

			return
		}
	}

	err = c.upsertSoftware(software, url, aliases, publiccodeYml, valid)
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

		activityIndex, _, err := git.CalculateRepoActivity(repository, activityDays, time.Now())
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
	for idx := range numCPUs {
		c.repositoriesWg.Add(1)

		go func(workerID int) {
			log.Debugf("Starting ProcessRepositories() goroutine (#%d)", workerID)
			c.ProcessRepositories(reposChan)
		}(idx)
	}

	for repo := range c.repositories {
		reposChan <- repo
	}

	close(reposChan)
	c.repositoriesWg.Wait()

	fetchFailed := metrics.GetCounterValue("repository_fetch_failed", c.Index)

	summary := fmt.Sprintf(
		"Summary: Total repos scanned: %v. With good publiccode.yml file: %v. With bad publiccode.yml file: %v\n"+
			"Repos with good publiccode.yml file: New repos: %v, Known repos: %v, Failures saving to API: %v",
		metrics.GetCounterValue("repository_processed", c.Index),
		metrics.GetCounterValue("repository_good_publiccodeyml", c.Index),
		metrics.GetCounterValue("repository_bad_publiccodeyml", c.Index),
		metrics.GetCounterValue("repository_new", c.Index),
		metrics.GetCounterValue("repository_known", c.Index),
		metrics.GetCounterValue("repository_upsert_failures", c.Index),
	)

	if fetchFailed > 0 {
		summary += fmt.Sprintf(
			"\nWARNING: %v repos could not be fetched (non-404, likely rate limited or network error)"+
				" — search logs for \"failed to fetch publiccode.yml\"",
			fetchFailed,
		)
	}

	log.Info(summary)

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

	// When a Publisher has an alternativeId, it takes precedence over the
	// autogenerated UUID one and it's exposed as publisher.ID.
	//
	// //nolint:godox
	// TODO:: What if some catalogs want to have UUIDs as alternativeIDs?
	idIsUUID, _ := regexp.MatchString("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", publisher.ID)

	var organisationURI string
	if parsed.Organisation != nil {
		organisationURI = parsed.Organisation.URI
	}

	if !idIsUUID && !strings.EqualFold(
		strings.TrimSpace(publisher.ID),
		strings.TrimSpace(organisationURI),
	) {
		return fmt.Errorf(
			"organisation is '%s', but '%s' was expected for '%s' in %s. "+
				"Set organisation.uri to '%s'",
			organisationURI,
			publisher.ID,
			publisher.Name,
			fileRawURL,
			publisher.ID,
		)
	}

	return nil
}

// upsertSoftware posts or patches a software entry depending on whether it already exists.
func (c *Crawler) upsertSoftware(
	software *apiclient.Software,
	repoURL string,
	aliases []string,
	publiccodeYml []byte,
	valid bool,
) error {
	if software == nil {
		// New software to add.
		metrics.GetCounter("repository_new", c.Index).Inc()

		if c.DryRun {
			return nil
		}

		// Add the software even if publiccode.yml is invalid, setting active to
		// false so that we know about the new software and for example
		// [publiccode-issueopener](https://github.com/italia/publiccode-issueopener) can
		// notify maintainers about the errors.
		_, err := c.apiClient.PostSoftware(repoURL, aliases, string(publiccodeYml), valid)

		return err
	}

	// Known software: merge any aliases from the API that we don't already have.
	aliases = mergeNewAliases(aliases, software.Aliases)

	metrics.GetCounter("repository_known", c.Index).Inc()

	if !c.DryRun {
		return c.apiClient.PatchSoftware(software.ID, repoURL, aliases, string(publiccodeYml))
	}

	return nil
}

// mergeNewAliases appends aliases from newAliases that are not already in existing.
func mergeNewAliases(existing, newAliases []string) []string {
	for _, alias := range newAliases {
		if !slices.Contains(existing, alias) {
			existing = append(existing, alias)
		}
	}

	return existing
}
