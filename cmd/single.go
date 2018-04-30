package cmd

import (
	"fmt"
	"io/ioutil"
	"github.com/italia/developers-italia-backend/crawler"
	"github.com/italia/developers-italia-backend/metrics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(singleCmd)
}

var singleCmd = &cobra.Command{
	Use:   "single [hosting]",
	Short: "Crawl publiccode.yml from [hosting].",
	Long: `Start the crawler on [hosting] host defined on hosting.yml file.
Beware! May take days to complete.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := args[0]

		// Init Prometheus for metrics.
		processedCounter := metrics.PrometheusCounter(fmt.Sprintf("repository_processed_%s", serviceName), fmt.Sprintf("Number of repository processed on %s.", serviceName))

		// Open and read hosting file list.
		hostingFile := "hosting.yml"
		data, err := ioutil.ReadFile(hostingFile)
		if err != nil {
			panic(fmt.Sprintf("error in reading %s file: %v", hostingFile, err))
		}
		// Parse hosting file list.
		hostings, err := crawler.ParseHostingFile(data)
		if err != nil {
			panic(fmt.Sprintf("error in parsing %s file: %v", hostingFile, err))
		}
		log.Debug("Loaded and parsed hosting.yml")

		// Initiate a channel of repositories.
		repositories := make(chan crawler.Repository)

		// For each host parsed from hosting, Process the repositories.
		for _, hosting := range hostings {
			if hosting.ServiceName == serviceName {
				go crawler.ProcessHosting(hosting, repositories)
			}
		}

		// Process the repositories in order to retrieve publiccode.yml.
		crawler.ProcessRepositories(repositories, processedCounter)
	},
}
