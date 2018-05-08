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
	Use:   "plugins",
	Short: "List existing plugins.",
	Long:  `....`,
	Run: func(cmd *cobra.Command, args []string) {
		// Register client API plugins.
		crawler.RegisterCrawlers()

		plugins := crawler.GetPlugins()

		for id, _ := range plugins {
			fmt.Println(id)
		}
	},
}
