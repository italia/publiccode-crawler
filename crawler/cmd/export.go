package cmd

import (
	"github.com/italia/developers-italia-backend/crawler/crawler"
	"github.com/italia/developers-italia-backend/crawler/jekyll"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(exportCmd)
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export YAML files.",
	Long:  `Export YAML files for the front end.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Make sure the output directory exists or spit an error
		if stat, err := os.Stat(viper.GetString("CRAWLER_DATADIR")); err != nil || !stat.IsDir() {
			log.Fatalf("The configured data directory (%v) does not exist: %v", viper.GetString("CRAWLER_DATADIR"), err)
		}

		// Elastic connection.
		log.Debug("Connecting to ElasticSearch...")
		elasticClient, err := crawler.ElasticClientFactory(
			viper.GetString("ELASTIC_URL"),
			viper.GetString("ELASTIC_USER"),
			viper.GetString("ELASTIC_PWD"))
		if err != nil {
			log.Fatal(err)
		}
		log.Debug("Successfully connected to ElasticSearch")

		// Generate the jekyll files.
		err = jekyll.GenerateJekyllYML(elasticClient)
		if err != nil {
			log.Errorf("Error generating Jekyll yml data: %v", err)
		}
	}}
