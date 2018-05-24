package crawler

import (
	"gopkg.in/yaml.v2"

	"errors"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

// PA is a Public Administration data.
type PA struct {
	CodiceIPA    string `yaml:"codiceIPA"`
	Name         string `yaml:"name"`
	Repositories []struct {
		API           string   `yaml:"api"`
		Organizations []string `yaml:"organizations"`
	} `yaml:"repositories"`
}

func ReadAndParseWhitelist(whitelistFile string) ([]PA, error) {
	// Open and read whitelist file.
	data, err := ioutil.ReadFile(whitelistFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in reading %s file: %v", whitelistFile, err))
	}
	// Parse whitelist file.
	whitelist, err := parseWhitelistFile(data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in parsing %s file: %v", whitelistFile, err))
	}
	log.Infof("Loaded and parsed %s", whitelistFile)

	return whitelist, nil
}

// parseWhitelistFile parses the whitelist file to build a slice of PA.
func parseWhitelistFile(data []byte) ([]PA, error) {
	whitelist := []PA{}

	// Unmarshal the yml in domains list.
	err := yaml.Unmarshal(data, &whitelist)
	if err != nil {
		return nil, err
	}

	return whitelist, nil
}
