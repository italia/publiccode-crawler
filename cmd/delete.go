package cmd

import (
	"github.com/italia/developers-italia-backend/crawler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete [repo url]",
	Short: "Delete from ElasticSearch one single [repo url].",
	Long: `Delete from ElasticSearch a single repository defined with [repo url].
		No organizations! Only single repositories!`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := crawler.NewCrawler(false)

		err := c.DeleteByQueryFromES(args[0])
		if err != nil {
			log.Error(err)
		}

		// Generate the data files for Jekyll.
		err = c.ExportForJekyll()
		if err != nil {
			log.Errorf("Error while exporting data for Jekyll: %v", err)
		}
	},
}
