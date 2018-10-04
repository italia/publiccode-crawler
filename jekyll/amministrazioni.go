package jekyll

import (
	"context"
	"os"
	"reflect"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/italia/developers-italia-backend/ipa"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// Administration is a simple description of an Administration.
type Administration struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	CodiceIPA string `json:"ipa"`
}

// AmministrazioniYML generate a yml file with all the amministrazioni in es.
func AmministrazioniYML(filename string, unsupportedCountries []string, elasticClient *elastic.Client) error {
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
	var administrations []Administration

	// UnsupportedCountries.
	uc := make([]interface{}, len(unsupportedCountries))
	for i, v := range unsupportedCountries {
		uc[i] = v
	}

	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTypeQuery("software"))
	query = query.MustNot(elastic.NewTermsQuery("intended-audience-unsupported-countries", uc...))

	searchResult, err := elasticClient.Search().
		Index("publiccode").     // search in index "publiccode"
		Query(query).            // specify the query
		Pretty(true).            // pretty print request and response JSON
		From(0).Size(10000).     // get first 10k elements. It can be changed.
		Do(context.Background()) // execute
	if err != nil {
		log.Error(err)
	}

	// Foreach search result check if codiceIPA is not empty.
	var pctype crawler.PublicCodeES
	for _, item := range searchResult.Each(reflect.TypeOf(pctype)) {
		i := item.(crawler.PublicCodeES)

		if i.ItRiusoCodiceIPA != "" {
			administrations = append(administrations, Administration{
				Name:      ipa.GetAdministrationName(i.ItRiusoCodiceIPA),
				URL:       i.LandingURL,
				CodiceIPA: i.ItRiusoCodiceIPA,
			})
		}

	}
	// Debug note if file will be empty.
	if len(administrations) == 0 {
		log.Warnf("%s is empty.", filename)
	}

	// Remove duplicates.
	administrations = removeDuplicates(administrations)

	// Marshal yml.
	d, err := yaml.Marshal(&administrations)
	if err != nil {
		return err
	}
	//Append data to file.
	if _, err = f.WriteString(string(d)); err != nil {
		return err
	}

	return err
}

func removeDuplicates(elements []Administration) []Administration {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []Administration{}

	for v := range elements {
		if encountered[elements[v].CodiceIPA] {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v].CodiceIPA] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}
