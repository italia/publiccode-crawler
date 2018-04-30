package cmd

import (
	"github.com/italia/developers-italia-backend/crawler"
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

		hostings, err := crawler.ReadAndParseHosting()
		if err != nil {
			panic(err)
		}

		// Initiate a channel of repositories.
		repositories := make(chan crawler.Repository)

		// For each host parsed from hosting, Process the repositories.
		for _, hosting := range hostings {
			if hosting.ServiceName == serviceName {
				go crawler.ProcessHosting(hosting, repositories)
			}
		}

		// Process the repositories in order to retrieve publiccode.yml.
		crawler.ProcessRepositories(repositories)
	},
}
