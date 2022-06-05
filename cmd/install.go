package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yo3jones/yconfig/install"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "run configured install scripts",
	Long:  "run configured install scripts",
	Run: func(_ *cobra.Command, _ []string) {
		config := install.Parse(viper.Get("install"))
		install.Print(config)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
