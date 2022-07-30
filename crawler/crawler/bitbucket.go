package crawler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"
	"time"

	httpclient "github.com/italia/httpclient-lib-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Bitbucket is the complete response for the Bitbucket all repositories list.
type Bitbucket struct {
	Pagelen int `json:"pagelen"`
	Values  []struct {
		Scm        string `json:"scm"`
		Website    string `json:"website"`
		HasWiki    bool   `json:"has_wiki"`
		Name       string `json:"name"`
		Links      Links  `json:"links"`
		ForkPolicy string `json:"fork_policy"`
		UUID       string `json:"uuid"`
		Language   string `json:"language"`
		CreatedOn  string `json:"created_on"`
		Mainbranch struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"mainbranch"`
		FullName  string `json:"full_name"`
		HasIssues bool   `json:"has_issues"`
		Owner     struct {
			Username    string `json:"username"`
			DisplayName string `json:"display_name"`
			Type        string `json:"type"`
			UUID        string `json:"uuid"`
			Links       struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Avatar struct {
					Href string `json:"href"`
				} `json:"avatar"`
			} `json:"links"`
		} `json:"owner"`
		UpdatedOn   string `json:"updated_on"`
		Size        int    `json:"size"`
		Type        string `json:"type"`
		Slug        string `json:"slug"`
		IsPrivate   bool   `json:"is_private"`
		Description string `json:"description"`
		Project     struct {
			Key   string `json:"key"`
			Type  string `json:"type"`
			UUID  string `json:"uuid"`
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Avatar struct {
					Href string `json:"href"`
				} `json:"avatar"`
			} `json:"links"`
			Name string `json:"name"`
		} `json:"project,omitempty"`
		Parent struct {
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Avatar struct {
					Href string `json:"href"`
				} `json:"avatar"`
			} `json:"links"`
			Type     string `json:"type"`
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			UUID     string `json:"uuid"`
		} `json:"parent,omitempty"`
	} `json:"values"`
	Next string `json:"next"`
}

// BitbucketRepo is the complete response for the Bitbucket single repository.
type BitbucketRepo struct {
	Scm        string    `json:"scm"`
	Website    string    `json:"website"`
	HasWiki    bool      `json:"has_wiki"`
	Name       string    `json:"name"`
	Links      Links     `json:"links"`
	ForkPolicy string    `json:"fork_policy"`
	UUID       string    `json:"uuid"`
	Language   string    `json:"language"`
	CreatedOn  time.Time `json:"created_on"`
	Mainbranch struct {
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"mainbranch"`
	FullName  string `json:"full_name"`
	HasIssues bool   `json:"has_issues"`
	Owner     struct {
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		Type        string `json:"type"`
		UUID        string `json:"uuid"`
		Links       struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"links"`
	} `json:"owner"`
	UpdatedOn   time.Time `json:"updated_on"`
	Size        int       `json:"size"`
	Type        string    `json:"type"`
	Slug        string    `json:"slug"`
	IsPrivate   bool      `json:"is_private"`
	Description string    `json:"description"`
}

// Links is the list of Links associated to the repository.
type Links struct {
	Watchers struct {
		Href string `json:"href"`
	} `json:"watchers"`
	Branches struct {
		Href string `json:"href"`
	} `json:"branches"`
	Tags struct {
		Href string `json:"href"`
	} `json:"tags"`
	Commits struct {
		Href string `json:"href"`
	} `json:"commits"`
	Clone []struct {
		Href string `json:"href"`
		Name string `json:"name"`
	} `json:"clone"`
	Self struct {
		Href string `json:"href"`
	} `json:"self"`
	Source struct {
		Href string `json:"href"`
	} `json:"source"`
	HTML struct {
		Href string `json:"href"`
	} `json:"html"`
	Avatar struct {
		Href string `json:"href"`
	} `json:"avatar"`
	Hooks struct {
		Href string `json:"href"`
	} `json:"hooks"`
	Forks struct {
		Href string `json:"href"`
	} `json:"forks"`
	Downloads struct {
		Href string `json:"href"`
	} `json:"downloads"`
	Pullrequests struct {
		Href string `json:"href"`
	} `json:"pullrequests"`
}

// RegisterBitbucketAPI register the crawler function for Bitbucket API.
func RegisterBitbucketAPI() OrganizationHandler {
	return func(domain Domain, link string, repositories chan Repository, pa PA) (string, error) {
		// Set BasicAuth header.
		headers := make(map[string]string)
		if domain.BasicAuth != nil {
			n, err := generateRandomInt(len(domain.BasicAuth))
			if err != nil {
				return link, err
			}
			headers["Authorization"] = domain.BasicAuth[n]
		}

		// Parse url.
		u, err := url.Parse(link)
		if err != nil {
			return link, err
		}
		// Set domain host to new host.
		domain.Host = u.Hostname()

		// Get List of repositories.
		resp, err := httpclient.GetURL(link, headers)
		if err != nil {
			return link, err
		}
		if resp.Status.Code != http.StatusOK {
			log.Warnf("Request returned: %s", string(resp.Body))
			return "", errors.New("request returned an incorrect http.Status: " + resp.Status.Text)
		}

		// Fill response as list of values (repositories data).
		var result Bitbucket
		err = json.Unmarshal(resp.Body, &result)
		if err != nil {
			return link, err
		}

		// Add repositories to the channel that will perform the check on everyone.
		for _, v := range result.Values {

			// Join file raw URL.
			u, err := url.Parse(v.Links.HTML.Href)
			if err != nil {
				return link, err
			}
			u.Path = path.Join(u.Path, "raw", v.Mainbranch.Name, viper.GetString("CRAWLED_FILENAME"))

			// Marshal all the repository metadata.
			metadata, err := json.Marshal(v)
			if err != nil {
				log.Errorf("bitbucket metadata: %v", err)
			}

			// If the repository was never used, the Mainbranch is empty ("").
			if v.Mainbranch.Name != "" {
				repositories <- Repository{
					Name:        v.FullName,
					Hostname:    u.Hostname(),
					FileRawURL:  u.String(),
					GitCloneURL: v.Links.Clone[0].Href,
					GitBranch:   v.Mainbranch.Name,
					Domain:      domain,
					Pa:          pa,
					Headers:     headers,
					Metadata:    metadata,
				}
			}
		}

		// if last page for this organization, the result.Next is empty.
		if len(result.Next) == 0 {
			return "", nil
		}

		// Return next url.
		return result.Next, nil
	}
}

// RegisterSingleBitbucketAPI register the crawler function for single Bitbucket repository.
func RegisterSingleBitbucketAPI() SingleRepoHandler {
	return func(domain Domain, link string, repositories chan Repository, pa PA) error {
		// Set BasicAuth header
		headers := make(map[string]string)
		if domain.BasicAuth != nil {
			n, err := generateRandomInt(len(domain.BasicAuth))
			if err != nil {
				return err
			}
			headers["Authorization"] = domain.BasicAuth[n]
		}

		// Parse url.
		u, err := url.Parse(link)
		if err != nil {
			return err
		}

		// Set domain host to new host.
		domain.Host = u.Hostname()

		u.Path = path.Join("/2.0/repositories", u.Path)
		u.Host = "api." + u.Host

		linkRepo := u.String()

		// Get single Repo
		resp, err := httpclient.GetURL(linkRepo, headers)
		if err != nil {
			return err
		}
		if resp.Status.Code != http.StatusOK {
			log.Warnf("Request returned: %s", string(resp.Body))
			return errors.New("request returned an incorrect http.Status: " + resp.Status.Text)
		}

		// Fill response as list of values (repositories data).
		var result BitbucketRepo
		err = json.Unmarshal(resp.Body, &result)
		if err != nil {
			return err
		}

		// Join file raw URL.
		u, err = url.Parse(link)
		if err != nil {
			return err
		}
		fullURL := path.Join(u.Hostname(), result.FullName, "raw", result.Mainbranch.Name, viper.GetString("CRAWLED_FILENAME"))

		// Marshal all the repository metadata.
		metadata, err := json.Marshal(result)
		if err != nil {
			log.Errorf("bitbucket metadata: %v", err)
		}
		// If the repository was never used, the Mainbranch is empty ("").
		if result.Mainbranch.Name != "" {
			repositories <- Repository{
				Name:       result.FullName,
				Hostname:   u.Hostname(),
				FileRawURL: "https://" + fullURL,
				GitBranch:  result.Mainbranch.Name,
				Domain:     domain,
				Pa:         pa,
				Headers:    headers,
				Metadata:   metadata,
			}
		} else {
			return errors.New("repository is: empty")
		}

		return nil
	}
}

// GenerateBitbucketAPIURL returns the api url of given Bitbucket  organization link.
// IN: https://bitbucket.org/Soft
// OUT:https://api.bitbucket.org/2.0/repositories/Soft?pagelen=100
func GenerateBitbucketAPIURL() GeneratorAPIURL {
	return func(in string) (out []string, err error) {
		u, err := url.Parse(in)
		if err != nil {
			return []string{in}, err
		}
		u.Path = path.Join("/2.0/repositories", u.Path)
		u.Host = "api." + u.Host

		out = append(out, u.String())
		return
	}
}
