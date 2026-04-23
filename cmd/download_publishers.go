package cmd

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/italia/publiccode-crawler/v4/common"
	ymlurl "github.com/italia/publiccode-crawler/v4/internal"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	rootCmd.AddCommand(downloadPublishersCmd)
}

type repolistType struct {
	Registrati []struct {
		IPA string `yaml:"ipa"`
		URL string `yaml:"url"`
		PEC string `yaml:"pec"`
	} `yaml:"registrati"`
}

var downloadPublishersCmd = &cobra.Command{
	Use:   "download-publishers REPOLIST_URL DEST_FILE",
	Short: "Download the list of repos and orgs from the onboarding portal.",
	Long:  `Download the list of repos and orgs from the onboarding portal and convert it into a publishers.yml.`,
	Args:  cobra.ExactArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		var publishers []common.Publisher
		if _, err := os.Stat(args[1]); err == nil {
			data, err := os.ReadFile(args[1])
			if err != nil {
				log.Fatalf("error in reading %s: %v", args[1], err)
			}
			//nolint:musttag // false positive
			_ = yaml.Unmarshal(data, &publishers)
		}

		resp, err := http.Get(args[0])
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var repolist repolistType
		err = yaml.Unmarshal(bodyBytes, &repolist)
		if err != nil {
			log.Fatal(err)
		}

	REPOLIST:
		for _, entry := range repolist.Registrati {
			for idx, publisher := range publishers {
				if publisher.ID == entry.IPA {
					parsedURL, _ := url.Parse(entry.URL)
					// If this Id is already known, append this URL to the existing item
					publishers[idx].Organizations = append(publisher.Organizations, (ymlurl.URL)(*parsedURL))

					continue REPOLIST
				}
			}

			parsedURL, _ := url.Parse(entry.URL)
			// If this IPA code is not known, append a new publisher item
			publishers = append(publishers, common.Publisher{
				Name:          entry.IPA,
				ID:            entry.IPA,
				Organizations: []ymlurl.URL{(ymlurl.URL)(*parsedURL)},
			})
		}

		// Write to the destination file
		outFile, err := os.Create(args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()
		data, err := yaml.Marshal(publishers)
		if err != nil {
			log.Fatal(err)
		}
		if _, err = outFile.Write(data); err != nil {
			log.Fatal(err)
		}
	},
}
