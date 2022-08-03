package cmd

import (
	"net/url"
	"regexp"

	"github.com/italia/developers-italia-backend/crawler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	oneCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "perform a dry run with no changes made")

	rootCmd.AddCommand(oneCmd)
}

var oneCmd = &cobra.Command{
	Use:   "one [repo url] publishers.*.yml",
	Short: "Crawl publiccode.yml from one single [repo url].",
	Long: `Crawl publiccode.yml from a single repository defined with [repo url]
		according to the supplied file(s).
		No organizations! Only single repositories!`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// check if repo url is not present in blacklist
		// if so report error and exit.
		if crawler.IsRepoInBlackList(args[0]) {
			return
		}

		c := crawler.NewCrawler(dryRun)

		paths := args[1:]
		url, err := url.Parse(args[0])
		if err != nil {
			log.Error(err)
		}

		err = c.CrawlRepo(*url, getPublisher(*url, paths))
		if err != nil {
			log.Error(err)
		}

		// Generate the data files for Jekyll.
		err = c.ExportForJekyll()
		if err != nil {
			log.Errorf("Error while exporting data for Jekyll: %v", err)
		}
	},
}

func getPublisher(repoURL url.URL, paths []string) (p crawler.Publisher) {
	var publishers []crawler.Publisher
	for _, path := range paths {
		p, err := crawler.LoadPublishers(path)
		if err != nil {
			log.Fatal(err)
		}
		publishers = append(publishers, p...)
	}

	for _, publisher := range publishers {
		// looking into repositories
		for _, repo := range publisher.Repositories {
			log.Tracef("matching %s with %s", repo.String(), repoURL.String())
			if (url.URL)(repo) == repoURL {
				log.Debugf("Publisher found %+v", publisher)
				return publisher
			}
		}
		// looking into organizations
		for _, repo := range publisher.Organizations {
			log.Tracef("matching %s.* with %s", repo.String(), repoURL.String())
			if matched, _ := regexp.MatchString(repo.String()+".*", repoURL.String()); matched {
				log.Debugf("Publisher found %+v", publisher)
				return publisher
			}
		}
	}

	log.Warn("Publisher not found in publishers list, slug will be generated without Id")

	return p
}
