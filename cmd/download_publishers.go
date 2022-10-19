package cmd

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/italia/publiccode-crawler/v3/common"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	ymlurl "github.com/italia/publiccode-crawler/v3/internal"
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
	Run: func(cmd *cobra.Command, args []string) {
		var publishers []common.Publisher
		if _, err := os.Stat(args[1]); err == nil {
			data, err := ioutil.ReadFile(args[1])
			if err != nil {
				log.Fatalf("error in reading %s: %v", args[1], err)
			}
			_ = yaml.Unmarshal(data, &publishers)
		}

		resp, err := http.Get(args[0])
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var repolist repolistType
		err = yaml.Unmarshal(bodyBytes, &repolist)
		if err != nil {
			log.Fatal(err)
		}

	REPOLIST:
		for _, i := range repolist.Registrati {
			for idx, publisher := range publishers {
				if publisher.Id == i.IPA {
					u, _ := url.Parse(i.URL)
					// If this Id is already known, append this URL to the existing item
					publishers[idx].Organizations = append(publisher.Organizations, (ymlurl.URL)(*u))
					continue REPOLIST
				}
			}

			u, _ := url.Parse(i.URL)
			// If this IPA code is not known, append a new publisher item
			publishers = append(publishers, common.Publisher{
				Name:          i.IPA,
				Id:            i.IPA,
				Organizations: []ymlurl.URL{(ymlurl.URL)(*u)},
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
		if _, err = f.Write(data); err != nil {
			log.Fatal(err)
		}
	}}
