package cmd

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yo3jones/yconfig/install"
)

var installCmd = &cobra.Command{
	Use:   "install [GROUPS] [OPTIONS]",
	Short: "run configured install scripts",
	Long:  "run configured install scripts",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			all  bool
			tags []string
			err  error
		)

		if all, err = cmd.Flags().GetBool("all"); err != nil {
			panic(err)
		}

		installError := false
		program := tea.NewProgram(install.InitModel())

		go func() {
			config := viper.Get("install")
			installer, err := install.New(&config)
			if err != nil {
				panic(err)
			}

			err = installer.
				Groups(args).
				All(all).
				Tags(tags).
				OnProgress(func(inst *install.Install) {
					program.Send(inst)
				}).
				Install()
			if err != nil {
				installError = true
				program.Send(install.ProgressMsg)
				program.Send(install.ErrorMsg)
			} else {
				program.Send(install.ProgressMsg)
				program.Send(install.CompleteMsg)
			}
		}()

		if err := program.Start(); err != nil {
			panic(err)
		}

		if installError {
			os.Exit(1)
		}
	},
}

func init() {
	installCmd.Flags().Bool("all", false, "set to insall all groups")
	installCmd.Flags().StringSlice(
		"tag",
		[]string{},
		"tags to filter groups and commands to install",
	)

	rootCmd.AddCommand(installCmd)
}
