package crawler

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"path"
	"strings"

	"golang.org/x/oauth2"

	log "github.com/sirupsen/logrus"
	"github.com/google/go-github/v43/github"
)

func githubBasicAuth(domain Domain) string {
	if len(domain.BasicAuth) > 0 {
		auth := domain.BasicAuth[rand.Intn(len(domain.BasicAuth))]
		return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	}
	return ""
}

func token(domain Domain) string {
	if len(domain.BasicAuth) > 0 {
		auth := domain.BasicAuth[rand.Intn(len(domain.BasicAuth))]

		return strings.Split(auth, ":")[1]
	}

	return ""
}

// RegisterGithubAPI register the crawler function for Github API.
// It get the list of repositories on "link" url.
// If a next page is available return its url.
// Otherwise returns an empty ("") string.
func RegisterGithubAPI() OrganizationHandler {
	return func(domain Domain, url url.URL, repositories chan Repository, publisher Publisher) (*url.URL, error) {
		// Set BasicAuth header
		headers := make(map[string]string)
		headers["Authorization"] = githubBasicAuth(domain)

		// Set domain host to new host.
		domain.Host = url.Hostname()

		ctx := context.Background()

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token(domain)},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)
		opt := &github.RepositoryListByOrgOptions{};

		orgName := strings.Split(url.Path, "/")[1]

		for {
			repos, resp, err := client.Repositories.ListByOrg(ctx, orgName, opt)
			if (err != nil) {
				log.Errorf("Can't list repositories in `%s'", orgName)
				return nil, err
			}

			// Add repositories to the channel that will perform the check on everyone.
			for _, r := range repos {
				if *r.Private || *r.Archived {
					log.Warnf("Skipping %s: repo is private or archived", *r.FullName)
					continue
				}

				file, _, _, err := client.Repositories.GetContents(ctx, orgName, *r.Name, "publiccode.yml", nil)
				if (err != nil) {
					log.Infof("[%s]: no publiccode.yml", *r.FullName)
					continue
				}
				if file != nil {
					repositories <- Repository{
						Name:        *r.FullName,
						Hostname:    domain.Host,
						FileRawURL:  *file.DownloadURL,
						GitCloneURL: *r.CloneURL,
						GitBranch:   *r.DefaultBranch,
						Domain:      domain,
						Publisher:   publisher,
						Headers:     headers,
					}
				}
			}

			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}

		return nil, nil
	}
}

// RegisterSingleGithubAPI register the crawler function for single repository Github API.
// Return nil if the repository was successfully added to repositories channel.
// Otherwise return the generated error.
func RegisterSingleGithubAPI() SingleRepoHandler {
	return func(domain Domain, url url.URL, repositories chan Repository, publisher Publisher) error {
		// Set BasicAuth header.
		headers := make(map[string]string)
		headers["Authorization"] = githubBasicAuth(domain)

		// Set domain host to new host.
		domain.Host = url.Hostname()

		ctx := context.Background()

		// token := strings.Split(url.Path, "/")[1]
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token(domain)},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)

		splitted := strings.Split(strings.Trim(url.Path, "/"), "/")
		orgName := splitted[0]
		repoName := splitted[1]

		repo, _, err := client.Repositories.Get(ctx, orgName, repoName)
		if (err != nil) {
			log.Errorf("Can't get repo `%s/%s'", orgName, repoName)
			return err
		}

		if *repo.Private || *repo.Archived {
			log.Warnf("Skipping %s: repo is private or archived", *repo.FullName)
			return errors.New("Skipping private or archived repo")
		}

		file, _, _, err := client.Repositories.GetContents(ctx, orgName, *repo.Name, "publiccode.yml", nil)
		if (err != nil) {
			return fmt.Errorf("[%s]: failed to get contents", *repo.FullName)
		}
		if file != nil {
			repositories <- Repository{
				Name:        *repo.FullName,
				Hostname:    domain.Host,
				FileRawURL:  *file.DownloadURL,
				GitCloneURL: *repo.CloneURL,
				GitBranch:   *repo.DefaultBranch,
				Domain:      domain,
				Publisher:   publisher,
				Headers:     headers,
			}
		}
		return nil
	}
}

// GenerateGithubAPIURL returns the api url of given Gitlab organization link.
// IN: https://github.com/italia
// OUT:https://api.github.com/orgs/italia/repos,https://api.github.com/users/italia/repos
func GenerateGithubAPIURL() GeneratorAPIURL {
	return func(in url.URL) (out []url.URL, err error) {
		u := in
		u.Path = path.Join("orgs", u.Path, "repos")
		u.Path = strings.Trim(u.Path, "/")
		u.Host = "api." + u.Host
		out = append(out, u)

		u2 := in
		u2.Path = path.Join("users", u2.Path, "repos")
		u2.Path = strings.Trim(u2.Path, "/")
		u2.Host = "api." + u2.Host
		out = append(out, u2)

		return
	}
}
