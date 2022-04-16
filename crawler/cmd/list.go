package cmd

import (
	"log"
	"os"
	"strconv"

	"github.com/italia/developers-italia-backend/crawler/crawler"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list whitelist.yml",
	Short: "List all the PA in the whitelist file.",
	Long:  `List all the Public Administrations in the whitelist file.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Read and parse the whitelist.
		whitelist, err := crawler.ReadAndParseWhitelist(args[0])
		if err != nil {
			log.Fatal(err)
		}

		// Prepare data table.
		var data [][]string

		// Process every item in whitelist.
		for _, publisher := range whitelist {
			// And add to data table.
			data = append(data, []string{publisher.Name, publisher.Id, ""})
			for _, org := range publisher.Organizations {
				data = append(data, []string{publisher.Name, publisher.Id, org.String()})
			}
		}

		// Write data and render as table in os.Stdout.
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Codice iPA", "Repository"})
		table.SetFooter([]string{"Total Public Administrations: " + strconv.Itoa(len(whitelist)), "", ""})
		table.SetAutoMergeCells(true)
		table.SetRowLine(true)
		table.AppendBulk(data)
		table.Render()

	}}
