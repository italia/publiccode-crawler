package cmd

import (
	"net/url"
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
	Use:   "one [repo url]",
	Short: "Crawl publiccode.yml from one single [repo url].",
	Long: `Crawl publiccode.yml from a single repository defined with [repo url].
No organizations! Only single repositories!`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Read repository URL.
		repo := args[0]

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

		log.Debugf("Processing Single Repo: %s", repo)

		knownHost := false
		domain := crawler.Domain{}
		// Parse as url.URL.
		u, err := url.Parse(repo)
		if err != nil {
			log.Errorf("invalid host: %v", err)
		}

		// Check if host is in list of "famous" hosts.
		for _, d := range domains {
			if u.Hostname() == d.Host {
				// Process this host
				knownHost = true
				domain = d
			}
		}

		if knownHost {
			log.Infof("%s - API known:%s", repo, u.Hostname())
			// Host is detected.
			err = crawler.ProcessSingleRepository(repo, domain, repositories)
			if err != nil {
				log.Error(err)
				return
			}

		} else {
			// host unknown, needs to be inferred.
			if crawler.IsGithub(repo) {
				log.Infof("%s - API inferred:%s", repo, "github.com")
				err = crawler.ProcessSingleRepository(repo, crawler.Domain{Host: "github.com"}, repositories)
				if err != nil {
					log.Error(err)
					return
				}
			} else if crawler.IsBitbucket(repo) {
				log.Infof("%s - API inferred:%s", repo, "bitbucket.org")
				err = crawler.ProcessSingleRepository(repo, crawler.Domain{Host: "bitbucket.org"}, repositories)
				if err != nil {
					log.Error(err)
					return
				}
			} else if crawler.IsGitlab(repo) {
				log.Infof("%s - API inferred:%s", repo, "gitlab.com")
				err = crawler.ProcessSingleRepository(repo, crawler.Domain{Host: "gitlab.com"}, repositories)
				if err != nil {
					log.Error(err)
					return
				}
			}
		}

		// Start the metrics server.
		go metrics.StartPrometheusMetricsServer()

		// WaitingLoop check and close the repositories channel
		go crawler.WaitingLoop(repositories, &wg)

		// Process the repositories in order to retrieve the file.
		// ProcessRepositories is blocking (wait until repositories is closed by WaitingLoop).
		crawler.ProcessRepositories(repositories, index, &wg, elasticClient)

		// Update Elastic alias.
		err = crawler.ElasticAliasUpdate(index, "publiccode", elasticClient)
		if err != nil {
			log.Errorf("Error updating Elastic Alias: %v", err)
		}

	},
}
