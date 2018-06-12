package cmd

import (
	"os"
	"strconv"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(clientsCmd)
}

var clientsCmd = &cobra.Command{
	Use:   "clients",
	Short: "List existing clients.",
	Long:  `List existing clients registered in the crawler.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve crawler clients.
		clients := crawler.GetClients()

		// Prepare data table.
		var data [][]string

		// Iterate over the crawler clients.
		for id := range clients {
			data = append(data, []string{id})
		}

		// Write data and render as table in os.Stdout.
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Client ID"})
		table.SetFooter([]string{"Total Client IDs: " + strconv.Itoa(len(clients))})
		table.SetAutoMergeCells(true)
		table.SetRowLine(true)
		table.AppendBulk(data)
		table.Render()

	},
}
