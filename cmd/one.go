package cmd

import (
	"strconv"
	"sync"
	"time"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/italia/developers-italia-backend/metrics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(oneCmd)
}

var oneCmd = &cobra.Command{
	Use:   "one [domain ID] [repo url]",
	Short: "Crawl publiccode.yml from one single [repo url] using [domain ID] configs.",
	Long: `Crawl publiccode.yml from one [repo url] using [domain ID] configs.
	The domainID should be one in the domains.yml list`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Read args from command.
		domainID := args[0]
		repo := args[1]

		// Register client API plugins.
		crawler.RegisterClientAPIs()

		// Elastic connection.
		elasticClient, err := crawler.ElasticClientFactory(
			viper.GetString("ELASTIC_URL"),
			viper.GetString("ELASTIC_USER"),
			viper.GetString("ELASTIC_PWD"))
		if err != nil {
			panic(err)
		}

		// Read and parse list of domains.
		domains, err := crawler.ReadAndParseDomains(domainsFile)
		if err != nil {
			panic(err)
		}

		// Initiate a channel of repositories.
		repositories := make(chan crawler.Repository, 1)
		// Prepare WaitGroup.
		var wg sync.WaitGroup

		// Index for actual process.
		index := strconv.FormatInt(time.Now().Unix(), 10)

		// Register Prometheus metrics.
		metrics.RegisterPrometheusCounter("repository_processed", "Number of repository processed.", index)
		metrics.RegisterPrometheusCounter("repository_file_saved", "Number of file saved.", index)
		metrics.RegisterPrometheusCounter("repository_file_indexed", "Number of file indexed.", index)
		//metrics.RegisterPrometheusCounter("repository_file_saved_valid", "Number of valid file saved.", index)

		// Process each domain service.
		for _, domain := range domains {
			// get the correct domain ID
			if domain.ID == domainID {
				log.Debugf("Processing domain: %s - Repo: %s", domainID, repo)
				err = crawler.ProcessSingleRepository(repo, domain, repositories)
				if err != nil {
					log.Error(err)
					return
				}

			}
		}

		// Start the metrics server.
		go metrics.StartPrometheusMetricsServer()

		// WaitingLoop check and close the repositories channel
		go crawler.WaitingLoop(repositories, index, &wg, elasticClient)

		// Process the repositories in order to retrieve the file.
		crawler.ProcessRepositories(repositories, index, &wg, elasticClient)

	},
}
