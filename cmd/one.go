package cmd

import (
	"net/url"
	"regexp"

	"github.com/italia/publiccode-crawler/v3/common"
	"github.com/italia/publiccode-crawler/v3/crawler"
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
	},
}

func getPublisher(repoURL url.URL, paths []string) (p common.Publisher) {
	var publishers []common.Publisher
	for _, path := range paths {
		p, err := common.LoadPublishers(path)
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
