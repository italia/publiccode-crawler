package cmd

import (
	"github.com/italia/developers-italia-backend/apiclient"
	"github.com/italia/developers-italia-backend/common"
	"github.com/italia/developers-italia-backend/crawler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	crawlCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "perform a dry run with no changes made")

	rootCmd.AddCommand(crawlCmd)
}

var crawlCmd = &cobra.Command{
	Use:   "crawl publishers.yml [directory/*.yml ...]",
	Short: "Crawl publiccode.yml files in publishers' repos.",
	Long:  `Crawl publiccode.yml files in publishers' repos.
				When run with no arguments, the publishers are fetched from the API,
				otherwise the passed YAML files are used.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		c := crawler.NewCrawler(dryRun)

		var publishers []common.Publisher

		if len(args) == 0 {
			var err error

			apiclient := apiclient.NewClient()

			publishers, err = apiclient.GetPublishers()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			for _, yamlFile := range args {
				filePublishers, err := common.LoadPublishers(yamlFile)
				if err != nil {
					log.Fatal(err)
				}

				publishers = append(publishers, filePublishers...)
			}
		}

		if err := c.CrawlPublishers(publishers); err != nil {
			log.Fatal(err)
		}

		// Generate the data files for Jekyll.
		err := c.ExportForJekyll()
		if err != nil {
			log.Errorf("Error while exporting data for Jekyll: %v", err)
		}
	},
}
