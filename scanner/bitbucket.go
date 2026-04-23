package scanner

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/italia/publiccode-crawler/v4/common"
	"github.com/ktrysmt/go-bitbucket"
	log "github.com/sirupsen/logrus"
)

type BitBucketScanner struct {
	client *bitbucket.Client
}

func NewBitBucketScanner() Scanner {
	client, err := bitbucket.NewBasicAuth("", "")
	if err != nil {
		panic(err)
	}

	return BitBucketScanner{client: client}
}

// ScanGroupOfRepos scans a Bitbucket workspace represented by url.
func (scanner BitBucketScanner) ScanGroupOfRepos(
	url url.URL, publisher common.Publisher, repositories chan common.Repository,
) error {
	log.Debugf("BitBucketScanner.ScanGroupOfRepos(%s)", url.String())

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
		return fmt.Errorf("can't list repositories in %s: %w", url.String(), err)
	}

	for _, item := range res.Items {
		if item.Is_private {
			log.Warnf("Skipping %s: repo is private", item.Full_name)

			continue
		}

		opt := &bitbucket.RepositoryFilesOptions{
			Owner:    owner,
			RepoSlug: item.Slug,
			Ref:      item.Mainbranch.Name,
			Path:     "publiccode.yml",
		}

		fileContent, err := scanner.client.Repositories.Repository.GetFileContent(opt)
		if err != nil {
			log.Infof("[%s]: no publiccode.yml: %s", item.Full_name, err.Error())

			continue
		}

		if fileContent != nil {
			repoURL, err := url.Parse(fmt.Sprintf("https://bitbucket.org/%s/%s.git", owner, item.Slug))
			if err != nil {
				return fmt.Errorf("failed to get canonical repo URL for %s: %w", url.String(), err)
			}

			rawURL := fmt.Sprintf(
				"https://bitbucket.org/%s/%s/raw/%s/publiccode.yml",
				owner, item.Slug, item.Mainbranch.Name,
			)
			repositories <- common.Repository{
				Name:         item.Full_name,
				FileRawURL:   rawURL,
				URL:          *repoURL,
				CanonicalURL: *repoURL,
				GitBranch:    item.Mainbranch.Name,
				Publisher:    publisher,
			}
		}
	}

	return nil
}

// ScanRepo scans a single Bitbucket repository represented by url.
func (scanner BitBucketScanner) ScanRepo(
	url url.URL, publisher common.Publisher, repositories chan common.Repository,
) error {
	log.Debugf("BitBucketScanner.ScanRepo(%s)", url.String())

	splitted := strings.Split(strings.TrimSuffix(strings.Trim(url.Path, "/"), ".git"), "/")
	if len(splitted) != 2 {
		return fmt.Errorf("bitbucket URL %s doesn't look like a repo", url.String())
	}

	owner := splitted[0]
	slug := splitted[1]

	opt := &bitbucket.RepositoryOptions{
		Owner:    owner,
		RepoSlug: slug,
	}

	repo, err := scanner.client.Repositories.Repository.Get(opt)
	if err != nil {
		return err
	}

	filesOpt := &bitbucket.RepositoryFilesOptions{
		Owner:    owner,
		RepoSlug: slug,
		Ref:      "HEAD",
		Path:     "publiccode.yml",
	}

	res, err := scanner.client.Repositories.Repository.GetFileContent(filesOpt)
	if err != nil {
		return fmt.Errorf("[%s]: no publiccode.yml: %w", url.String(), err)
	}

	if res != nil {
		canonicalURL, err := url.Parse(fmt.Sprintf("https://bitbucket.org/%s/%s.git", owner, repo.Slug))
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
