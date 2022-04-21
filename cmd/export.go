package cmd

import (
	"github.com/italia/developers-italia-backend/crawler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(exportCmd)
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export YAML files.",
	Long:  `Export YAML files for the front end.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := crawler.NewCrawler(false)

		// Generate the data files for Jekyll.
		err := c.ExportForJekyll()
		if err != nil {
			log.Errorf("Error while exporting data for Jekyll: %v", err)
		}
	}}
