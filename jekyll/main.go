package jekyll

import (
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

// GenerateJekyllYML generate all the yml files that will be used by Jekyll to generate the static site.
func GenerateJekyllYML(elasticClient *elastic.Client) error {
	// Create and populate amministrazioni.yml
	amministrazioniFilePath := "jekyll/generated/amministrazioni.yml"
	err := AmministrazioniYML(amministrazioniFilePath, elasticClient)
	if err != nil {
		log.Error(err)
	}

	// Create and populate software-riuso.yml
	softwareRiusoFilePath := "jekyll/generated/software-riuso.yml"
	numberOfSoftwareRiuso := 4
	err = FirstSoftwareRiuso(softwareRiusoFilePath, numberOfSoftwareRiuso, elasticClient)
	if err != nil {
		log.Error(err)
	}

	// Create and populate software-open-source.yml
	softwareOSFilePath := "jekyll/generated/software-open-source.yml"
	numberOfSoftwareOS := 4
	err = FirstSoftwareOpenSource(softwareOSFilePath, numberOfSoftwareOS, elasticClient)
	if err != nil {
		log.Error(err)
	}

	// Create and populate softwares.yml
	softwaresFilePath := "jekyll/generated/softwares.yml"
	numberOfSimilarSoftware := 4
	err = AllSoftwareYML(softwaresFilePath, numberOfSimilarSoftware, elasticClient)
	if err != nil {
		log.Errorf("Error exporting jekyll file of all the software : %v", err)
	}

	return err
}
