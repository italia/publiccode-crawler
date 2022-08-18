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
	Use:   "crawl publishers.yml directory/*.yml ...",
	Short: "Crawl publiccode.yml files in publishers' repos.",
	Long:  `Crawl publiccode.yml files according to the supplied publisher file(s).`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		orgs := make(map[string]bool)
		c := crawler.NewCrawler(dryRun)

		apiclient := apiclient.NewClient()

		var dedupedPublishers []common.Publisher
		for id := range args {
			publishers, err := apiclient.GetPublishers()
			if err != nil {
				log.Fatal(err)
			} else {
				log.Infof("Loaded and parsed %s", args[id])
			}

		Publisher:
			for _, publisher := range publishers {
				for _, orgURL := range publisher.Organizations {
					if orgs[orgURL.String()] {
						log.Warnf("Skipping publisher '%s': organization '%s' already present", publisher.Name, orgURL.String())
						continue Publisher
					} else {
						orgs[orgURL.String()] = true
					}
				}
				dedupedPublishers = append(dedupedPublishers, publisher)
			}
		}

		if err := c.CrawlPublishers(dedupedPublishers); err != nil {
			log.Fatal(err)
		}

		// Generate the data files for Jekyll.
		err := c.ExportForJekyll()
		if err != nil {
			log.Errorf("Error while exporting data for Jekyll: %v", err)
		}
	},
}
