package cmd

import (
	"fmt"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pluginsCmd)
}

var pluginsCmd = &cobra.Command{
	Use:   "clients",
	Short: "List existing clients.",
	Long:  `....`,
	Run: func(cmd *cobra.Command, args []string) {
		clients := crawler.GetClients()

		for id, _ := range clients {
			fmt.Println(id)
		}
	},
}
