package cmd

import (
	"fmt"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the Public Administrations in the whitelist.",
	Long:  `List all the Public Administrations in whitelist.yml file.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Read and parse the whitelist.
		whitelistFile := "whitelist.yml"
		whitelist, err := crawler.ReadAndParseWhitelist(whitelistFile)
		if err != nil {
			panic(err)
		}

		// Process every item in whitelist
		for _, pa := range whitelist {
			// and Print
			fmt.Printf("%s \t %s\n", pa.Name, pa.CodiceIPA)
			for _, repository := range pa.Repositories {
				fmt.Printf("\t%s\n", repository.API)
				for _, orgs := range repository.Organizations {
					fmt.Printf("\t\t- %s\n", orgs)
				}
			}
		}

	}}
