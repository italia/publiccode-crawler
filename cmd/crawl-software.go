package cmd

import (
	"github.com/italia/publiccode-crawler/v4/common"
	"github.com/italia/publiccode-crawler/v4/crawler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	crawlSoftwareCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "perform a dry run with no changes made")

	rootCmd.AddCommand(crawlSoftwareCmd)
}

var crawlSoftwareCmd = &cobra.Command{
	Use:   "crawl-software [SOFTWARE_ID | SOFTWARE_URL] PUBLISHER_ID",
	Short: "Crawl a single software by its id.",
	Long: `Crawl a single software by its id.

Crawl a single software given its API id and its publisher.`,
	Example: "# Crawl just the specified software\n" +
		"publiccode-crawler crawl-software" +
		" https://api.developers.italia.it/v1/software/af6056fc-b2b2-4d31-9961-c9bd94e32bd4 PCM",

	Args: cobra.ExactArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		if token := viper.GetString("GITHUB_TOKEN"); token == "" {
			log.Fatal("Please set GITHUB_TOKEN, it's needed to use the GitHub API'")
		}

		c := crawler.NewCrawler(dryRun)

		publisher := common.Publisher{
			ID: args[1],
		}

		if err := c.CrawlSoftwareByID(args[0], publisher); err != nil {
			log.Fatal(err)
		}
	},
}
