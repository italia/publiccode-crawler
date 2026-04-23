package crawler

import (
	"context"
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
	"github.com/italia/publiccode-crawler/v4/catalog"
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
	catalogsWg     sync.WaitGroup
	repositoriesWg sync.WaitGroup

	gitHubScanner    scanner.Scanner
	gitLabScanner    scanner.Scanner
	bitBucketScanner scanner.Scanner
	giteaScanner     scanner.Scanner

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
	metrics.RegisterPrometheusCounter(
		"repository_fetch_failed", "Number of repositories where fetching publiccode.yml failed (non-404)",
		c.Index,
	)

	c.gitHubScanner = scanner.NewGitHubScanner()
	c.gitLabScanner = scanner.NewGitLabScanner()
	c.bitBucketScanner = scanner.NewBitBucketScanner()
	c.giteaScanner = scanner.NewGiteaScanner()

	c.apiClient = apiclient.NewClient()

	return &c
}

// CrawlSoftwareByID crawls a single software.
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
	sourcesNum := 0
	for _, publisher := range publishers {
		sourcesNum += len(publisher.Sources)
	}

	log.Infof("Scanning %d publishers (%d catalog sources)", len(publishers), sourcesNum)

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

// ScanPublisher scans all the publisher's catalog sources and sends discovered
// repositories to the repositories channel.
func (c *Crawler) ScanPublisher(publisher common.Publisher) {
	log.Infof("Processing publisher: %s", publisher.Name)

	defer c.publishersWg.Done()

	for _, src := range publisher.Sources {
		if err := c.scanSource(src, publisher, c.repositories); err != nil {
			if errors.Is(err, scanner.ErrPubliccodeNotFound) {
				log.Warnf("[%s] %s", src.URL.String(), err.Error())
			} else {
				log.Error(err)
			}
		}
	}
}

// CrawlCatalogs processes a list of catalogs.
func (c *Crawler) CrawlCatalogs(catalogs []common.Catalog) error {
	sourcesNum := 0
	for _, cat := range catalogs {
		sourcesNum += len(cat.Sources)
	}

	log.Infof("Scanning %d catalogs (%d sources)", len(catalogs), sourcesNum)

	for _, cat := range catalogs {
		c.catalogsWg.Add(1)

		go c.ScanCatalog(cat)
	}

	go func() {
		c.catalogsWg.Wait()
		close(c.repositories)
	}()

	return c.crawl()
}

