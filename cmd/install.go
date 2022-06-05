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
		var (
			inst *install.Install
			err  error
		)

		if inst, err = install.Parse(viper.Get("install")); err != nil {
			panic(err)
		}

		install.Print(inst)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
