package scanner

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v43/github"
	log "github.com/sirupsen/logrus"

	"github.com/italia/publiccode-crawler/v3/common"
)

type GitHubScanner struct {
	client *github.Client
	ctx    context.Context
}

// NewGitHubScanner returns a new GitHubScanner using the
// authentication token from the GITHUB_TOKEN environment variable or,
// if not set, the tokens in domains.yml.
func NewGitHubScanner() Scanner {
	ctx := context.Background()

	token := os.Getenv("GITHUB_TOKEN")

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return GitHubScanner{client: client, ctx: ctx}
}

// ScanGroupOfRepos scans a GitHub organization represented by url, associated to
// publisher and sends any repository containing a publiccode.yml to the repositories
// channel as a [common.Repository].
// It returns any error encountered if any, otherwise nil.
func (scanner GitHubScanner) ScanGroupOfRepos(url url.URL, publisher common.Publisher, repositories chan common.Repository) error {
	opt := &github.RepositoryListByOrgOptions{}

	splitted := strings.Split(strings.Trim(url.Path, "/"), "/")
	if len(splitted) != 1 {
		return fmt.Errorf("doesn't look like a GitHub org %s", url.String())
	}

	orgName := splitted[0]

	for {
Retry:
		repos, resp, err := scanner.client.Repositories.ListByOrg(scanner.ctx, orgName, opt)
		if _, ok := err.(*github.RateLimitError); ok {
			log.Infof("GitHub rate limit hit, sleeping until %s", resp.Rate.Reset.Time.String())
			time.Sleep(time.Until(resp.Rate.Reset.Time))
			goto Retry
		} else if err != nil {
			return fmt.Errorf("can't list repositories in %s (not an GitHub organization?): %w", url.String(), err)
		}

		// Add repositories to the channel that will perform the check on everyone.
		for _, r := range repos {
			repoURL, err := url.Parse(*r.HTMLURL)
			if err != nil {
				log.Errorf("can't parse URL %s: %s", *r.URL, err.Error())
				continue
			}

			if err = scanner.ScanRepo(*repoURL, publisher, repositories); err != nil {
				if errors.Is(err, ErrPubliccodeNotFound) {
					log.Warnf("can't scan repository %s: %s", repoURL.String(), err.Error())
				} else {
					log.Errorf("can't scan repository %s: %s", repoURL.String(), err.Error())
				}

				continue
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil
}

// ScanRepo scans a GitHub repository represented by url, associated to
// publisher and, if it contains a publiccode.yml, sends it as a [common.Repository]
// repositories channel.
// It returns any error encountered if any, otherwise nil.
func (scanner GitHubScanner) ScanRepo(url url.URL, publisher common.Publisher, repositories chan common.Repository) error {
	splitted := strings.Split(strings.Trim(url.Path, "/"), "/")
	if len(splitted) != 2 {
		return fmt.Errorf("doesn't look like a GitHub repo %s", url.String())
	}

	orgName := splitted[0]
	repoName := splitted[1]

Retry:
	repo, resp, err := scanner.client.Repositories.Get(scanner.ctx, orgName, repoName)
	if _, ok := err.(*github.RateLimitError); ok {
		log.Infof("GitHub rate limit hit, sleeping until %s", resp.Rate.Reset.Time.String())
		time.Sleep(time.Until(resp.Rate.Reset.Time))
		goto Retry
	}
	if err != nil {
		return fmt.Errorf("can't get repo %s: %w", url.String(), err)
	}

	if *repo.Private || *repo.Archived {
		return fmt.Errorf("skipping private or archived repo %s", *repo.FullName)
	}

	file, _, resp, err := scanner.client.Repositories.GetContents(scanner.ctx, orgName, repoName, "publiccode.yml", nil)
	if _, ok := err.(*github.RateLimitError); ok {
		log.Infof("GitHub rate limit hit, sleeping until %s", resp.Rate.Reset.Time.String())
		time.Sleep(time.Until(resp.Rate.Reset.Time))
		goto Retry
	}

	if err != nil {
		if resp.StatusCode == 404 {
			return ErrPubliccodeNotFound
		}

		return fmt.Errorf("[%s]: failed to get publiccode.yml: %w", *repo.FullName, err)
	}
	if file != nil {
		canonicalURL, err := url.Parse(*repo.CloneURL)
		if err != nil {
			return fmt.Errorf("failed to get canonical repo URL for %s: %w", url.String(), err)
		}

		repositories <- common.Repository{
			Name:         *repo.FullName,
			FileRawURL:   *file.DownloadURL,
			URL:          url,
			CanonicalURL: *canonicalURL,
			GitBranch:    *repo.DefaultBranch,
			Publisher:    publisher,
			Headers:      make(map[string]string),
		}
	}

	return nil
}
