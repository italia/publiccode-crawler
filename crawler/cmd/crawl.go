package cmd

import (
	"github.com/italia/developers-italia-backend/crawler/crawler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(crawlCmd)
}

var crawlCmd = &cobra.Command{
	Use:   "crawl whitelist.yml whitelist/*.yml",
	Short: "Crawl publiccode.yml files from given domains.",
	Long:  `Crawl publiccode.yml files according to the supplied whitelist file(s).`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := crawler.NewCrawler()

		// Read the supplied whitelists.
		var publishers []crawler.PA
		for id := range args {
			readWhitelist, err := crawler.ReadAndParseWhitelist(args[id])
			if err != nil {
				log.Fatal(err)
			}
			publishers = append(publishers, readWhitelist...)
		}

		// Crawl
		err := c.CrawlOrgs(publishers)
		if err != nil {
			log.Fatal(err)
		}

		// Generate the data files for Jekyll.
		err = c.ExportForJekyll()
		if err != nil {
			log.Errorf("Error while exporting data for Jekyll: %v", err)
		}
	}}
