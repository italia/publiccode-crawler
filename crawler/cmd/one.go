package cmd

import (
	"regexp"

	"github.com/italia/developers-italia-backend/crawler/crawler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
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

		c := crawler.NewCrawler()

		repoURL, whitelists := args[0], args[1:]
		err := c.CrawlRepo(repoURL, getPAfromWhiteList(repoURL, whitelists))
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

func getPAfromWhiteList(repoURL string, args []string) (pa crawler.PA) {
	// Read the supplied whitelists.
	var publishers []crawler.PA
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
			log.Debugf("matching %s with %s", paWlRepo, repoURL)
			if paWlRepo == repoURL {
				log.Debugf("PA found in whitelist %+v", paWl)
				return paWl
			}
		}
		// looking into organizations
		for _, paWlRepo := range paWl.Organizations {
			log.Debugf("matching %s.* with %s", repoURL, paWlRepo)
			if matched, _ := regexp.MatchString(paWlRepo+".*", repoURL); matched {
				log.Debugf("PA found in whitelist %+v", paWl)
				return paWl
			}
		}
	}

	log.Warn("PA not found in whitelist, slug will be generated without coideIPA")
	// since this routine is called by command: `<command_name> one ...`
	// that is not aware about whitelists
	// this hack will skip IPA code match with those lists
	pa.UnknownIPA = true
	return pa
}
