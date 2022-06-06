package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yo3jones/yconfig/install"
)

var installCmd = &cobra.Command{
	Use:   "install [GROUPS]",
	Short: "run configured install scripts",
	Long:  "run configured install scripts",
	Args:  cobra.ArbitraryArgs,
	Run: func(_ *cobra.Command, args []string) {
		var err error

		config := viper.Get("install")

		installer, err := install.New(&config)
		if err != nil {
			panic(err)
		}

		err = installer.
			Groups(args).
			OnProgress(func(_ *install.Install) {
				fmt.Println("progress")
			}).
			Install()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
