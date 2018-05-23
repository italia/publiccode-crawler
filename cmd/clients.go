package cmd

import (
	"os"
	"strconv"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pluginsCmd)
}

var pluginsCmd = &cobra.Command{
	Use:   "clients",
	Short: "List existing clients.",
	Long:  `List existing clients registered in the crawler.`,
	Run: func(cmd *cobra.Command, args []string) {
		clients := crawler.GetClients()

		var data [][]string

		for id, _ := range clients {
			data = append(data, []string{id})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Client ID"})
		table.SetFooter([]string{"Total Client IDs: " + strconv.Itoa(len(clients))})
		table.SetAutoMergeCells(true)
		table.SetRowLine(true)
		table.AppendBulk(data)
		table.Render()

	},
}
