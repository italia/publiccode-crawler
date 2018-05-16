package persistency

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// SaveToFile save the chosen <file_name> in ./data/<source>/<vendor>/<repo>/<crawler_timestamp>_<file_name>.
func SaveToFile(source, name string, data []byte, fileTimestamp int64) {
	fileName := strconv.FormatInt(fileTimestamp, 10) + "_" + os.Getenv("CRAWLED_FILENAME")
	vendor, repo := splitFullName(name)

	path := filepath.Join("./data", source, vendor, repo)

	// MkdirAll will create all the folder path, if not exists.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	err := ioutil.WriteFile(filepath.Join(path, fileName), data, 0644)
	if err != nil {
		log.Error(err)
	}
}

// splitFullName split a git FullName format to vendor and repo strings.
func splitFullName(fullName string) (string, string) {
	s := strings.Split(fullName, "/")
	return s[0], s[1]
}
