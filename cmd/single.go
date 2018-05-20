package cmd

import (
	"sync"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/italia/developers-italia-backend/metrics"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(singleCmd)
	singleCmd.Flags().BoolVarP(&restartCrawling, "restart", "r", false, "Ignore interrupted jobs and restart from the beginning.")
}

var restartCrawling bool

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
		redisClient, err := crawler.RedisClientFactory(viper.GetString("REDIS_URL"))
		if err != nil {
			panic(err)
		}

		// Elastic connection.
		elasticClient, err := crawler.ElasticClientFactory(
			viper.GetString("ELASTIC_URL"),
			viper.GetString("ELASTIC_USER"),
			viper.GetString("ELASTIC_PWD"))
		if err != nil {
			panic(err)
		}

		domainsFile := "domains.yml"
		domains, err := crawler.ReadAndParseDomains(domainsFile, redisClient, restartCrawling)
		if err != nil {
			panic(err)
		}

		// Index for current running id.
		index, err := crawler.UpdateIndex(domains, redisClient, elasticClient)
		if err != nil {
			panic(err)
		}

		log.Debugf("Index %s", index)

		// Initiate a channel of repositories.
		repositories := make(chan crawler.Repository, 1000)
		// Prepare WaitGroup.
		var wg sync.WaitGroup

		// Process each domain service.
		for _, domain := range domains {
			if domain.Id == domainID {
				wg.Add(1)
				// Register single domain metrics.
				metrics.RegisterPrometheusCounter(domain.Id, "Counter for "+domain.Id)
				// Start the process of repositories list.
				go crawler.ProcessDomain(domain, repositories, &wg)
			}
		}

		// Process the repositories in order to retrieve publiccode.yml.
		go crawler.ProcessRepositories(repositories, index, &wg, elasticClient)

		// Wait until all the domains and repositories are processed.
		crawler.WaitingLoop(repositories, index, &wg, elasticClient)
	},
}
