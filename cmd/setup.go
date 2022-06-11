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
		var (
			s   *setup.Setup
			err error
		)

		config := viper.Get("setup")
		if s, err = setup.Parse(&config); err != nil {
			panic(err)
		}

		s.Print()
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
