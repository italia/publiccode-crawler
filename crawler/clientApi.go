package crawler

import (
	"fmt"
	"sync"
)

// Handler returns the client handler for an organization/team/group page (every domain has a different handler implementation).
type Handler func(domain Domain, url string, repositories chan Repository, wg *sync.WaitGroup) (string, error)

// SingleHandler returns the client handler for an a single repository (every domain has a different handler implementation).
type SingleHandler func(domain Domain, url string, repositories chan Repository) error

// SingleHandler returns the client handler for an a single repository (every domain has a different handler implementation).
type GeneratorURL func(url string) (string, error)

var (
	clientAPIs       map[string]Handler
	clientSingleAPIs map[string]SingleHandler
	clientAPIURLs    map[string]GeneratorURL
)

// RegisterClientAPIs register all the client APIs for all the clients.
func RegisterClientAPIs() {
	clientAPIs = make(map[string]Handler)
	clientSingleAPIs = make(map[string]SingleHandler)
	clientAPIURLs = make(map[string]GeneratorURL)

	// Client APIs for repository list.
	clientAPIs["bitbucket.org"] = RegisterBitbucketAPI()
	clientAPIs["github.com"] = RegisterGithubAPI()
	clientAPIs["gitlab.com"] = RegisterGitlabAPI()

	// Client APIs for a single repository.
	clientSingleAPIs["bitbucket.org"] = RegisterSingleBitbucketAPI()
	clientSingleAPIs["github.com"] = RegisterSingleGithubAPI()
	clientSingleAPIs["gitlab.com"] = RegisterSingleGitlabAPI()

	// Client APIs for a single repository.
	clientAPIURLs["bitbucket.org"] = GenerateBitbucketAPIURL()
	clientAPIURLs["github.com"] = GenerateGithubAPIURL()
	clientAPIURLs["gitlab.com"] = GenerateGitlabAPIURL()

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
	if crawler, ok := clientSingleAPIs[clientAPI]; ok {
		return crawler, nil
	}
	return nil, fmt.Errorf("no single client found for %s", clientAPI)
}

// GetAPIURL
func GetAPIURL(clientAPI string) (GeneratorURL, error) {
	if crawler, ok := clientAPIURLs[clientAPI]; ok {
		return crawler, nil
	}
	return nil, fmt.Errorf("no client found for %s", clientAPI)
}

// GetClients returns a list of all registered clientAPI.
func GetClients() map[string]Handler {
	return clientAPIs
}
