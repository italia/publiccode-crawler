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

// SoftwareOpenSource is a simple description of a Software without it/riuso/codiceIPA key.
type SoftwareOpenSource struct {
	Name      string `json:"name"`
	Logo      string `json:"logo"`
	URL       string `json:"url"`
	CodiceIPA string `json:"ipa"`
}

// FirstSoftwareOpenSource generate a yml file with simplified info about SoftwareRiuso, ordered by releaseDate.
func FirstSoftwareOpenSource(filename string, results int, elasticClient *elastic.Client) error {
	log.Debugf("Generating %s", filename)

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
	var softwareOS []SoftwareOpenSource

	// Extract all the softwares.
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTypeQuery("software"))

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

		if i.ItRiusoCodiceIPA == "" {
			softwareOS = append(softwareOS, SoftwareOpenSource{
				Name:      i.Name,
				Logo:      concatenateLink(rawBaseDir, i.Logo),
				URL:       i.URL,
				CodiceIPA: i.ItRiusoCodiceIPA,
			})

		}
	}
	// Debug note if file will be empty.
	if len(softwareOS) == 0 {
		log.Warnf("%s is empty.", filename)
	}

	// Remove duplicates.
	softwareOS = removeDuplicatesSoftwareOS(softwareOS)

	// Marshal yml.
	d, err := yaml.Marshal(&softwareOS)
	if err != nil {
		return err
	}
	//Append data to file.
	if _, err = f.WriteString(string(d)); err != nil {
		return err
	}

	return err
}

func removeDuplicatesSoftwareOS(elements []SoftwareOpenSource) []SoftwareOpenSource {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []SoftwareOpenSource{}

	for v := range elements {
		if encountered[elements[v].Name] {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v].Name] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}
