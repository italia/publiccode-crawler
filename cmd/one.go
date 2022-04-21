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
	Use:   "one [repo url] whitelist.yml whitelist/*.yml",
	Short: "Crawl publiccode.yml from one single [repo url].",
	Long: `Crawl publiccode.yml from a single repository defined with [repo url] 
		according to the supplied whitelist file(s).
		No organizations! Only single repositories!`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// check if repo url is not present in blacklist
		// if so report error and exit.
		if crawler.IsRepoInBlackList(args[0]) {
			return
		}

		c := crawler.NewCrawler(dryRun)

		whitelists := args[1:]
		url, err := url.Parse(args[0])
		if err != nil {
			log.Error(err)
		}

		err = c.CrawlRepo(*url, getPAfromWhiteList(*url, whitelists))
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

func getPAfromWhiteList(repoURL url.URL, args []string) (p crawler.Publisher) {
	// Read the supplied whitelists.
	var publishers []crawler.Publisher
	for id := range args {
		readWhitelist, err := crawler.ReadAndParseWhitelist(args[id])
		if err != nil {
			log.Fatal(err)
		}
		publishers = append(publishers, readWhitelist...)
	}

	for _, paWl := range publishers {
		// looking into repositories
		for _, paWlRepo := range paWl.Repositories {
			log.Tracef("matching %s with %s", paWlRepo.String(), repoURL.String())
			if (url.URL)(paWlRepo) == repoURL {
				log.Debugf("Publisher found in whitelist %+v", paWl)
				return paWl
			}
		}
		// looking into organizations
		for _, paWlRepo := range paWl.Organizations {
			log.Tracef("matching %s.* with %s", paWlRepo.String(), repoURL.String())
			if matched, _ := regexp.MatchString(paWlRepo.String()+".*", repoURL.String()); matched {
				log.Debugf("Publisher found in whitelist %+v", paWl)
				return paWl
			}
		}
	}

	log.Warn("Publisher not found in whitelist, slug will be generated without Id")

	return p
}
