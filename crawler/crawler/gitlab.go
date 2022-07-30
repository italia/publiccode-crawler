package crawler

import (
	"encoding/json"
	"math/rand"
	"net/url"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

// isGitlabGroup returns true if the API URL points to a group.
func isGitlabGroup(u url.URL) bool {
	return strings.ToLower(u.Hostname()) == "gitlab.com" ||
		// When u.Path is /api/v4/groups there's no group, otherwise
		// it would have been /api/v4/groups/$GROUPNAME.
		u.Path != "/api/v4/groups"
}

// RegisterGitlabAPI register the crawler function for Gitlab API.
func RegisterGitlabAPI() OrganizationHandler {
	return func(domain Domain, link string, repositories chan Repository, pa PA) (string, error) {
		log.Debugf("RegisterGitlabAPI: %s ", link)

		headers := make(map[string]string)

		token := ""
		if domain.BasicAuth != nil {
			token = domain.BasicAuth[rand.Intn(len(domain.BasicAuth))]
			// TODO: refactor in order to not need to pass Headers around
			headers["Authorization"] = token
		}

		u, err := url.Parse(link)
		if err != nil {
			return link, err
		}
		// Set domain host to new host.
		domain.Host = u.Hostname()

		apiURL, _ := u.Parse("/api/v4")
		git, err := gitlab.NewClient(token, gitlab.WithBaseURL(apiURL.String()))
		if err != nil {
			return link, err
		}

		if isGitlabGroup(*u) {
			groupName := strings.Replace(u.Path, "/api/v4/groups/", "", -1)

			group, _, err := git.Groups.GetGroup(groupName, &gitlab.GetGroupOptions{})
			if err != nil {
				return link, err
			}

			err = addGroupProjects(*group, domain, pa, headers, repositories, git)
			if err != nil {
				return link, err
			}
		} else {
			opts := &gitlab.ListProjectsOptions{
				ListOptions: gitlab.ListOptions{Page: 1},
			}

			for {
				projects, res, err := git.Projects.ListProjects(opts)
				if err != nil {
					return link, err
				}
				for _, prj := range projects {
					err = addProject(*prj, domain, pa, headers, repositories)
					if err != nil {
						return link, err
					}
				}

				if res.NextPage == 0 {
					break
				}
				opts.Page = res.NextPage
			}
		}

		return "", nil
	}
}

// RegisterSingleGitlabAPI register the crawler function for single Bitbucket API.
func RegisterSingleGitlabAPI() SingleRepoHandler {
	return func(domain Domain, link string, repositories chan Repository, pa PA) error {
		headers := make(map[string]string)

		token := ""
		if domain.BasicAuth != nil {
			token = domain.BasicAuth[rand.Intn(len(domain.BasicAuth))]
			// TODO: refactor in order to not need to pass Headers around
			headers["Authorization"] = token
		}

		u, err := url.Parse(link)
		if err != nil {
			return err
		}
		// Set domain host to new host.
		domain.Host = u.Hostname()

		apiURL, _ := u.Parse("/api/v4")
		git, err := gitlab.NewClient(token, gitlab.WithBaseURL(apiURL.String()))
		if err != nil {
			return err
		}

		projectName := strings.Trim(u.Path, "/")
		prj, _, err := git.Projects.GetProject(projectName, &gitlab.GetProjectOptions{})
		if err != nil {
			return err
		}

		err = addProject(*prj, domain, pa, headers, repositories)
		if err != nil {
			return err
		}

		return nil
	}
}

// generateGitlabRawURL returns the file Gitlab specific file raw url.
func generateGitlabRawURL(baseURL, defaultBranch string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, "raw", defaultBranch, viper.GetString("CRAWLED_FILENAME"))

	return u.String(), err
}

// addGroupProjects sends all the projects in a GitLab group, including all subgroups, to
// the repositories channel
func addGroupProjects(group gitlab.Group, domain Domain, pa PA, headers map[string]string, repositories chan Repository, client *gitlab.Client) error {
	opts := &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{Page: 1},
	}

	for {
		projects, res, err := client.Groups.ListGroupProjects(group.ID, opts)
		if err != nil {
			return err
		}
		for _, prj := range projects {
			err = addProject(*prj, domain, pa, headers, repositories)
			if err != nil {
				return err
			}
		}

		if res.NextPage == 0 {
			break
		}
		opts.Page = res.NextPage
	}

	dg_opts := &gitlab.ListDescendantGroupsOptions{
		ListOptions: gitlab.ListOptions{Page: 1},
	}
	for {
		groups, res, err := client.Groups.ListDescendantGroups(group.ID, dg_opts)
		if err != nil {
			return err
		}
		for _, g := range groups {
			err = addGroupProjects(*g, domain, pa, headers, repositories, client)
			if err != nil {
				return err
			}
		}

		if res.NextPage == 0 {
			break
		}
		dg_opts.Page = res.NextPage
	}

	return nil
}

// addGroupProjects sends the GitLab project the repositories channel
func addProject(project gitlab.Project, domain Domain, pa PA, headers map[string]string, repositories chan Repository) error {
	// Join file raw URL string.
	rawURL, err := generateGitlabRawURL(project.WebURL, project.DefaultBranch)
	if err != nil {
		return err
	}

	// Marshal all the repository metadata.
	metadata, err := json.Marshal(project)
	if err != nil {
		log.Errorf("gitlab metadata: %v", err)
		return err
	}

	if project.DefaultBranch != "" {
		repositories <- Repository{
			Name:        project.PathWithNamespace,
			Hostname:    domain.Host,
			FileRawURL:  rawURL,
			GitCloneURL: project.HTTPURLToRepo,
			GitBranch:   project.DefaultBranch,
			Domain:      domain,
			Pa:          pa,
			Headers:     headers,
			Metadata:    metadata,
		}
	}

	return nil
}

// GenerateGitlabAPIURL returns the api url of given Gitlab organization link.
// IN: https://gitlab.org/blockninja
// OUT:https://gitlab.com/api/v4/groups/blockninja
func GenerateGitlabAPIURL() GeneratorAPIURL {
	return func(in string) (out []string, err error) {
		u, err := url.Parse(in)
		if err != nil {
			return []string{in}, err
		}
		u.Path = path.Join("api/v4/groups", u.Path)

		out = append(out, u.String())
		return
	}
}
