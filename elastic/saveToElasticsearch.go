package elastic

import (
	"context"
	"errors"
	"net/url"
	"path"
	"time"

	"github.com/alranel/go-vcsurl/v2"
	"github.com/ghodss/yaml"
	"github.com/italia/developers-italia-backend/metrics"
	publiccode "github.com/italia/publiccode-parser-go/v3"
	elastic "github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/italia/developers-italia-backend/common"
)

type administration struct {
	Name      string `json:"it-riuso-codiceIPA-label"`
	CodiceIPA string `json:"it-riuso-codiceIPA"`
	Type      string `json:"type"`
}

func SaveToES(client *elastic.Client, index string, repo common.Repository, activityIndex float64, vitality []int, parser publiccode.Parser) error {
	// softwareES represents a software record in Elasticsearch
	type softwareES struct {
		FileRawURL            string            `json:"fileRawURL"`
		ID                    string            `json:"id"`
		CrawlTime             string            `json:"crawltime"`
		ItRiusoCodiceIPALabel string            `json:"it-riuso-codiceIPA-label"`
		Slug                  string            `json:"slug"`
		PublicCode            interface{}       `json:"publiccode"`
		VitalityScore         float64           `json:"vitalityScore"`
		VitalityDataChart     []int             `json:"vitalityDataChart"`
		OEmbedHTML            map[string]string `json:"oEmbedHTML"`
		Type                  string            `json:"type"`
	}

	// TODO: We should probably get rid of this and maintain the original
	// publiccode.yml in the database, and expand the logo and screenshots paths
	// client side.
	publiccode := &parser.PublicCode
	rawRoot, err := vcsurl.GetRawRoot((*url.URL)(parser.PublicCode.URL), parser.Branch)
	if err != nil {
		return err
	}

	if publiccode.Logo != "" {
		logoURL, _ := url.Parse(publiccode.Logo)
		if !logoURL.IsAbs() {
			*logoURL = *rawRoot
			logoURL.Path = path.Join(rawRoot.Path, publiccode.Logo)
		}
		publiccode.Logo = logoURL.String()
	}
	for lang, desc := range publiccode.Description {
		for idx, screenshot := range desc.Screenshots {
			screenshotURL, _ := url.Parse(screenshot)
			if !screenshotURL.IsAbs() {
				*screenshotURL = *rawRoot
				screenshotURL.Path = path.Join(rawRoot.Path, screenshot)
			}

			publiccode.Description[lang].Screenshots[idx] = screenshotURL.String()
		}
	}

	yml, err := parser.ToYAML()
	if err != nil {
		return err
	}

	// Create a softwareES object and populate it
	file := softwareES{
		FileRawURL:            repo.FileRawURL,
		ID:                    repo.GenerateID(),
		CrawlTime:             time.Now().Format(time.RFC3339),
		Slug:                  repo.GenerateSlug(),
		ItRiusoCodiceIPALabel: GetAdministrationName(publiccode.It.Riuso.CodiceIPA),
		VitalityScore:         activityIndex,
		VitalityDataChart:     vitality,
		Type:                  "software",
	}

	err = yaml.Unmarshal(yml, &file.PublicCode)
	if err != nil {
		return err
	}

	// Put publiccode data in ES.
	ctx := context.Background()
	_, err = client.Index().
		Index(index).
		Id(file.ID).
		BodyJson(file).
		Do(ctx)
	if err != nil {
		return err
	}

	metrics.GetCounter("repository_file_indexed", index).Inc()

	// Add administration data.
	if publiccode.It.Riuso.CodiceIPA != "" {
		// Put administrations data in ES.
		_, err = client.Index().
			Index(viper.GetString("ELASTIC_PUBLISHERS_INDEX")).
			Id(publiccode.It.Riuso.CodiceIPA).
			BodyJson(administration{
				Name:      file.ItRiusoCodiceIPALabel,
				CodiceIPA: publiccode.It.Riuso.CodiceIPA,
				Type:      "administration",
			}).
			Do(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteByQueryFromES delete record from elasticsearch
// that will match search string for publiccode.url field
func DeleteByQueryFromES(client *elastic.Client, search string, index string) error {
	// Search with a term query
	termQuery := elastic.NewTermQuery("publiccode.url", search)

	// Put publiccode data in ES.
	ctx := context.Background()
	searchResult, err := client.DeleteByQuery().
		Index(index).
		Query(termQuery). // specify the query
		Do(ctx)           // execute
	if err != nil {
		return err
	}

	if searchResult == nil {
		return errors.New("Generic error on DeleteByQueryFromES()")
	}

	if searchResult.Deleted == 0 {
		return errors.New("No records deleted for searched query")
	}

	log.Infof("Deleted %d record from ES linked to %s", searchResult.Deleted, search)
	return nil
}
