package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var dryRun bool
var rootCmd = &cobra.Command{
	Use:   "crawler",
	Short: "A crawler for publiccode.yml files.",
	Long: `A fast and robust publiccode.yml file crawler.
Complete documentation is available at https://github.com/italia/developers-italia-backend`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			log.Fatal(err)
		}
	},
}

// Execute is the entrypoint for cmd package Cobra.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
