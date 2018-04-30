package cmd

import (
	"github.com/italia/developers-italia-backend/crawler"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(allCmd)
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Crawl publiccode.yml from hostings.",
	Long: `Start the crawler on every host written on hosting.yml file.
Beware! May take days to complete.`,
	Run: func(cmd *cobra.Command, args []string) {
		hostings, err := crawler.ReadAndParseHosting()
		if err != nil {
			panic(err)
		}

		// Initiate a channel of repositories.
		repositories := make(chan crawler.Repository)

		// Process each hosting service.
		for _, hosting := range hostings {
			go crawler.ProcessHosting(hosting, repositories)
		}

		// Process repositories in order to retrieve publiccode.yml.
		crawler.ProcessRepositories(repositories)
	},
}
