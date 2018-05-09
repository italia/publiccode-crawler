package crawler

import (
	"errors"
	"fmt"
)

var (
	clientApis      map[string]func(domain Domain, url string, repositories chan Repository) (string, error)
	clientSingleApi map[string]func(domain Domain, url string, repositories chan Repository) error
)

func RegisterClientApis() {
	clientApis = make(map[string]func(domain Domain, url string, repositories chan Repository) (string, error))
	clientSingleApi = make(map[string]func(domain Domain, url string, repositories chan Repository) error)

	clientApis["bitbucket"] = RegisterBitbucketAPI()
	clientApis["github"] = RegisterGithubAPI()
	clientApis["gitlab"] = RegisterGitlabAPI()

	clientSingleApi["bitbucket"] = RegisterSingleBitbucketAPI()
}

func GetClientApiCrawler(clientApi string) (func(domain Domain, url string, repositories chan Repository) (string, error), error) {
	if crawler, ok := clientApis[clientApi]; ok {
		return crawler, nil
	} else {
		return nil, errors.New(fmt.Sprintf("no client found for %s", clientApi))
	}
}

func GetSingleClientApiCrawler(clientApi string) (func(domain Domain, url string, repositories chan Repository) error, error) {
	if crawler, ok := clientSingleApi[clientApi]; ok {
		return crawler, nil
	} else {
		return nil, errors.New(fmt.Sprintf("no client found for %s", clientApi))
	}
}
