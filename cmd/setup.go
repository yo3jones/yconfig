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
			pms []setup.PackageManager
			err error
		)

		config := viper.Get("setup")
		if s, err = setup.Parse(&config); err != nil {
			panic(err)
		}

		s.Print()

		pmsConfig := viper.Get("packageManagers")
		if pms, err = setup.ParsePackageManagers(&pmsConfig); err != nil {
			panic(err)
		}

		for _, pm := range pms {
			pm.Print()
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
