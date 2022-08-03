package crawler

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/italia/developers-italia-backend/metrics"
	log "github.com/sirupsen/logrus"
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

	fileName := index + "_" + viper.GetString("CRAWLED_FILENAME")
	vendor, repo := splitFullName(name)

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

// splitFullName split a git FullName format to vendor and repo strings.
func splitFullName(fullName string) (string, string) {
	s := strings.Split(fullName, "/")
	return s[0], s[1]
}

// Save the bad publiccode.yaml url to a file used by the publiccode-issueopener script.
func logBadYamlToFile(fileRawURL string) {
	log.Errorf("Appending the bad file URL to the list: %s", fileRawURL)

	filePath := path.Join(viper.GetString("CRAWLER_DATADIR"), "bad_publiccodes.lst")

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Errorf(err.Error())
	}

	_, err = f.WriteString(time.Now().Format("2006-01-02T15:04:05") + " - " + fileRawURL + "\r\n")

	if err != nil {
		log.Errorf(err.Error())
	}

	f.Sync()        // nolint: errcheck
	defer f.Close() // nolint: errcheck
}
