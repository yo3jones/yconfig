package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yo3jones/yconfig/setup"
)

var setupCmd = &cobra.Command{
	Use:   "setup [Name]",
	Short: "execute setup steps",
	Long:  "execute setup steps",
	Run: func(_ *cobra.Command, _ []string) {
		var err error

		scriptsConfig := viper.Get("scripts")
		packageManagersConfig := viper.Get("packageManagers")
		config := viper.Get("setup")

		if err = setup.New().
			ScriptsConfig(&scriptsConfig).
			PackageManagersConfig(&packageManagersConfig).
			Config(&config).
			Setup(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
