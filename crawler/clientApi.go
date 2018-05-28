package crawler

import (
	"fmt"
	"sync"
)

// Handler returns the client handler for an organization/team/group page (every domain has a different handler implementation).
type Handler func(domain Domain, url string, repositories chan Repository, wg *sync.WaitGroup) (string, error)

// SingleHandler returns the client handler for an a single repository (every domain has a different handler implementation).
type SingleHandler func(domain Domain, url string, repositories chan Repository) error

var (
	clientAPIs      map[string]Handler
	clientSingleAPI map[string]SingleHandler
)

// RegisterClientAPIs register all the client APIs for all the clients.
func RegisterClientAPIs() {
	clientAPIs = make(map[string]Handler)
	clientSingleAPI = make(map[string]SingleHandler)

	// Client APIs for repository list.
	clientAPIs["bitbucket"] = RegisterBitbucketAPI()
	clientAPIs["github"] = RegisterGithubAPI()
	clientAPIs["gitlab"] = RegisterGitlabAPI()

	// Client APIs for a single repository.
	clientSingleAPI["bitbucket"] = RegisterSingleBitbucketAPI()
	clientSingleAPI["github"] = RegisterSingleGithubAPI()
	clientSingleAPI["gitlab"] = RegisterSingleGitlabAPI()

}

// GetClientAPICrawler checks if the API client for the requested organization clientAPI exists and return its handler.
func GetClientAPICrawler(clientAPI string) (Handler, error) {
	if crawler, ok := clientAPIs[clientAPI]; ok {
		return crawler, nil
	}
	return nil, fmt.Errorf("no client found for %s", clientAPI)

}

// GetSingleClientAPICrawler checks if the API client for the requested singlle repository clientAPI exists and return its handler.
func GetSingleClientAPICrawler(clientAPI string) (SingleHandler, error) {
	if crawler, ok := clientSingleAPI[clientAPI]; ok {
		return crawler, nil
	}
	return nil, fmt.Errorf("no single client found for %s", clientAPI)

}

// GetClients returns a list of all registered clientAPI.
func GetClients() map[string]Handler {
	return clientAPIs
}
