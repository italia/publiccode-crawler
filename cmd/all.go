package cmd

import (
	"os"
	"sync"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/italia/developers-italia-backend/metrics"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(allCmd)
	allCmd.Flags().BoolVarP(&restartCrawling, "restart", "r", false, "Ignore interrupted jobs and restart from the beginning.")
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Crawl publiccode.yml from domains.",
	Long: `Start the crawler on every host written on domains.yml file.
Beware! May take days to complete.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Register client API plugins.
		crawler.RegisterCrawlers()

		// Redis connection.
		redisClient, err := crawler.RedisClientFactory(os.Getenv("REDIS_URL"))
		if err != nil {
			panic(err)
		}

		domainsFile := "domains.yml"
		domains, err := crawler.ReadAndParseDomains(domainsFile, redisClient, restartCrawling)
		if err != nil {
			panic(err)
		}

		// Initiate a channel of repositories.
		repositories := make(chan crawler.Repository, 1000)
		// Prepare WaitGroup.
		var wg sync.WaitGroup

		// Process each domain service.
		for _, domain := range domains {
			wg.Add(1)
			// Register single domain metrics.
			metrics.RegisterPrometheusCounter(domain.Id, "Counter for "+domain.Id)
			// Start the process of repositories list.
			go crawler.ProcessDomain(domain, repositories, &wg)
		}

		// Process the repositories in order to retrieve publiccode.yml.
		go crawler.ProcessRepositories(repositories, &wg)

		// Start metrics server.
		go metrics.StartPrometheusMetricsServer()

		// Wait until all the domains and repositories are processed.
		crawler.WaitingLoop(repositories, &wg)
	},
}
