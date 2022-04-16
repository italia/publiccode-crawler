package crawler

import (
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	url "github.com/italia/developers-italia-backend/crawler/internal"
)

var fileReaderInject = ioutil.ReadFile

// Whitelist contain a list of Public Administrations.
type Whitelist []Publisher

type Publisher struct {
	Id            string    `yaml:"id"`
	Name          string    `yaml:"name"`
	Organizations []url.URL `yaml:"orgs"`
	Repositories  []url.URL `yaml:"repos"`
}

// ReadAndParseWhitelist read the whitelist and return the parsed content in a slice of PA.
func ReadAndParseWhitelist(whitelistFile string) ([]Publisher, error) {
	// Open and read whitelist file.
	data, err := fileReaderInject(whitelistFile)
	if err != nil {
		return nil, fmt.Errorf("error in reading %s file: %v", whitelistFile, err)
	}

	// Parse whitelist file.
	whitelist, err := parseWhitelistFile(data)
	if err != nil {
		return nil, fmt.Errorf("error in parsing %s file: %v", whitelistFile, err)
	}
	log.Infof("Loaded and parsed %s", whitelistFile)

	return whitelist, err
}

// parseWhitelistFile parses the whitelist file to build a slice of PA.
func parseWhitelistFile(data []byte) ([]Publisher, error) {
	var whitelist []Publisher

	// Unmarshal the yml in domains list.
	err := yaml.Unmarshal(data, &whitelist)
	if err != nil {
		return nil, err
	}

	return whitelist, err
}
