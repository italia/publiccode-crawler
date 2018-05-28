package crawler

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/italia/developers-italia-backend/metrics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// SaveToFile save the chosen <file_name> in ./data/<source>/<vendor>/<repo>/<crawler_timestamp>_<file_name>.
func SaveToFile(domain Domain, name string, data []byte, index string) {
	fileName := index + "_" + viper.GetString("CRAWLED_FILENAME")
	vendor, repo := splitFullName(name)

	path := filepath.Join("./data", domain.ID, vendor, repo)

	// MkdirAll will create all the folder path, if not exists.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	err := ioutil.WriteFile(filepath.Join(path, fileName), data, 0644)
	if err != nil {
		log.Error(err)
	}

	metrics.GetCounter("repository_file_saved", index).Inc()
}

// splitFullName split a git FullName format to vendor and repo strings.
func splitFullName(fullName string) (string, string) {
	s := strings.Split(fullName, "/")
	return s[0], s[1]
}
