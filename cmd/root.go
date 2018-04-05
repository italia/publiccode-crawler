package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "crawler",
	Short: "Crawler is a crawler for publiccode.yml file.",
	Long: `A Fast and Robust publiccode.yml file crawler.
        Complete documentation is available at https://github.com/italia/developers-italia-backend`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute is the entrypoint for cmd package Cobra.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
