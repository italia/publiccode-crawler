package crawler

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Blacklist contain a list of blocked repositories.
type Blacklist struct {
	Repos []Repo `yaml:"repos"`
}

// Repo matches a single repository.
type Repo struct {
	URL         string `yaml:"url"`
	Reason      string `yaml:"reason"`
	Description string `yaml:"description"`
}

// ReadAndParseBlacklist read the blacklist and return the parsed content in a slice of PA.
func ReadAndParseBlacklist(blacklistFile string) ([]Repo, error) {
	// Open and read blacklist file.
	// fileReaderInject is a var in the package, whitelist.go
	data, err := fileReaderInject(blacklistFile)
	if err != nil {
		return nil, fmt.Errorf("error in reading %s file: %v", blacklistFile, err)
	}
	// Parse blacklist file.
	blacklist, err := parseBlacklistFile(data)
	if err != nil {
		return nil, fmt.Errorf("error in parsing %s file: %v", blacklistFile, err)
	}
	log.Infof("Loaded and parsed %s", blacklistFile)

	return blacklist.Repos, err
}

// parseBlacklistFile parses the blacklist file to build a slice of Repo.
func parseBlacklistFile(data []byte) (Blacklist, error) {
	var blacklist Blacklist

	// Unmarshal the yml in domains list.
	err := yaml.Unmarshal(data, &blacklist)
	if err != nil {
		return Blacklist{}, err
	}

	return blacklist, err
}
