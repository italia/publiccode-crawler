package common

import (
	"fmt"
	"net/url"
	"os"

	internalURL "github.com/italia/publiccode-crawler/v4/internal"
	"gopkg.in/yaml.v2"
)

var fileReaderInject = os.ReadFile

// CatalogSource is a single source of repository URLs.
// Driver identifies the driver ("github", "gitlab", "bitbucket", "gitea", "json", ...).
// Args holds positional arguments for the driver, e.g. a JSONPath expression for "json".
// Group distinguishes org/group scans (true) from single-repo scans (false).
type CatalogSource struct {
	URL    url.URL
	Driver string
	Args   []string
	Group  bool
}

type Publisher struct {
	ID            string
	AlternativeID string
	Name          string
	Sources       []CatalogSource
}

// publisherYAML is the on-disk representation. Driver is inferred from the URL.
type publisherYAML struct {
	ID            string            `yaml:"id"`
	Name          string            `yaml:"name"`
	Organizations []internalURL.URL `yaml:"orgs"`
	Repositories  []internalURL.URL `yaml:"repos"`
}

// LoadPublishers loads the publishers YAML file and returns a slice of Publisher.
func LoadPublishers(path string) ([]Publisher, error) {
	data, err := fileReaderInject(path)
	if err != nil {
		return nil, fmt.Errorf("error in reading `%s': %w", path, err)
	}

	var raw []publisherYAML

	err = yaml.Unmarshal(data, &raw)
	if err != nil {
		return nil, fmt.Errorf("error in parsing `%s': %w", path, err)
	}

	publishers := make([]Publisher, 0, len(raw))

	for _, rawPub := range raw {
		pub := Publisher{
			ID:   rawPub.ID,
			Name: rawPub.Name,
		}

		for _, org := range rawPub.Organizations {
			stdURL := (url.URL)(org)
			pub.Sources = append(pub.Sources, CatalogSource{
				URL:    stdURL,
				Driver: InferVCSDriver(stdURL),
				Group:  true,
			})
		}

		for _, repo := range rawPub.Repositories {
			stdURL := (url.URL)(repo)
			pub.Sources = append(pub.Sources, CatalogSource{
				URL:    stdURL,
				Driver: InferVCSDriver(stdURL),
				Group:  false,
			})
		}

		publishers = append(publishers, pub)
	}

	return publishers, nil
}
