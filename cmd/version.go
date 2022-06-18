package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yo3jones/yconfig/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of yconfig",
	Long:  "Print the version of yconfig",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("yconfig\nversion: %s\n", version.Get())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
