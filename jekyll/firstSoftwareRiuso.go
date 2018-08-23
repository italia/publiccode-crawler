package jekyll

import (
	"context"
	"os"
	"reflect"
	"strings"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// SoftwareRiuso is a simple description of a Software with  it/riuso/codiceIPA key.
type SoftwareRiuso struct {
	Name      string `json:"name"`
	Logo      string `json:"logo"`
	URL       string `json:"url"`
	CodiceIPA string `json:"ipa"`
}

// FirstSoftwareRiuso generate a yml file with simplified info about SoftwareRiuso, ordered by releaseDate.
func FirstSoftwareRiuso(filename string, results int, elasticClient *elastic.Client) error {
	log.Infof("Generating %s", filename)

	// Create file if not exists.
	if _, err := os.Stat(filename); os.IsExist(err) {
		err := os.Remove(filename)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	// Open file.
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close() // nolint: errcheck

	// Administrations data.
	var softwareRiuso []SoftwareRiuso

	// Generate query.
	query := elastic.NewExistsQuery("it-riuso-codice-ipa")

	// Extract all the documents.
	searchResult, err := elasticClient.Search().
		Index("publiccode").        // search in index "publiccode"
		Query(query).               // specify the query
		Sort("releaseDate", false). // sort by releaseDate, from newest to oldest.
		Pretty(true).               // pretty print request and response JSON
		From(0).Size(results).      // get first 10k elements. It can be changed.
		Do(context.Background())    // execute
	if err != nil {
		log.Error(err)
	}

	// Foreach search result check if codiceIPA is not empty.
	var pctype crawler.PublicCodeES
	for _, item := range searchResult.Each(reflect.TypeOf(pctype)) {
		i := item.(crawler.PublicCodeES)

		rawBaseDir := strings.TrimRight(i.FileRawURL, viper.GetString("CRAWLED_FILENAME"))

		if i.ItRiusoCodiceIPA != "" {
			softwareRiuso = append(softwareRiuso, SoftwareRiuso{
				Name:      i.Name,
				Logo:      concatenateLink(rawBaseDir, i.Logo),
				URL:       i.URL,
				CodiceIPA: i.ItRiusoCodiceIPA,
			})

		}
	}
	// Debug note if file will be empty.
	if len(softwareRiuso) == 0 {
		log.Warnf("%s is empty.", filename)
	}

	// Marshal yml.
	d, err := yaml.Marshal(&softwareRiuso)
	if err != nil {
		return err
	}
	//Append data to file.
	if _, err = f.WriteString(string(d)); err != nil {
		return err
	}

	return err
}
