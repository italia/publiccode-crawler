package crawler

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strings"

	httpclient "github.com/italia/httpclient-lib-go"
	"github.com/ktrysmt/go-bitbucket"
	log "github.com/sirupsen/logrus"
)

// RegisterBitbucketAPI register the crawler function for Bitbucket API.
func RegisterBitbucketAPI() OrganizationHandler {
	return func(domain Domain, url url.URL, repositories chan Repository, publisher Publisher) (*url.URL, error) {
		// Set BasicAuth header.
		headers := make(map[string]string)

		token := ""
		if domain.BasicAuth != nil {
			token = domain.BasicAuth[rand.Intn(len(domain.BasicAuth))]
			// TODO: refactor in order to not need to pass Headers around
			headers["Authorization"] = token
		}

		// Set domain host to new host.
		domain.Host = url.Hostname()


		team := strings.Split(url.Path, "/")[3]
		opt := &bitbucket.RepositoriesOptions{
			Owner: team,
		}

		client := bitbucket.NewBasicAuth("", "")
		res, err := client.Repositories.ListForAccount(opt)

		if err != nil {
			return nil, fmt.Errorf("Can't list repositories in `%s'", team)
		}

		// Add repositories to the channel that will perform the check on everyone.
		for _, r := range res.Items {
			if r.Is_private {
				log.Warnf("Skipping %s: repo is private", r.Full_name)
				continue
			}

			opt := &bitbucket.RepositoryFilesOptions {
				Owner: team,
				RepoSlug: r.Slug,
				Ref: r.Mainbranch.Name,
				Path: "publiccode.yml",
			}
			res, err := client.Repositories.Repository.GetFileContent(opt)
			if (err != nil) {
				log.Infof("[%s]: no publiccode.yml: %s", r.Full_name, err.Error())
				continue
			}
			if res != nil {
				repositories <- Repository{
					Name:        r.Full_name,
					Hostname:    domain.Host,
					FileRawURL:  fmt.Sprintf("https://bitbucket.org/%s/%s/raw/%s/publiccode.yml", team,r.Slug, r.Mainbranch.Name),
					GitCloneURL: fmt.Sprintf("https://bitbucket.org/%s/%s.git", team, r.Slug),
					GitBranch:   r.Mainbranch.Name,
					Domain:      domain,
					Publisher:   publisher,
					Headers:     headers,
				}
			}
		}

		return nil, nil
	}
}

// RegisterSingleBitbucketAPI register the crawler function for single Bitbucket repository.
func RegisterSingleBitbucketAPI() SingleRepoHandler {
	return func(domain Domain, url url.URL, repositories chan Repository, publisher Publisher) error {
		// Set BasicAuth header
		headers := make(map[string]string)

		token := ""
		if domain.BasicAuth != nil {
			token = domain.BasicAuth[rand.Intn(len(domain.BasicAuth))]
			// TODO: refactor in order to not need to pass Headers around
			headers["Authorization"] = token
		}

		// Set domain host to new host.
		domain.Host = url.Hostname()

		repoName := strings.Split(strings.Trim(url.Path, "/"), "/")
		opt := &bitbucket.RepositoryOptions{
			Owner: repoName[0],
			RepoSlug: repoName[1],
		}

		client := bitbucket.NewBasicAuth("", "")

		repo , err := client.Repositories.Repository.Get(opt)
		if err != nil {
			return err
		}

		filesOpt := &bitbucket.RepositoryFilesOptions {
			Owner: repoName[0],
			RepoSlug: repoName[1],
			Path: "publiccode.yml",
		}
		res, err := client.Repositories.Repository.GetFileContent(filesOpt)
		if (err != nil) {
			return err
		}
		if res != nil {
			repositories <- Repository{
				Name:        repo.Full_name,
				Hostname:    domain.Host,
				FileRawURL:  fmt.Sprintf("https://bitbucket.org/%s/%s/raw/%s/publiccode.yml", repoName[0], repoName[1], repo.Mainbranch.Name),
				GitCloneURL: fmt.Sprintf("https://bitbucket.org/%s/%s.git", repoName[0], repoName[1]),
				GitBranch:   repo.Mainbranch.Name,
				Domain:      domain,
				Publisher:   publisher,
				Headers:     headers,
			}
		}

		return nil
	}
}

// GenerateBitbucketAPIURL returns the api url of given Bitbucket  organization link.
// IN: https://bitbucket.org/Soft
// OUT:https://api.bitbucket.org/2.0/repositories/Soft?pagelen=100
func GenerateBitbucketAPIURL() GeneratorAPIURL {
	return func(u url.URL) (out []url.URL, err error) {
		u.Path = path.Join("/2.0/repositories", u.Path)
		u.Host = "api." + u.Host

		out = append(out, u)
		return
	}
}
