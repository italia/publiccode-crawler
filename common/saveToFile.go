package common

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/italia/developers-italia-backend/metrics"
	"github.com/spf13/viper"
)

// SaveToFile save the chosen <file_name> in DATADIR/repos/<source>/<vendor>/<repo>/<crawler_timestamp>_<file_name>.
func SaveToFile(domain Domain, hostname string, name string, data []byte, index string) error {
	if domain.Host == "" {
		return errors.New("cannot save a file without domain host")
	}
	if name == "" {
		return errors.New("cannot save a file without name")
	}

	fileName := index + "_publiccode.yml"
	vendor, repo := SplitFullName(name)

	path := filepath.Join(viper.GetString("CRAWLER_DATADIR"), hostname, vendor, repo)

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

// SplitFullName split a git FullName format to vendor and repo strings.
func SplitFullName(fullName string) (string, string) {
	s := strings.Split(fullName, "/")
	return s[0], s[1]
}
