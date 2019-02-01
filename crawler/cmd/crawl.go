package cmd

import (
	"github.com/italia/developers-italia-backend/crawler/crawler"
	"github.com/italia/developers-italia-backend/crawler/ipa"
	"github.com/italia/developers-italia-backend/crawler/jekyll"
	"github.com/italia/developers-italia-backend/crawler/metrics"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(crawlCmd)
}

var crawlCmd = &cobra.Command{
	Use:   "crawl whitelist.yml whitelist/*.yml",
	Short: "Crawl publiccode.yml files from given domains.",
	Long:  `Crawl publiccode.yml files according to the supplied whitelist file(s).`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Make sure the output directory exists or spit an error
		if stat, err := os.Stat(viper.GetString("CRAWLER_DATADIR")); err != nil || !stat.IsDir() {
			log.Fatalf("The configured data directory (%v) does not exist: %v", viper.GetString("CRAWLER_DATADIR"), err)
		}

		// Update ipa to lastest data.
		err := ipa.UpdateFromIndicePA()
		if err != nil {
			log.Error(err)
		}

		// Elastic connection.
		log.Debug("Connecting to ElasticSearch...")
		elasticClient, err := crawler.ElasticClientFactory(
			viper.GetString("ELASTIC_URL"),
			viper.GetString("ELASTIC_USER"),
			viper.GetString("ELASTIC_PWD"))
		if err != nil {
			log.Fatal(err)
		}
		log.Debug("Successfully connected to ElasticSearch")

		// Create ES index with mapping for PublicCode.
		const index = "publiccodes"  // Index for actual process.
		err = crawler.ElasticIndexMapping(index, elasticClient)
		if err != nil {
			log.Fatal(err)
		}

		// Create ES index with mapping "administration-codiceIPA".
		err = crawler.ElasticAdministrationsMapping("administration", elasticClient)
		if err != nil {
			log.Fatal(err)
		}

		// Read and parse list of domains.
		domains, err := crawler.ReadAndParseDomains(domainsFile)
		if err != nil {
			log.Fatal(err)
		}

		// Read and parse the whitelist.
		var whitelist []crawler.PA

		// Fill the whitelist with all the args whitelists.
		for id := range args {
			readWhitelist, err := crawler.ReadAndParseWhitelist(args[id])
			if err != nil {
				log.Fatal(err)
			}
			whitelist = append(whitelist, readWhitelist...)
		}

		// Count configured orgs
		orgCount := 0
		for _, pa := range whitelist {
			orgCount += len(pa.Organizations)
		}
		log.Infof("%v organizations belonging to %v publishers are going to be scanned",
			orgCount, len(whitelist))

		// Register Prometheus metrics.
		metrics.RegisterPrometheusCounter("repository_processed", "Number of repository processed.", index)
		metrics.RegisterPrometheusCounter("repository_file_saved", "Number of file saved.", index)
		metrics.RegisterPrometheusCounter("repository_file_indexed", "Number of file indexed.", index)
		metrics.RegisterPrometheusCounter("repository_cloned", "Number of repository cloned", index)
		// Uncomment when validating publiccode.yml
		//metrics.RegisterPrometheusCounter("repository_file_saved_valid", "Number of valid file saved.", index)
		
		// Initiate a channel of repositories.
		repositories := make(chan crawler.Repository, 1000)

		// Prepare WaitGroup.
		var wg sync.WaitGroup

		// Process every item in whitelist.
		for _, pa := range whitelist {
			wg.Add(1)
			go crawler.ProcessPA(pa, domains, repositories, &wg)
		}

		// Start the metrics server.
		go metrics.StartPrometheusMetricsServer()

		// WaitingLoop check and close the repositories channel.
		go crawler.WaitingLoop(repositories, &wg)

		// Process the repositories in order to retrieve the file.
		// ProcessRepositories loop is blocking (wait until repositories is closed by WaitingLoop).
		crawler.ProcessRepositories(repositories, index, &wg, elasticClient)

		// ElasticFlush to flush all the operations on ES.
		err = crawler.ElasticFlush(index, elasticClient)
		if err != nil {
			log.Errorf("Error flushing ElasticSearch: %v", err)
		}

		// Update Elastic alias.
		err = crawler.ElasticAliasUpdate("administration", viper.GetString("ELASTIC_ALIAS"), elasticClient)
		if err != nil {
			log.Errorf("Error updating Elastic Alias: %v", err)
		}
		err = crawler.ElasticAliasUpdate(index, viper.GetString("ELASTIC_ALIAS"), elasticClient)
		if err != nil {
			log.Errorf("Error updating Elastic Alias: %v", err)
		}

		// Generate the jekyll files.
		err = jekyll.GenerateJekyllYML(elasticClient)
		if err != nil {
			log.Errorf("Error generating Jekyll yml data: %v", err)
		}
	}}
