package cmd

import (
	"net/url"
	"sync"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/italia/developers-italia-backend/ipa"
	"github.com/italia/developers-italia-backend/jekyll"
	"github.com/italia/developers-italia-backend/metrics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(oneCmd)
}

var oneCmd = &cobra.Command{
	Use:   "one [repo url]",
	Short: "Crawl publiccode.yml from one single [repo url].",
	Long: `Crawl publiccode.yml from a single repository defined with [repo url].
No organizations! Only single repositories!`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Update ipa to lastest data.
		err := ipa.UpdateFile("./ipa/amministrazioni.txt", "http://www.indicepa.gov.it/public-services/opendata-read-service.php?dstype=FS&filename=amministrazioni.txt")
		if err != nil {
			log.Fatal(err)
		}

		// Read repository URL.
		repo := args[0]
		// Index for actual process.
		index := "publiccode"

		// Elastic connection.
		elasticClient, err := crawler.ElasticClientFactory(
			viper.GetString("ELASTIC_URL"),
			viper.GetString("ELASTIC_USER"),
			viper.GetString("ELASTIC_PWD"))
		if err != nil {
			log.Fatal(err)
		}
		err = crawler.ElasticIndexMapping(index, elasticClient)
		if err != nil {
			log.Fatal(err)
		}

		// Read and parse list of domains.
		domains, err := crawler.ReadAndParseDomains(domainsFile)
		if err != nil {
			log.Fatal(err)
		}

		// Initiate a channel of repositories.
		repositories := make(chan crawler.Repository, 1)
		// Prepare WaitGroup.
		var wg sync.WaitGroup

		// Register Prometheus metrics.
		metrics.RegisterPrometheusCounter("repository_processed", "Number of repository processed.", index)
		metrics.RegisterPrometheusCounter("repository_file_saved", "Number of file saved.", index)
		metrics.RegisterPrometheusCounter("repository_file_indexed", "Number of file indexed.", index)
		metrics.RegisterPrometheusCounter("repository_cloned", "Number of repository cloned", index)
		//metrics.RegisterPrometheusCounter("repository_file_saved_valid", "Number of valid file saved.", index)

		log.Infof("Processing Single Repo: %s", repo)

		// Parse as url.URL.
		u, err := url.Parse(repo)
		if err != nil {
			log.Errorf("invalid host: %v", err)
		}

		// Check if current host is in known in domains.yml hosts.
		domain, err := crawler.KnownHost(repo, u.Hostname(), domains)
		if err != nil {
			log.Error(err)
		}

		// Process single repository.
		log.Infof("Start ProcessSingleRepository '%s'", repo)
		err = crawler.ProcessSingleRepository(repo, domain, repositories)
		if err != nil {
			log.Error(err)
			return
		}

		// Start the metrics server.
		go metrics.StartPrometheusMetricsServer()

		// WaitingLoop check and close the repositories channel
		go crawler.WaitingLoop(repositories, &wg)

		// Process the repositories in order to retrieve the file.
		// ProcessRepositories is blocking (wait until repositories is closed by WaitingLoop).
		crawler.ProcessRepositories(repositories, index, &wg, elasticClient)

		log.Infof("End ProcessSingleRepository '%s'", repo)

		// ElasticFlush to flush all the operations on ES.
		err = crawler.ElasticFlush(index, elasticClient)
		if err != nil {
			log.Errorf("Error flushing ElasticSearch: %v", err)
		}

		// Generate the jekyll files.
		err = jekyll.GenerateJekyllYML(elasticClient)
		if err != nil {
			log.Errorf("Error generating Jekyll yml data: %v", err)
		}

	},
}
