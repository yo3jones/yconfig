package cmd

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yo3jones/yconfig/setup"
)

var (
	tags  []string
	delay int
)

var setupCmd = &cobra.Command{
	Use:   "setup [Name]",
	Short: "execute setup steps",
	Long:  "execute setup steps",
	Run: func(_ *cobra.Command, _ []string) {
		run()
	},
}

func init() {
	setupCmd.Flags().
		StringSliceVar(&tags, "tag", []string{}, "tags used for filtering")
	setupCmd.Flags().
		IntVar(&delay, "delay", 0, "add delay between setup entries")

	rootCmd.AddCommand(setupCmd)
}

func run() {
	var setupErr error

	program := tea.NewProgram(setup.InitModel())

	go func() {
		var err error

		scriptsConfig := viper.Get("scripts")
		packageManagersConfig := viper.Get("packageManagers")
		config := viper.Get("setup")

		err = setup.New().
			ScriptsConfig(&scriptsConfig).
			PackageManagersConfig(&packageManagersConfig).
			Config(&config).
			Tags(tags).
			Delay(delay).
			OnProgress(func(progress []*setup.Progress) {
				program.Send(progress)
			}).
			Setup()
		if err != nil {
			setupErr = err
		}

		program.Send(setup.MsgRefresh)
		program.Send(setup.MsgDone)
	}()

	if err := program.Start(); err != nil {
		panic(err)
	}

	if setupErr != nil {
		os.Exit(1)
	}
}
