package crawler

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/italia/developers-italia-backend/metrics"
	"github.com/spf13/viper"
)

// SaveToFile save the chosen <file_name> in ./data/<source>/<vendor>/<repo>/<crawler_timestamp>_<file_name>.
func SaveToFile(domain Domain, hostname string, name string, data []byte, index string) error {
	if domain.Host == "" {
		return errors.New("cannot save a file without domain host")
	}
	if name == "" {
		return errors.New("cannot save a file without name")
	}

	fileName := index + "_" + viper.GetString("CRAWLED_FILENAME")
	vendor, repo := splitFullName(name)

	path := filepath.Join("./data", hostname, vendor, repo)

	// MkdirAll will create all the folder path, if not exists.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	err := ioutil.WriteFile(filepath.Join(path, fileName), data, 0644)
	if err != nil {
		return err
	}

	metrics.GetCounter("repository_file_saved", index).Inc()
	return err
}

// splitFullName split a git FullName format to vendor and repo strings.
func splitFullName(fullName string) (string, string) {
	s := strings.Split(fullName, "/")
	return s[0], s[1]
}
