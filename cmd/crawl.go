package cmd

import (
	"strconv"
	"sync"
	"time"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/italia/developers-italia-backend/ipa"
	"github.com/italia/developers-italia-backend/jekyll"
	"github.com/italia/developers-italia-backend/metrics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(crawlCmd)
}

var crawlCmd = &cobra.Command{
	Use:   "crawl whitelist.yml [whitelistGeneric.yml whitelistPA.yml ...]",
	Short: "Crawl publiccode.yml file from domains in whitelist file.",
	Long:  `Start whitelist file. It's possible to add multiple files adding them as args.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Update ipa to lastest data.
		err := ipa.UpdateFile("./ipa/amministrazioni.txt", "http://www.indicepa.gov.it/public-services/opendata-read-service.php?dstype=FS&filename=amministrazioni.txt")
		if err != nil {
			log.Fatal(err)
		}
		// Index for actual process.
		index := strconv.FormatInt(time.Now().Unix(), 10)

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

		// Initiate a channel of repositories.
		repositories := make(chan crawler.Repository, 1000)
		// Prepare WaitGroup.
		var wg sync.WaitGroup

		// Register Prometheus metrics.
		metrics.RegisterPrometheusCounter("repository_processed", "Number of repository processed.", index)
		metrics.RegisterPrometheusCounter("repository_file_saved", "Number of file saved.", index)
		metrics.RegisterPrometheusCounter("repository_file_indexed", "Number of file indexed.", index)
		metrics.RegisterPrometheusCounter("repository_cloned", "Number of repository cloned", index)
		// Uncomment when validating publiccode.yml
		//metrics.RegisterPrometheusCounter("repository_file_saved_valid", "Number of valid file saved.", index)

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
		for repository := range repositories {
			wg.Add(1)
			go crawler.CheckAvailability(repository, index, &wg, elasticClient)
		}
		wg.Wait()

		// ElasticFlush to flush all the operations on ES.
		err = crawler.ElasticFlush(index, elasticClient)
		if err != nil {
			log.Errorf("Error flushing ElasticSearch: %v", err)
		}

		// Update Elastic alias.
		err = crawler.ElasticAliasUpdate(index, "publiccode", elasticClient)
		if err != nil {
			log.Errorf("Error updating Elastic Alias: %v", err)
		}

		// Generate the jekyll files.
		err = jekyll.GenerateJekyllYML(elasticClient)
		if err != nil {
			log.Errorf("Error generating Jekyll yml data: %v", err)
		}
	}}
