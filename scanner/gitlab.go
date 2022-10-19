package scanner

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/xanzy/go-gitlab"

	"github.com/italia/publiccode-crawler/v3/common"
)

type GitLabScanner struct {
}

func NewGitLabScanner() Scanner {
	return GitLabScanner{}
}

// RegisterGitlabAPI register the crawler function for Gitlab API.
func (scanner GitLabScanner) ScanGroupOfRepos(url url.URL, publisher common.Publisher, repositories chan common.Repository) error {
	apiURL, _ := url.Parse("/api/v4")
	git, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"), gitlab.WithBaseURL(apiURL.String()))
	if err != nil {
		return err
	}

	if isGitlabGroup(url) {
		groupName := strings.Trim(url.Path, "/")

		group, _, err := git.Groups.GetGroup(groupName, &gitlab.GetGroupOptions{})
		if err != nil {
			return fmt.Errorf("can't get GitLag group '%s': %w", groupName, err)
		}

		if err = addGroupProjects(*group, publisher, repositories, git); err != nil {
			return err
		}
	} else {
		opts := &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{Page: 1},
		}

		for {
			projects, res, err := git.Projects.ListProjects(opts)
			if err != nil {
				return err
			}
			for _, prj := range projects {
				if err = addProject(nil, *prj, publisher, repositories); err != nil {
					return err
				}
			}

			if res.NextPage == 0 {
				break
			}
			opts.Page = res.NextPage
		}
	}

	return nil
}

// RegisterSingleGitlabAPI register the crawler function for single Bitbucket API.
func (scanner GitLabScanner) ScanRepo(url url.URL, publisher common.Publisher, repositories chan common.Repository) error {
	apiURL, _ := url.Parse("/api/v4")
	git, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"), gitlab.WithBaseURL(apiURL.String()))
	if err != nil {
		return err
	}

	projectName := strings.Trim(url.Path, "/")
	prj, _, err := git.Projects.GetProject(projectName, &gitlab.GetProjectOptions{})
	if err != nil {
		return err
	}

	if err = addProject(&url, *prj, publisher, repositories); err != nil {
		return err
	}

	return nil
}

// isGitlabGroup returns true if the API URL points to a group.
func isGitlabGroup(u url.URL) bool {
	return (
		// Always assume it's a group if the projects are hosted on gitlab.com,
		// because we only want to support groups (ie. not repos belonging to a user)
		strings.ToLower(u.Hostname()) == "gitlab.com" ||
		// Assume an on-premise GitLab's URL is a group if the path is not the root
		// path (/) or empty
		len(u.Path) > 1)
}

// generateGitlabRawURL returns the file Gitlab specific file raw url.
func generateGitlabRawURL(baseURL, defaultBranch string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, "raw", defaultBranch, "publiccode.yml")

	return u.String(), err
}

// addGroupProjects sends all the projects in a GitLab group, including all subgroups, to
// the repositories channel
func addGroupProjects(group gitlab.Group, publisher common.Publisher, repositories chan common.Repository, client *gitlab.Client) error {
	opts := &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{Page: 1},
	}

	for {
		projects, res, err := client.Groups.ListGroupProjects(group.ID, opts)
		if err != nil {
			return err
		}
		for _, prj := range projects {
			err = addProject(nil, *prj, publisher, repositories)
			if err != nil {
				return err
			}
		}

		if res.NextPage == 0 {
			break
		}
		opts.Page = res.NextPage
	}

	dgOpts := &gitlab.ListDescendantGroupsOptions{
		ListOptions: gitlab.ListOptions{Page: 1},
	}
	for {
		groups, res, err := client.Groups.ListDescendantGroups(group.ID, dgOpts)
		if err != nil {
			return err
		}
		for _, g := range groups {
			err = addGroupProjects(*g, publisher, repositories, client)
			if err != nil {
				return err
			}
		}

		if res.NextPage == 0 {
			break
		}
		dgOpts.Page = res.NextPage
	}

	return nil
}

// addGroupProjects sends the GitLab project the repositories channel
func addProject(originalURL *url.URL, project gitlab.Project, publisher common.Publisher, repositories chan common.Repository) error {
	// Join file raw URL string.
	rawURL, err := generateGitlabRawURL(project.WebURL, project.DefaultBranch)
	if err != nil {
		return err
	}


	if project.DefaultBranch != "" {
		canonicalURL, err := url.Parse(project.HTTPURLToRepo)
		if err != nil {
			return fmt.Errorf("failed to get canonical repo URL for %s: %w", project.WebURL, err)
		}
		if originalURL == nil {
			originalURL = canonicalURL
		}

		repositories <- common.Repository{
			Name:         project.PathWithNamespace,
			FileRawURL:   rawURL,
			URL:          *originalURL,
			CanonicalURL: *canonicalURL,
			GitBranch:    project.DefaultBranch,
			Publisher:    publisher,
		}
	}

	return nil
}
