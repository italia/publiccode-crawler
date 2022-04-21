package cmd

import (
	"github.com/italia/developers-italia-backend/crawler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	crawlCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "perform a dry run with no changes made")

	rootCmd.AddCommand(crawlCmd)
}

var crawlCmd = &cobra.Command{
	Use:   "crawl whitelist.yml whitelist/*.yml",
	Short: "Crawl publiccode.yml files from given domains.",
	Long:  `Crawl publiccode.yml files according to the supplied whitelist file(s).`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		orgs := make(map[string]bool)
		c := crawler.NewCrawler(dryRun)

		// Read the supplied whitelists.
		var publishers []crawler.Publisher
		for id := range args {
			readWhitelist, err := crawler.ReadAndParseWhitelist(args[id])
			if err != nil {
				log.Fatal(err)
			}

		Publisher:
			for _, publisher := range readWhitelist {
				for _, orgURL := range publisher.Organizations {
					if orgs[orgURL.String()] {
						log.Warnf("Skipping publisher '%s': organization '%s' already present", publisher.Name, orgURL.String())
						continue Publisher
					} else {
						orgs[orgURL.String()] = true
					}
				}
				publishers = append(publishers, publisher)
			}
		}

		toBeRemoved, err := c.CrawlPublishers(publishers)
		if err != nil {
			log.Fatal(err)
		}

		// I should call delete for items in blacklist
		// to ensure they are not present in ES and then in
		// jekyll datafile
		for _, repo := range toBeRemoved {
			log.Warnf("blacklisted, going to remove from ES %s", repo)
			err = c.DeleteByQueryFromES(repo)
			if err != nil {
				log.Errorf("Error while deleting data from ES: %v", err)
			}
		}

		// Generate the data files for Jekyll.
		err = c.ExportForJekyll()
		if err != nil {
			log.Errorf("Error while exporting data for Jekyll: %v", err)
		}
	}}
