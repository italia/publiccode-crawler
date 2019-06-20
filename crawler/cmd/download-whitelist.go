package cmd

import (
	"os"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"github.com/spf13/cobra"
	"log"
	"github.com/italia/developers-italia-backend/crawler/crawler"
	"github.com/thoas/go-funk"
)

func init() {
	rootCmd.AddCommand(downloadWhitelistCmd)
}

type repolistType struct {
	Registrati []struct {
		IPA string `yaml:"ipa"`
		URL string `yaml:"url"`
		PEC string `yaml:"pec"`
	} `yaml:"registrati"`
}

var downloadWhitelistCmd = &cobra.Command{
	Use:   "download-whitelist REPOLIST_URL DEST_FILE",
	Short: "Download the list of repos and orgs from the onboarding portal.",
	Long:  `Download the list of repos and orgs from the onboarding portal and convert it into a yml whitelist file.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Read the current destinatin whitelist, if any
		var publishers crawler.Whitelist
		if _, err := os.Stat(args[1]); err == nil {
			data, err := ioutil.ReadFile(args[1])
			if err != nil {
				log.Fatalf("error in reading %s: %v", args[1], err)
			}
			yaml.Unmarshal(data, &publishers)
		}

		// Download the repo-list file
		resp, err := http.Get(args[0])
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
	
		// Parse the repo-list file
		var repolist repolistType
		err = yaml.Unmarshal(bodyBytes, &repolist)
		if err != nil {
			log.Fatal(err)
		}

		// Merge the repo-list file into the whitelist
		REPOLIST:
		for _, i := range repolist.Registrati {
			for _, publisher := range publishers {
				if publisher.CodiceIPA == i.IPA {
					// If this IPA code is already known, append this URL to the existing item
					publisher.Organizations = funk.UniqString(append(publisher.Organizations, i.URL))
					continue REPOLIST
				}
			}

			// If this IPA code is not known, append a new publisher item
			publishers = append(publishers, crawler.PA{
				Name: i.IPA,
				CodiceIPA: i.IPA,
				Organizations: []string{i.URL},
			})
		}

		// Write to the destination file
		f, err := os.Create(args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		data, err := yaml.Marshal(publishers)
		if err != nil {
			log.Fatal(err)
		}
		f.Write(data)
	}}
