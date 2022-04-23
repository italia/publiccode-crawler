package crawler

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

// GetAllBlackListedRepos return all blacklisted repositories
func GetAllBlackListedRepos() map[string]string {
	files := viper.GetString("BLACKLIST_FOLDER")
	pattern := viper.GetString("BLACKLIST_PATTERN")
	if files == "" || pattern == "" {
		log.Warn("BLACKLIST_* vars are not defined in config.toml, please define both")
		return nil
	}

	readBlacklist, err := scanBlacklists(files, pattern)
	if err != nil {
		log.Errorf("path not exists or you don't have permission: %s", err)
		return nil
	}
	var repoListed = make(map[string]string)
	for _, repo := range readBlacklist {
		repoListed[appendGitExt(repo.URL)] = repo.URL
	}
	return repoListed
}

// IsRepoInBlackList checks whether a repo is in blacklist
func IsRepoInBlackList(repoURL string) bool {
	files := viper.GetString("BLACKLIST_FOLDER")
	pattern := viper.GetString("BLACKLIST_PATTERN")
	if files == "" || pattern == "" {
		log.Warn("BLACKLIST_* var are not defined in config.toml, please define both")
		return false
	}

	readBlacklist, err := scanBlacklists(files, pattern)
	if err != nil {
		log.Errorf("path not exists or you don't have permission: %s", err)
		return false
	}
	for _, repo := range readBlacklist {
		if repo.URL == repoURL {
			log.Warnf("PA found in blacklist with reason: "+
				"%s and description: %s, skipping...", repo.Reason, repo.Description)
			return true
		}
	}
	return false
}

func appendGitExt(repo string) string {
	re := regexp.MustCompile(`\.git$`)
	if re.MatchString(repo) {
		return repo
	}
	return repo + ".git"
}

// ReadAndParseBlacklist read the blacklist and return the parsed content in a slice of PA.
func ReadAndParseBlacklist(blacklistFile string) ([]Repo, error) {
	// Open and read blacklist file.
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

func scanBlacklists(dir string, pattern string) ([]Repo, error) {
	files, err := WalkMatch(dir, pattern)
	if err != nil {
		return nil, err
	}
	var repos []Repo
	for _, file := range files {
		blacklistSlice, _ := ReadAndParseBlacklist(file)
		repos = append(repos, blacklistSlice...)
	}
	return repos, nil
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

// WalkMatch util func
func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
