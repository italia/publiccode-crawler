package cmd

import (
	"fmt"

	"github.com/italia/publiccode-crawler/v3/internal"
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
		//nolint:forbidigo
		fmt.Println("Version:\t", internal.VERSION)
		//nolint:forbidigo
		fmt.Println("Build time:\t", internal.BuildTime)
	},
}
