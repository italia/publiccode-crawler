package cmd

import (
	"os"

	"github.com/italia/developers-italia-backend/crawler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sync"
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
		domainID := args[0]
		repo := args[1]

		// Register client API plugins.
		crawler.RegisterClientApis()

		// Redis connection.
		redisClient, err := crawler.RedisClientFactory(os.Getenv("REDIS_URL"))
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

		// Read and parse list of domains.
		domainsFile := "domains.yml"
		domains, err := crawler.ReadAndParseDomains(domainsFile, redisClient, false)
		if err != nil {
			panic(err)
		}

		// Retrieve the current index
		// TODO: implement this
		index := ""

		log.Debugf("Index %s", index)

		// Initiate a channel of repositories.
		repositories := make(chan crawler.Repository, 1)
		// Prepare WaitGroup.
		var wg sync.WaitGroup

		// Process each domain service.
		for _, domain := range domains {
			// get the correct domain ID
			if domain.Id == domainID {
				err = crawler.ProcessSingleRepository(repo, domain, repositories)
				if err != nil {
					log.Error(err)
					return
				}

			}
		}

		// Process the repositories in order to retrieve publiccode.yml.
		crawler.ProcessRepositories(repositories, index, &wg, elasticClient)
	},
}
