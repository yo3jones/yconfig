package cmd

import (
	"fmt"

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
			s       *setup.Setup
			pms     []setup.PackageManager
			scripts []setup.Script
			err     error
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

		fmt.Println("packageManagers")
		for _, pm := range pms {
			pm.Print()
		}

		scriptsConfig := viper.Get("scripts")
		if scripts, err = setup.ParseScripts(&scriptsConfig); err != nil {
			panic(err)
		}

		fmt.Println("scripts")
		for _, script := range scripts {
			script.Print()
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
