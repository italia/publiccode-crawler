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
		// Prepare WaitGroup.
		var wg sync.WaitGroup

		// Process each domain service.
		for _, domain := range domains {
			if domain.Id == domainID {
				wg.Add(1)
				go crawler.ProcessDomain(domain, repositories, &wg)
			}
		}

		// Goroutine that constatly check if all the processes are terminated.
		go crawler.WaitingLoop(repositories, &wg)

		// Process the repositories in order to retrieve publiccode.yml.
		crawler.ProcessRepositories(repositories, &wg)
	},
}
