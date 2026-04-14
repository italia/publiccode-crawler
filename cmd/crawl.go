package cmd

import (
	"github.com/italia/publiccode-crawler/v4/apiclient"
	"github.com/italia/publiccode-crawler/v4/common"
	"github.com/italia/publiccode-crawler/v4/crawler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	crawlCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "perform a dry run with no changes made")

	rootCmd.AddCommand(crawlCmd)
}

var crawlCmd = &cobra.Command{
	Use:   "crawl [publishers.yml] [directory/*.yml ...]",
	Short: "Crawl publiccode.yml files from catalogs or publishers.",
	Long: `Crawl publiccode.yml files from catalogs or publishers.

When run with no arguments, the catalogs are fetched from the API.
If no catalogs are found, the publishers are fetched as fallback.
When YAML files are passed, they are used as publisher definitions.`,
	Example: `
# Crawl catalogs fetched from the API (with publisher fallback)
crawl

# Crawl using a specific publishers.yml file
crawl publishers.yml

# Crawl all YAML files in a specific directory
crawl directory/*.yml`,

	Args: cobra.MinimumNArgs(0),
	Run: func(_ *cobra.Command, args []string) {
		if token := viper.GetString("GITHUB_TOKEN"); token == "" {
			log.Fatal("Please set GITHUB_TOKEN, it's needed to use the GitHub API'")
		}

		c := crawler.NewCrawler(dryRun)

		if len(args) > 0 {
			crawlFromYAML(c, args)

			return
		}

		crawlFromAPI(c)
	},
}

func crawlFromAPI(c *crawler.Crawler) {
	client := apiclient.NewClient()

	catalogs, err := client.GetCatalogs()
	if err != nil {
		log.Warnf("Failed to get catalogs: %s, falling back to publishers", err)
	}

	if len(catalogs) > 0 {
		if err := c.CrawlCatalogs(catalogs); err != nil {
			log.Fatal(err)
		}

		return
	}

	log.Info("No catalogs found, falling back to publishers")

	publishers, err := client.GetPublishers()
	if err != nil {
		log.Fatal(err)
	}

	if err := c.CrawlPublishers(publishers); err != nil {
		log.Fatal(err)
	}
}

func crawlFromYAML(c *crawler.Crawler, args []string) {
	var publishers []common.Publisher

	for _, yamlFile := range args {
		filePublishers, err := common.LoadPublishers(yamlFile)
		if err != nil {
			log.Fatal(err)
		}

		publishers = append(publishers, filePublishers...)
	}

	if err := c.CrawlPublishers(publishers); err != nil {
		log.Fatal(err)
	}
}
