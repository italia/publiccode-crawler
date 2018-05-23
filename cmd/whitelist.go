package cmd

import (
	"sync"
	"time"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(whitelistCmd)
}

var whitelistCmd = &cobra.Command{
	Use:   "whitelist",
	Short: "Crawl publiccode.yml from domains in whitelist.",
	Long:  `Start the crawler on whitelist.yml file.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Elastic connection.
		elasticClient, err := crawler.ElasticClientFactory(
			viper.GetString("ELASTIC_URL"),
			viper.GetString("ELASTIC_USER"),
			viper.GetString("ELASTIC_PWD"))
		if err != nil {
			panic(err)
		}

		// Read and parse list of domains.
		domainsFile := "domains.yml"
		domains, err := crawler.ReadAndParseDomains(domainsFile)
		if err != nil {
			panic(err)
		}

		// Read and parse the whitelist.
		whitelistFile := "whitelist.yml"
		whitelist, err := crawler.ReadAndParseWhitelist(whitelistFile)
		if err != nil {
			panic(err)
		}

		// Initiate a channel of repositories.
		repositories := make(chan crawler.Repository, 1000)
		// Prepare WaitGroup.
		var wg sync.WaitGroup

		// Index for actual process.
		index := time.Now().String()

		// Process every item in whitelist
		for _, pa := range whitelist {
			wg.Add(1)
			go crawler.ProcessPA(pa, domains, repositories, index, &wg)
		}

		go crawler.ProcessRepositories(repositories, index, &wg, elasticClient)
		//
		//go crawler.WaitingLoop(repositories, index, wg, elasticClient)

		wg.Wait()
	}}
