package crawler

import (
	"gopkg.in/yaml.v2"

	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

// Whitelist contain a list of Public Administrations.
type Whitelist []PA

// PA is a Public Administration.
type PA struct {
	ID            string   `yaml:"id"`
	CodiceIPA     string   `yaml:"codice-iPA"`
	Organizations []string `yaml:"organizations"`
}

// ReadAndParseWhitelist read the whitelist and return the parsed content in a slice of PA.
func ReadAndParseWhitelist(whitelistFile string) ([]PA, error) {
	// Open and read whitelist file.
	data, err := ioutil.ReadFile(whitelistFile)
	if err != nil {
		return nil, fmt.Errorf("error in reading %s file: %v", whitelistFile, err)
	}
	// Parse whitelist file.
	whitelist, err := parseWhitelistFile(data)
	if err != nil {
		return nil, fmt.Errorf("error in parsing %s file: %v", whitelistFile, err)
	}
	log.Infof("Loaded and parsed %s", whitelistFile)

	return whitelist, nil
}

// parseWhitelistFile parses the whitelist file to build a slice of PA.
func parseWhitelistFile(data []byte) ([]PA, error) {
	var whitelist []PA

	// Unmarshal the yml in domains list.
	err := yaml.Unmarshal(data, &whitelist)
	if err != nil {
		return nil, err
	}

	return whitelist, nil
}
