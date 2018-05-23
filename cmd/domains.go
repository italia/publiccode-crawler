package cmd

import (
	"os"
	"strconv"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(domainsCmd)
}

var domainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "List all the Domains.",
	Long:  `List all the Domains from domains.yml`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read and parse the whitelist.
		domainsFile := "domains.yml"
		domains, err := crawler.ReadAndParseDomains(domainsFile)
		if err != nil {
			panic(err)
		}
		// Prepare data table.
		var data [][]string

		// Process every item in domains.
		for _, domain := range domains {
			// Check if basicAuth is set.
			basicAuth := "no"
			if len(domain.BasicAuth) > 0 {
				basicAuth = "yes"
			}
			data = append(data, []string{domain.Id, domain.Description, domain.ClientApi, basicAuth})

		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Description", "API", "BasicAuth?"})
		table.SetFooter([]string{"Total Domains: " + strconv.Itoa(len(domains)), "", "", ""})
		table.SetRowLine(true)
		table.AppendBulk(data)
		table.Render()

	}}
