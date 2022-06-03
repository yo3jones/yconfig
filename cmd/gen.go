package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yo3jones/yconfig/generate"
)

const (
	nameTemplateRoot string = "template-root"
	pathTempalteRoot string = "generate.templateRoot"

	nameInclude string = "include"
	pathInclude string = "generate.include"

	nameExclude string = "exclude"
	pathExclude string = "generate.exclude"
)

var genCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate config files from templates",
	Long:  "generate config files from templates",
	Run: func(_ *cobra.Command, _ []string) {
		templateRoot := viper.GetString(pathTempalteRoot)
		include := viper.GetStringSlice(pathInclude)
		exclude := viper.GetStringSlice(pathExclude)

		generate.New().
			TemplateRoot(templateRoot).
			Include(include).
			Exclude(exclude).
			Progress(func(progress *generate.Progress) {
				for _, t := range progress.TemplatesProgress {
					fmt.Printf("%s : %s\n", t.Progress, t.Path)
				}
			}).
			Generate()
	},
}

func init() {
	genCmd.Flags().
		String(nameTemplateRoot, ".", "root path for the template files")
	viper.BindPFlag(pathTempalteRoot, genCmd.Flags().Lookup(nameTemplateRoot))

	genCmd.Flags().StringSlice(
		nameInclude,
		[]string{},
		"globs of config file templates to include",
	)
	viper.BindPFlag(pathInclude, genCmd.Flags().Lookup(nameInclude))

	genCmd.Flags().StringSlice(
		nameExclude,
		[]string{},
		"globs of config file templates to exclude",
	)
	viper.BindPFlag(pathExclude, genCmd.Flags().Lookup(nameExclude))

	rootCmd.AddCommand(genCmd)
}
