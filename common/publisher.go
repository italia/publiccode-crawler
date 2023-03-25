package common

import (
	"fmt"
	"io/ioutil"

	url "github.com/italia/publiccode-crawler/v3/internal"
	"gopkg.in/yaml.v2"
)

var fileReaderInject = ioutil.ReadFile

type Publisher struct {
	Id            string    `yaml:"id"`
	Name          string    `yaml:"name"`
	Organizations []url.URL `yaml:"orgs"`
	Repositories  []url.URL `yaml:"repos"`
}

// LoadPublishers loads the publishers YAML file and returns a slice of Publisher.
func LoadPublishers(path string) ([]Publisher, error) {
	data, err := fileReaderInject(path)
	if err != nil {
		return nil, fmt.Errorf("error in reading `%s': %w", path, err)
	}

	var publishers []Publisher
	err = yaml.Unmarshal(data, &publishers)
	if err != nil {
		return nil, fmt.Errorf("error in parsing `%s': %w", path, err)
	}

	return publishers, nil
}
