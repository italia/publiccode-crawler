package crawler

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/italia/developers-italia-backend/crawler/ipa"
	"github.com/italia/developers-italia-backend/crawler/metrics"
	pcode "github.com/italia/publiccode-parser-go"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type administration struct {
	Name      string `json:"it-riuso-codiceIPA-label"`
	CodiceIPA string `json:"it-riuso-codiceIPA"`
}

// saveToES save the chosen data []byte in elasticsearch
// data contains the raw publiccode.yml file
func (c *Crawler) saveToES(repo Repository, activityIndex float64, vitality []int, data []byte) error {
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
	}

	// Parse the publiccode.yml file
	parser := pcode.NewParser()
	parser.Strict = false
	parser.RemoteBaseURL = strings.TrimRight(repo.FileRawURL, viper.GetString("CRAWLED_FILENAME"))
	err := parser.Parse(data)
	if err != nil {
		log.Errorf("Error parsing publiccode.yml: %v", err)
	}

	// Create a softwareES object and populate it
	file := softwareES{
		FileRawURL:            repo.FileRawURL,
		ID:                    repo.generateID(),
		CrawlTime:             time.Now().Format(time.RFC3339),
		Slug:                  repo.generateSlug(),
		ItRiusoCodiceIPALabel: ipa.GetAdministrationName(parser.PublicCode.It.Riuso.CodiceIPA),
		VitalityScore:         activityIndex,
		VitalityDataChart:     vitality,
		OEmbedHTML:            parser.OEmbed,
	}

	// Convert parser.PublicCode to YAML and parse it again into the softwareES record
	yml, err := parser.ToYAML()
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yml, &file.PublicCode)

	// Put publiccode data in ES.
	ctx := context.Background()
	_, err = c.es.Index().
		Index(c.index).
		Type("software").
		Id(file.ID).
		BodyJson(file).
		Do(ctx)
	if err != nil {
		return err
	}

	metrics.GetCounter("repository_file_indexed", c.index).Inc()

	// Add administration data.
	if parser.PublicCode.It.Riuso.CodiceIPA != "" {
		// Put administrations data in ES.
		_, err = c.es.Index().
			Index(viper.GetString("ELASTIC_PUBLISHERS_INDEX")).
			Type("administration").
			Id(parser.PublicCode.It.Riuso.CodiceIPA).
			BodyJson(administration{
				Name:      file.ItRiusoCodiceIPALabel,
				CodiceIPA: parser.PublicCode.It.Riuso.CodiceIPA,
			}).
			Do(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// generateID generates a hash based on unique git repo URL.
func (repo *Repository) generateID() string {
	hash := sha1.New()
	_, err := hash.Write([]byte(repo.GitCloneURL))
	if err != nil {
		log.Errorf("Error generating the repository hash: %+v", err)
		return ""
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// generateSlug generates a readable unique string based on repository name.
func (repo *Repository) generateSlug() string {
	vendorAndName := strings.Replace(repo.Name, "/", "-", -1)
	vendorAndName = strings.ReplaceAll(vendorAndName, ".", "_")

	if repo.Pa.CodiceIPA == "" {
		ID := repo.generateID()
		return fmt.Sprintf("%s-%s", vendorAndName, ID[0:6])
	}

	return fmt.Sprintf("%s-%s", repo.Pa.CodiceIPA, vendorAndName)
}

// DeleteByQueryFromES delete record from elasticsearch
// that will match search sting for publiccode.url field
func (c *Crawler) DeleteByQueryFromES(search string) error {
	// Search with a term query
	termQuery := elastic.NewTermQuery("publiccode.url", search)

	// Put publiccode data in ES.
	ctx := context.Background()
	searchResult, err := c.es.DeleteByQuery().
		Index(c.index).
		Type("software").
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

	log.Infof("Deleted %d record from ES", searchResult.Deleted)
	return nil
}
