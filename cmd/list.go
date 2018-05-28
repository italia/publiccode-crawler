package cmd

import (
	"os"
	"strconv"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/olekukonko/tablewriter"
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
		whitelist, err := crawler.ReadAndParseWhitelist(whitelistFile)
		if err != nil {
			panic(err)
		}

		// Prepare data table.
		var data [][]string

		// Process every item in whitelist.
		for _, pa := range whitelist {
			// And add to data table.
			data = append(data, []string{pa.CodiceIPA, pa.Name, "", ""})
			for _, repository := range pa.Repositories {
				data = append(data, []string{pa.CodiceIPA, pa.Name, repository.API, ""})
				for _, org := range repository.Organizations {
					data = append(data, []string{pa.CodiceIPA, pa.Name, repository.API, org})
				}
			}
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"iPA", "Description", "API", "Org"})
		table.SetFooter([]string{"Total Public Administrations: " + strconv.Itoa(len(whitelist)), "", "", ""})
		table.SetAutoMergeCells(true)
		table.SetRowLine(true)
		table.AppendBulk(data)
		table.Render()

	}}
