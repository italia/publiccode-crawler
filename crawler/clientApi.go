package crawler

import (
	"errors"
	"fmt"
	"sync"
)

type Handler func(domain Domain, url string, repositories chan Repository, wg *sync.WaitGroup) (string, error)
type SingleHandler func(domain Domain, url string, repositories chan Repository) error

var (
	clientApis      map[string]Handler
	clientSingleApi map[string]SingleHandler
)

// RegisterClientApis register all the client APIs for all the clients.
func RegisterClientApis() {
	clientApis = make(map[string]Handler)
	clientSingleApi = make(map[string]SingleHandler)

	// Client APIs for repository list.
	clientApis["bitbucket"] = RegisterBitbucketAPI()
	clientApis["github"] = RegisterGithubAPI()
	clientApis["gitlab"] = RegisterGitlabAPI()

	// Client APIs for a single repository.
	clientSingleApi["bitbucket"] = RegisterSingleBitbucketAPI()
	clientSingleApi["github"] = RegisterSingleGithubAPI()
	clientSingleApi["gitlab"] = RegisterSingleGitlabAPI()

}

// GetClientApiCrawler checks if the api client for the requested organization clientApi exists and return its handler.
func GetClientApiCrawler(clientApi string) (Handler, error) {
	if crawler, ok := clientApis[clientApi]; ok {
		return crawler, nil
	} else {
		return nil, errors.New(fmt.Sprintf("no client found for %s", clientApi))
	}
}

// GetSingleClientApiCrawler checks if the api client for the requested singlle repository clientApi exists and return its handler.
func GetSingleClientApiCrawler(clientApi string) (SingleHandler, error) {
	if crawler, ok := clientSingleApi[clientApi]; ok {
		return crawler, nil
	} else {
		return nil, errors.New(fmt.Sprintf("no single client found for %s", clientApi))
	}
}

// GetClients returns a list of all registered clientApi.
func GetClients() map[string]Handler {
	return clientApis
}
