package crawler

import (
	"fmt"
	"net/url"
)

// ClientAPI contains all the API function in a single Client.
type ClientAPI struct {
	Organization OrganizationHandler
	Single       SingleRepoHandler

	APIURL GeneratorAPIURL
}

// OrganizationHandler returns the client handler for an organization/team/group page (every domain has a different handler implementation).
type OrganizationHandler func(domain Domain, url url.URL, repositories chan Repository, publisher Publisher) (*url.URL, error)

// SingleRepoHandler returns the client handler for an a single repository (every domain has a different handler implementation).
type SingleRepoHandler func(domain Domain, url url.URL, repositories chan Repository, publisher Publisher) error

// GeneratorAPIURL returns the url in the api correct ecosystem.
type GeneratorAPIURL func(url url.URL) ([]url.URL, error)

var clientAPIs map[string]ClientAPI

// RegisterClientAPIs register all the client APIs for all the clients.
func RegisterClientAPIs() {

	clientAPIs = make(map[string]ClientAPI)

	clientAPIs["bitbucket"] = ClientAPI{
		Organization: RegisterBitbucketAPI(),
		Single:       RegisterSingleBitbucketAPI(),
		APIURL:       GenerateBitbucketAPIURL(),
	}

	clientAPIs["github"] = ClientAPI{
		Organization: RegisterGithubAPI(),
		Single:       RegisterSingleGithubAPI(),
		APIURL:       GenerateGithubAPIURL(),
	}

	clientAPIs["gitlab"] = ClientAPI{
		Organization: RegisterGitlabAPI(),
		Single:       RegisterSingleGitlabAPI(),
		APIURL:       GenerateGitlabAPIURL(),
	}

}

// GetClientAPICrawler checks if the API client for the requested organization clientAPI exists and return its handler.
func GetClientAPICrawler(clientAPI string) (OrganizationHandler, error) {
	if clientAPIs[clientAPI].Organization != nil {
		return clientAPIs[clientAPI].Organization, nil
	}

	return nil, fmt.Errorf("no organization client found for %s", clientAPI)

}

// GetSingleClientAPICrawler checks if the API client for the requested single repository clientAPI exists and return its handler.
func GetSingleClientAPICrawler(clientAPI string) (SingleRepoHandler, error) {
	if clientAPIs[clientAPI].Single != nil {
		return clientAPIs[clientAPI].Single, nil
	}
	return nil, fmt.Errorf("no single client found for %s", clientAPI)
}

// GetAPIURL checks if the API client for the requested API url exists and return its handler.
func GetAPIURL(clientAPI string) (GeneratorAPIURL, error) {
	if clientAPIs[clientAPI].APIURL != nil {
		return clientAPIs[clientAPI].APIURL, nil
	}
	return nil, fmt.Errorf("no api url generator client found for %s", clientAPI)
}

// GetClients returns a list of all registered clientAPI.
func GetClients() map[string]ClientAPI {
	return clientAPIs
}
