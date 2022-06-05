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
		var err error

		config := viper.Get("install")

		err = install.New(config).
			Install()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
