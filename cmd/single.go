package cmd

import (
	"os"
	"sync"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(singleCmd)
}

var singleCmd = &cobra.Command{
	Use:   "single [domain id]",
	Short: "Crawl publiccode.yml from [domain id].",
	Long: `Start the crawler on [domain id] host defined on domains.yml file.
Beware! May take days to complete.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domainID := args[0]

		// Register client API plugins.
		crawler.RegisterCrawlers()

		// Redis connection.
		redisClient, err := crawler.RedisClientFactory(os.Getenv("REDIS_URL"))
		if err != nil {
			panic(err)
		}

		// Elastic connection.
		elasticClient, err := crawler.ElasticClientFactory(
			os.Getenv("ELASTIC_URL"),
			os.Getenv("ELASTIC_USER"),
			os.Getenv("ELASTIC_PWD"))
		if err != nil {
			panic(err)
		}

		domainsFile := "domains.yml"
		domains, err := crawler.ReadAndParseDomains(domainsFile, redisClient, elasticClient)
		if err != nil {
			panic(err)
		}

		// Index for current running id.
		index, err := crawler.UpdateIndex(domains, redisClient, elasticClient)
		if err != nil {
			panic(err)
		}

		// Initiate a channel of repositories.
		repositories := make(chan crawler.Repository, 1000)
		// Prepare WaitGroup.
		var wg sync.WaitGroup

		// Process each domain service.
		for _, domain := range domains {
			if domain.Id == domainID {
				wg.Add(1)
				go crawler.ProcessDomain(domain, repositories, &wg)
			}
		}

		// Process the repositories in order to retrieve publiccode.yml.
		go crawler.ProcessRepositories(repositories, index, &wg, elasticClient)

		// Wait until all the domains and repositories are processed.
		crawler.WaitingLoop(repositories, index, &wg, elasticClient)
	},
}
