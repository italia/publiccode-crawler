package crawler

import (
	"errors"
	"fmt"
	"sync"
)

type Handler func(domain Domain, url string, repositories chan Repository, wg *sync.WaitGroup) (string, error)

// Crawler is the interface for crawler plugins.
type Crawler interface {
	Register() Handler
	GetId() string
}

var (
	clientApis map[string]Handler
)

// RegisterCrawlers registers all founded crawler plugins.
func RegisterClientApis() {
	clientApis = make(map[string]Handler)

	clientApis["bitbucket"] = RegisterBitbucketAPI()
	clientApis["gitlab"] = RegisterGitlabAPI()
	clientApis["github"] = RegisterGithubAPI()
}

// GetClientApiCrawler returns the handler func to process domain.
func GetClientApiCrawler(clientApi string) (Handler, error) {
	if crawler, ok := clientApis[clientApi]; ok {
		return crawler, nil
	} else {
		return nil, errors.New(fmt.Sprintf("no client found for %s", clientApi))
	}
}

// GetClients returns a list of all registered plugins.
func GetClients() map[string]Handler {
	return clientApis
}
