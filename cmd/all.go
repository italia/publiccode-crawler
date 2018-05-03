package cmd

import (
	"os"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(allCmd)
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Crawl publiccode.yml from domains.",
	Long: `Start the crawler on every host written on domains.yml file.
Beware! May take days to complete.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Register client API plugins.
		crawler.RegisterClientApis()

		// Redis connection.
		redisClient, err := crawler.RedisClientFactory(os.Getenv("REDIS_URL"))
		if err != nil {
			panic(err)
		}

		domainsFile := "domains.yml"
		domains, err := crawler.ReadAndParseDomains(domainsFile, redisClient)
		if err != nil {
			panic(err)
		}

		// Initiate a channel of repositories.
		repositories := make(chan crawler.Repository, 100)
		// Initiate a channel for domains status.
		domainsStatus := make(chan int, len(domains))

		// Process each domain service.
		for _, domain := range domains {
			domainsStatus <- 1
			go crawler.ProcessDomain(domain, repositories, domainsStatus)
		}

		// Process repositories in order to retrieve publiccode.yml.
		crawler.ProcessRepositories(repositories, domainsStatus)
	},
}
