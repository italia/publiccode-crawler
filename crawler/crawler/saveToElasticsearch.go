package crawler

import (
	"context"
	"strings"
	"time"

	"github.com/italia/developers-italia-backend/crawler/ipa"
	"github.com/italia/developers-italia-backend/crawler/metrics"
	pcode "github.com/italia/publiccode-parser-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/ghodss/yaml"
)

type administration struct {
	Name      string `json:"it-riuso-codiceIPA-label"`
	CodiceIPA string `json:"it-riuso-codiceIPA"`
}

// SaveToES save the chosen data []byte in elasticsearch
// data contains the raw publiccode.yml file
func (c *Crawler) SaveToES(fileRawURL, hashedRepoURL string, activityIndex float64, vitality []int, data []byte) error {
	// Parse the publiccode.yml file
	parser := pcode.NewParser()
	parser.RemoteBaseURL = strings.TrimRight(fileRawURL, viper.GetString("CRAWLED_FILENAME"))
	err := parser.Parse(data)
	if err != nil {
		log.Errorf("Error parsing publiccode.yml: %v", err)
	}

	// Create a SoftwareES object and populate it
	file := SoftwareES{
		FileRawURL:            fileRawURL,
		ID:                    hashedRepoURL,
		CrawlTime:             time.Now().Format(time.RFC3339),
		ItRiusoCodiceIPALabel: ipa.GetAdministrationName(parser.PublicCode.It.Riuso.CodiceIPA),
		VitalityScore:     activityIndex,
		VitalityDataChart: vitality,
		OEmbedHTML: parser.OEmbed,
	}

	// Convert parser.PublicCode to YAML and parse it again into the SoftwareES record
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
		Id(hashedRepoURL).
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



