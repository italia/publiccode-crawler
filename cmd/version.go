package cmd

import (
	"fmt"

	"github.com/italia/developers-italia-backend/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of the crawler.",
	Long:  `All software has versions. This too.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version " + version.VERSION)
	},
}
