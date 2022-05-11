package jekyll

import (
	"os"
	"path"

	elastic "github.com/olivere/elastic/v7"
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

	// Create and populate amministrazioni.yml
	amministrazioniFilePath := path.Join(outputDir, "amministrazioni.yml")
	err := AmministrazioniYML(amministrazioniFilePath, elasticClient)
	if err != nil {
		log.Error(err)
	}

	// Create and populate software-riuso.yml
	softwareRiusoFilePath := path.Join(outputDir, "software-riuso.yml")
	numberOfSoftwareRiuso := 4
	err = FirstSoftwareRiuso(softwareRiusoFilePath, numberOfSoftwareRiuso, elasticClient)
	if err != nil {
		log.Error(err)
	}

	// Create and populate software-open-source.yml
	softwareOSFilePath := path.Join(outputDir, "software-open-source.yml")
	numberOfSoftwareOS := 4
	err = FirstSoftwareOpenSource(softwareOSFilePath, numberOfSoftwareOS, elasticClient)
	if err != nil {
		log.Error(err)
	}

	// Create and populate softwares.yml
	softwaresFilePath := path.Join(outputDir, "softwares.yml")
	numberOfSimilarSoftware := 4
	numberOfPopularCategories := 5
	err = AllSoftwareYML(softwaresFilePath, numberOfSimilarSoftware, numberOfPopularCategories, elasticClient)
	if err != nil {
		log.Errorf("Error exporting jekyll file of all the software : %v", err)
	}

	// Export the list of distinct categories mentioned in the catalog
	err = CategoriesYML(path.Join(outputDir, "software_categories.yml"), elasticClient)
	if err != nil {
		log.Errorf("Error exporting jekyll file of software categories: %v", err)
	}

	// Export the list of distinct scopes mentioned in the catalog
	err = ScopesYML(path.Join(outputDir, "software_scopes.yml"), elasticClient)
	if err != nil {
		log.Errorf("Error exporting jekyll file of software scopes: %v", err)
	}

	return err
}