// ScanCatalog scans all sources in a catalog and sends discovered repositories
// to the repositories channel, tagging each with the catalog ID.
func (c *Crawler) ScanCatalog(cat common.Catalog) {
	log.Infof("Processing catalog: %s", cat.Name)

	defer c.catalogsWg.Done()

	const proxyChanSize = 100

	proxyCh := make(chan common.Repository, proxyChanSize)

	var proxyWg sync.WaitGroup

	proxyWg.Go(func() {
		for repo := range proxyCh {
			repo.CatalogID = cat.ID
			c.repositories <- repo
		}
	})

	publisher := common.Publisher{
		ID:   cat.ID,
		Name: cat.Name,
	}

	for _, src := range cat.Sources {
		if err := c.scanSource(src, publisher, proxyCh); err != nil {
			if errors.Is(err, scanner.ErrPubliccodeNotFound) {
				log.Warnf("[%s] %s", src.URL.String(), err.Error())
			} else {
				log.Error(err)
			}
		}
	}

	close(proxyCh)
	proxyWg.Wait()
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
	var err error

	defer func() {
		for _, e := range logEntries {
			log.Info(e)
		}

		if !c.DryRun {
			entries := strings.Join(logEntries, "\n")

			var err error

			switch {
			case repository.CatalogID != "" && software != nil:
				err = c.apiClient.PostCatalogSoftwareLog(repository.CatalogID, software.ID, entries)
			case repository.CatalogID != "":
				err = c.apiClient.PostCatalogLog(repository.CatalogID, entries)
			case software != nil:
				err = c.apiClient.PostSoftwareLog(software.ID, entries)
			default:
				err = c.apiClient.PostLog(entries)
			}

			if err != nil {
				log.Errorf("[%s]: %s", repository.Name, err.Error())
			}
		}
	}()

	// Increment counter for the number of repositories processed.
	metrics.GetCounter("repository_processed", c.Index).Inc()

	if repository.CatalogID != "" {
		software, err = c.apiClient.GetCatalogSoftwareByURL(repository.CatalogID, repository.URL.String())
	} else {
		software, err = c.apiClient.GetSoftwareByURL(repository.URL.String())
	}

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

	if software == nil {
		// New software to add
		metrics.GetCounter("repository_new", c.Index).Inc()

		if !c.DryRun {
			// Add the software even if publiccode.yml is invalid, setting active to
			// false so that we know about the new software and for example
			// [publiccode-issueopener](https://github.com/italia/publiccode-issueopener) can
			// notify maintainers about the errors.
			active := valid

			if repository.CatalogID != "" {
				software, err = c.apiClient.PostCatalogSoftware(
					repository.CatalogID, url, aliases, string(publiccodeYml), active,
				)
			} else {
				software, err = c.apiClient.PostSoftware(url, aliases, string(publiccodeYml), active)
			}
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
			if repository.CatalogID != "" {
				err = c.apiClient.PatchCatalogSoftware(
					repository.CatalogID, software.ID, url, aliases, string(publiccodeYml),
				)
			} else {
				err = c.apiClient.PatchSoftware(software.ID, url, aliases, string(publiccodeYml))
			}
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

// scanSource dispatches a single CatalogSource to the appropriate scanner.
func (c *Crawler) scanSource(
	src common.CatalogSource, publisher common.Publisher, repos chan common.Repository,
) error {
	if src.Driver == "" {
		return fmt.Errorf(
			"%s: unrecognized platform for %s, skipping",
			publisher.Name,
			src.URL.String(),
		)
	}

	switch src.Driver {
	case "github":
		if src.Group {
			return c.gitHubScanner.ScanGroupOfRepos(src.URL, publisher, repos)
		}

		return c.gitHubScanner.ScanRepo(src.URL, publisher, repos)
	case "gitlab":
		if src.Group {
			return c.gitLabScanner.ScanGroupOfRepos(src.URL, publisher, repos)
		}

		return c.gitLabScanner.ScanRepo(src.URL, publisher, repos)
	case "bitbucket":
		if src.Group {
			return c.bitBucketScanner.ScanGroupOfRepos(src.URL, publisher, repos)
		}

		return c.bitBucketScanner.ScanRepo(src.URL, publisher, repos)
	case "gitea", "forgejo":
		if src.Group {
			return c.giteaScanner.ScanGroupOfRepos(src.URL, publisher, repos)
		}

		return c.giteaScanner.ScanRepo(src.URL, publisher, repos)
	case "json":
		return c.scanJSONCatalog(src, publisher, repos)
	default:
		return fmt.Errorf(
			"%s: unknown catalog driver %q for %s",
			publisher.Name,
			src.Driver,
			src.URL.String(),
		)
	}
}

// scanJSONCatalog enumerates repository URLs from a JSON catalog and dispatches
// each one to the appropriate scanner.
func (c *Crawler) scanJSONCatalog(
	src common.CatalogSource, publisher common.Publisher, repos chan common.Repository,
) error {
	if len(src.Args) == 0 {
		return fmt.Errorf(
			"%s: json source %s is missing the JSONPath argument",
			publisher.Name,
			src.URL.String(),
		)
	}

	cat := catalog.NewJSON(src.Args[0])

	urls, err := cat.Enumerate(context.Background(), src.URL)
	if err != nil {
		return fmt.Errorf("%s: %w", publisher.Name, err)
	}

	for _, u := range urls {
		repoSrc := common.CatalogSource{
			URL:    u,
			Driver: common.InferDriver(u),
			Group:  false,
		}

		if err := c.scanSource(repoSrc, publisher, repos); err != nil {
			if errors.Is(err, scanner.ErrPubliccodeNotFound) {
				log.Warnf("[%s] %s", u.String(), err.Error())
			} else {
				log.Warnf("[%s] %s (skipping)", u.String(), err.Error())
			}
		}
	}

	return nil
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

	// When the publisher id is a UUID, it means that the Publisher didn't originally
	// have an explicit AlternativeId, which in turn means that the Publisher
	// is not an Italian Public Administration, since those are registered in
	// the API with an alternativeId set to their iPA code (Italian PA code).
	//
	// When a publisher has an alternativeId, it takes precedence over the
	// autogenerated one and it's exposed as publisher.ID.
	//
	// //nolint:godox
	// TODO: This is not ideal and also an Italian-specific check
	// (https://github.com/italia/publiccode-crawler/issues/298)
	idIsUUID, _ := regexp.MatchString("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", publisher.ID)

	var organisationURI string
	if parsed.Organisation != nil {
		organisationURI = parsed.Organisation.URI
	}

	if !idIsUUID && !strings.EqualFold(
		strings.TrimSpace("urn:x-italian-pa:"+publisher.ID),
		strings.TrimSpace(organisationURI),
	) {
		return fmt.Errorf(
			"organisation is '%s', but 'urn:x-italian-pa:%s' was expected for '%s' in %s. "+
				"Set organisation.uri to 'urn:x-italian-pa:%s'",
			organisationURI,
			publisher.ID,
			publisher.Name,
			fileRawURL,
			publisher.ID,
		)
	}

	return nil
}
