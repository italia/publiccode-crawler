package cmd

import (
	"context"
	"os"
	"sync"

	"github.com/italia/developers-italia-backend/crawler"
	log "github.com/sirupsen/logrus"
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
		index, err := crawler.UpdateIndex(domains, redisClient)
		if err != nil {
			panic(err)
		}
		elasticClient.Alias().Add("publiccode", index).Do(context.Background())
		log.Debugf("Index %s", index)

		// Initiate a channel of repositories.
		repositories := make(chan crawler.Repository, 1000)
		// Prepare WaitGroup.
		var wg sync.WaitGroup

		// Process each domain service.
		for _, domain := range domains {
			wg.Add(1)
			go crawler.ProcessDomain(domain, repositories, &wg)
		}

		// Process the repositories in order to retrieve publiccode.yml.
		go crawler.ProcessRepositories(repositories, index, &wg, elasticClient)

		// Wait until all the domains and repositories are processed.
		crawler.WaitingLoop(repositories, index, &wg, elasticClient)
	},
}
