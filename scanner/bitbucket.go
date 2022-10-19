package scanner

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/italia/publiccode-crawler/v3/common"
	"github.com/ktrysmt/go-bitbucket"
	log "github.com/sirupsen/logrus"
)

type BitBucketScanner struct {
	client *bitbucket.Client
}

func NewBitBucketScanner() Scanner {
	return BitBucketScanner{client: bitbucket.NewBasicAuth("", "")}
}

// RegisterBitbucketAPI register the crawler function for Bitbucket API.
func (scanner BitBucketScanner) ScanGroupOfRepos(url url.URL, publisher common.Publisher, repositories chan common.Repository) error {
	splitted := strings.Split(strings.Trim(url.Path, "/"), "/")

	if len(splitted) != 1 {
		return fmt.Errorf("bitbucket URL %s doesn't look like a group of repos", url.String())
	}

	owner := splitted[0]

	opt := &bitbucket.RepositoriesOptions{
		Owner: owner,
	}

	res, err := scanner.client.Repositories.ListForAccount(opt)

	if err != nil {
		return fmt.Errorf("Can't list repositories in %s: %w", url.String(), err)
	}

	for _, r := range res.Items {
		if r.Is_private {
			log.Warnf("Skipping %s: repo is private", r.Full_name)
			continue
		}

		opt := &bitbucket.RepositoryFilesOptions {
			Owner: owner,
			RepoSlug: r.Slug,
			Ref: r.Mainbranch.Name,
			Path: "publiccode.yml",
		}
		res, err := scanner.client.Repositories.Repository.GetFileContent(opt)
		if (err != nil) {
			log.Infof("[%s]: no publiccode.yml: %s", r.Full_name, err.Error())
			continue
		}
		if res != nil {
			u, err := url.Parse(fmt.Sprintf("https://bitbucket.org/%s/%s.git", owner, r.Slug))
			if err != nil {
				return fmt.Errorf("failed to get canonical repo URL for %s: %w", url.String(), err)
			}

			repositories <- common.Repository{
				Name:         r.Full_name,
				FileRawURL:   fmt.Sprintf("https://bitbucket.org/%s/%s/raw/%s/publiccode.yml", owner,r.Slug, r.Mainbranch.Name),
				URL:          *u,
				CanonicalURL: *u,
				GitBranch:    r.Mainbranch.Name,
				Publisher:    publisher,
			}
		}
	}

	return nil
}

// RegisterSingleBitbucketAPI register the crawler function for single Bitbucket repository.
func (scanner BitBucketScanner) ScanRepo(url url.URL, publisher common.Publisher, repositories chan common.Repository) error {
	splitted := strings.Split(strings.Trim(url.Path, "/"), "/")
	if len(splitted) != 2 {
		return fmt.Errorf("bitbucket URL %s doesn't look like a repo", url.String())
	}

	owner := splitted[0]
	slug := splitted[1]

	opt := &bitbucket.RepositoryOptions{
		Owner: owner,
		RepoSlug: slug,
	}

	repo , err := scanner.client.Repositories.Repository.Get(opt)
	if err != nil {
		return err
	}

	filesOpt := &bitbucket.RepositoryFilesOptions {
		Owner: owner,
		RepoSlug: slug,
		Path: "publiccode.yml",
	}
	res, err := scanner.client.Repositories.Repository.GetFileContent(filesOpt)
	if (err != nil) {
		return err
	}
	if res != nil {
		canonicalURL, err := url.Parse(fmt.Sprintf("https://bitbucket.org/%s/%s.git", repo.Owner, repo.Slug))
		if err != nil {
			return fmt.Errorf("failed to get canonical repo URL for %s: %w", url.String(), err)
		}

		repositories <- common.Repository{
			Name:         repo.Full_name,
			FileRawURL:   fmt.Sprintf("https://bitbucket.org/%s/%s/raw/%s/publiccode.yml", owner, slug, repo.Mainbranch.Name),
			URL:          url,
			CanonicalURL: *canonicalURL,
			GitBranch:    repo.Mainbranch.Name,
			Publisher:    publisher,
		}
	}

	return nil
}
