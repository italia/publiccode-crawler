package cmd

import (
	"github.com/italia/developers-italia-backend/elastic"
	"github.com/italia/developers-italia-backend/ipa"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(updateIPACmd)
}

var updateIPACmd = &cobra.Command{
	Use:   "updateipa",
	Short: "Update data from IndicePA.",
	Long:  `Download data from IndicePA and inject it into Elasticsearch.`,
	Run: func(cmd *cobra.Command, args []string) {
		es, err := elastic.ClientFactory(
			viper.GetString("ELASTIC_URL"),
			viper.GetString("ELASTIC_USER"),
			viper.GetString("ELASTIC_PWD"))
		if err != nil {
			log.Fatal(err)
		}

		err = ipa.UpdateFromIndicePA(es)
		if err != nil {
			log.Error(err)
		}
	}}
