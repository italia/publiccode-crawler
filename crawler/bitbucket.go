package crawler

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/italia/developers-italia-backend/httpclient"
	log "github.com/sirupsen/logrus"
	"sync"
)

// Bitbucket represent a complete result for the Bitbucket API respose from all repositories list.
type Bitbucket struct {
	Pagelen int `json:"pagelen"`
	Values  []struct {
		Scm     string `json:"scm"`
		Website string `json:"website"`
		HasWiki bool   `json:"has_wiki"`
		Name    string `json:"name"`
		Links   struct {
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
		} `json:"links"`
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

// BitbucketRepo represent a complete for the Bitbucket API respose from a single repository.
type BitbucketRepo struct {
	Scm     string `json:"scm"`
	Website string `json:"website"`
	HasWiki bool   `json:"has_wiki"`
	Name    string `json:"name"`
	Links   struct {
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
	} `json:"links"`
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

// RegisterBitbucketAPI register the crawler function for Bitbucket API.
func RegisterBitbucketAPI() Handler {
	return func(domain Domain, link string, repositories chan Repository, wg *sync.WaitGroup) (string, error) {
		// Set BasicAuth header
		headers := make(map[string]string)
		if domain.BasicAuth != nil {
			rand.Seed(time.Now().Unix())
			n := rand.Int() % len(domain.BasicAuth)
			headers["Authorization"] = "Basic " + domain.BasicAuth[n]
		}

		// Get List of repositories
		resp, err := httpclient.GetURL(link, headers)
		if err != nil {
			return link, err
		}
		if resp.Status.Code != http.StatusOK {
			log.Warnf("Request returned: %s", string(resp.Body))
			return link, errors.New("request returned an incorrect http.Status: " + resp.Status.Text)
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
			u.Path = path.Join(u.Path, "raw", v.Mainbranch.Name, os.Getenv("CRAWLED_FILENAME"))

			// If the repository was never used, the Mainbranch is empty ("")
			if v.Mainbranch.Name != "" {
				repositories <- Repository{
					Name:       v.FullName,
					FileRawURL: u.String(),
					Domain:     domain,
					Headers:    headers,
				}
			}
		}

		// Bitbucket end reached.
		if len(result.Next) == 0 {
			for len(repositories) != 0 {
				time.Sleep(time.Second)
			}
			log.Info("Bitbucket repositories status: end reached. Restart from domain value:" + domain.URL)

			// Restart.
			// return domain.URL, nil
			return "", nil

		}

		// Return next url
		return result.Next, nil
	}
}

// RegisterSingleBitbucketAPI register the crawler function for single Bitbucket API.
func RegisterSingleBitbucketAPI() SingleHandler {
	return func(domain Domain, link string, repositories chan Repository) error {
		// Set BasicAuth header
		headers := make(map[string]string)
		if domain.BasicAuth != nil {
			rand.Seed(time.Now().Unix())
			n := rand.Int() % len(domain.BasicAuth)
			headers["Authorization"] = "Basic " + domain.BasicAuth[n]
		}

		u, err := url.Parse(link)
		if err != nil {
			log.Error(err)
		}

		// Clear the url.
		fullName := u.Path
		if u.Path[:1] == "/" {
			fullName = fullName[1:]
		}
		if u.Path[len(u.Path)-1:] == "/" {
			fullName = fullName[:len(u.Path)-2]
		}

		fullURL := "https://api.bitbucket.org/2.0/repositories/" + fullName

		// Get single Repo
		resp, err := httpclient.GetURL(fullURL, headers)
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
		u, err = url.Parse(domain.RawBaseUrl)
		if err != nil {
			return err
		}
		u.Path = path.Join(u.Path, result.FullName, "raw", result.Mainbranch.Name, os.Getenv("CRAWLED_FILENAME"))

		// If the repository was never used, the Mainbranch is empty ("")
		if result.Mainbranch.Name != "" {
			repositories <- Repository{
				Name:       result.FullName,
				FileRawURL: u.String(),
				Domain:     domain,
				Headers:    headers,
			}
		} else {
			return errors.New("repository is: empty")
		}

		return nil
	}
}
