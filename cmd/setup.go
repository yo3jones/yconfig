package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yo3jones/yconfig/setup"
)

var tags []string

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
			Tags(tags).
			OnProgress(func(progress []*setup.Progress) {
				for _, p := range progress {
					fmt.Printf("%s %s %d\n", p.Value.GetName(), p.Status, len(p.Out))
				}
				fmt.Println()
			}).
			Setup(); err != nil {
			panic(err)
		}
	},
}

func init() {
	setupCmd.Flags().
		StringSliceVar(&tags, "tag", []string{}, "tags used for filtering")

	rootCmd.AddCommand(setupCmd)
}
