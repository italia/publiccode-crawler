package jekyll

import (
	"fmt"
	"os"
	"path"

	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// GenerateJekyllYML generate all the yml files that will be used by Jekyll to generate the static site.
func GenerateJekyllYML(elasticClient *elastic.Client) error {
	// Make sure the output directory exists or spit an error
	outputDir := viper.GetString("OUTPUT_DIR")
	if stat, err := os.Stat(outputDir); err != nil || !stat.IsDir() {
		log.Fatalf("The configured output directory (%v) does not exist: %v", outputDir, err)
	}

	// unsupportedCountries list.
	unsupportedCountries := viper.GetStringSlice("IGNORE_UNSUPPORTEDCOUNTRIES")
	fmt.Println(unsupportedCountries)

	// Create and populate amministrazioni.yml
	amministrazioniFilePath := path.Join(outputDir, "amministrazioni.yml")
	err := AmministrazioniYML(amministrazioniFilePath, unsupportedCountries, elasticClient)
	if err != nil {
		log.Error(err)
	}

	// Create and populate software-riuso.yml
	softwareRiusoFilePath := path.Join(outputDir, "software-riuso.yml")
	numberOfSoftwareRiuso := 4
	err = FirstSoftwareRiuso(softwareRiusoFilePath, numberOfSoftwareRiuso, unsupportedCountries, elasticClient)
	if err != nil {
		log.Error(err)
	}

	// Create and populate software-open-source.yml
	softwareOSFilePath := path.Join(outputDir, "software-open-source.yml")
	numberOfSoftwareOS := 4
	err = FirstSoftwareOpenSource(softwareOSFilePath, numberOfSoftwareOS, unsupportedCountries, elasticClient)
	if err != nil {
		log.Error(err)
	}

	// Create and populate softwares.yml
	softwaresFilePath := path.Join(outputDir, "softwares.yml")
	numberOfSimilarSoftware := 4
	numberOfPopularTags := 5
	err = AllSoftwareYML(softwaresFilePath, numberOfSimilarSoftware, numberOfPopularTags, unsupportedCountries, elasticClient)
	if err != nil {
		log.Errorf("Error exporting jekyll file of all the software : %v", err)
	}

	// Create and populate software_tags.yml
	softwaresTagsDestFilePath := path.Join(outputDir, "software_tags.yml")
	err = TagsYML(softwaresTagsDestFilePath, elasticClient)
	if err != nil {
		log.Errorf("Error exporting jekyll file of software tags: %v", err)
	}

	return err
}
