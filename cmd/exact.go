package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(exactCmd)
}

var exactCmd = &cobra.Command{
	Use:   "exact [domain ID] [repo url]",
	Short: "Crawl publiccode.yml from exact [repo url] using [domain ID] configs.",
	Long: `Crawl publiccode.yml from exact [repo url] using [domain ID] configs.
	The domainID should be `,
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

		domainsFile := "domains.yml"
		domains, err := crawler.ReadAndParseDomains(domainsFile, redisClient)
		if err != nil {
			panic(err)
		}

		// Initiate a channel of repositories.
		repositories := make(chan crawler.Repository, 1)

		// Process each domain service.
		for _, domain := range domains {

			// get the correct domain ID
			if domain.Id == domainID {
				u, err := url.Parse(repo)
				if err != nil {
					fmt.Println(err)
				}

				// Clear the url.
				fullName := u.Path
				if u.Path[:1] == "/" {
					fullName = fullName[1:]
				}
				if u.Path[len(u.Path)-1:] == "/" {
					fullName = fullName[:len(u.Path)-2]
				}

				err = crawler.ProcessSingleDomain(repo, domain, repositories)
				if err != nil {
					panic(err)
				}

			}
		}

		// Process the repositories in order to retrieve publiccode.yml.
		crawler.ProcessRepositories(repositories)
	},
}
