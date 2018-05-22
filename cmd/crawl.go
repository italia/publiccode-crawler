package cmd

import (
	"sync"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/italia/developers-italia-backend/metrics"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var restartCrawling bool
var domainToCrawl string

func init() {
	rootCmd.AddCommand(crawlCmd)
	crawlCmd.Flags().BoolVarP(&restartCrawling, "restart", "r", false, "Ignore interrupted jobs and restart from the beginning.")
	crawlCmd.Flags().StringVarP(&domainToCrawl, "domain", "d", "", "Import repositories only from this domain.")
}

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl publiccode.yml from domains.",
	Long: `Start the crawler on every host written on domains.yml file.
Beware! May take days to complete.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		// Register Prometheus metrics.
		metrics.RegisterPrometheusCounter("repository_processed", "Number of repository processed.", index)
		metrics.RegisterPrometheusCounter("repository_file_saved", "Number of file saved.", index)
		metrics.RegisterPrometheusCounter("repository_file_indexed", "Number of file indexed.", index)
		metrics.RegisterPrometheusCounter("repository_file_saved_valid", "Number of valid file saved.", index)

		// Process each domain service.
		for _, domain := range domains {
			// If domainToCrawl is empty crawl all domains, otherwise crawl only the one with Id equals to domainToCrawl.
			if (domainToCrawl == "") || (domainToCrawl != "" && domain.Id == domainToCrawl) {
				wg.Add(1)

				// Register Prometheus metrics.
				metrics.RegisterPrometheusCounter("repository_"+domain.Id+"_processed", "Counter for "+domain.Id, index)

				// Start the process of repositories list.
				go crawler.ProcessDomain(domain, redisClient, repositories, index, &wg)
			}
		}

		// Process the repositories in order to retrieve publiccode.yml.
		go crawler.ProcessRepositories(repositories, index, &wg, elasticClient)

		// Start the metrics server.
		go metrics.StartPrometheusMetricsServer()

		// Wait until all the domains and repositories are processed.
		crawler.WaitingLoop(repositories, index, &wg, elasticClient)
	},
}
