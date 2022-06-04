package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yo3jones/yconfig/generate"
)

const (
	nameTemplateRoot string = "template-root"
	pathTempalteRoot string = "generate.templateRoot"

	nameDestRoot string = "destination-root"
	pathDestRoot string = "generate.destinationRoot"

	nameInclude string = "include"
	pathInclude string = "generate.include"

	nameExclude string = "exclude"
	pathExclude string = "generate.exclude"
)

var genCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate config files from templates",
	Long:  "generate config files from templates",
	RunE: func(_ *cobra.Command, _ []string) error {
		templateRoot := viper.GetString(pathTempalteRoot)
		destinationRoot := viper.GetString(pathDestRoot)
		include := viper.GetStringSlice(pathInclude)
		exclude := viper.GetStringSlice(pathExclude)

		return generate.New().
			TemplateRoot(templateRoot).
			DesinationRoot(destinationRoot).
			Include(include).
			Exclude(exclude).
			OnProgress(func(progress *generate.Progress) {
				for _, t := range progress.TemplatesProgress {
					fmt.Printf("%s : %s\n", t.Progress, t.Path)
				}
				fmt.Println()
			}).
			Generate()
	},
}

func init() {
	genCmd.Flags().
		String(nameTemplateRoot, ".", "root path for the template files")
	err1 := viper.BindPFlag(
		pathTempalteRoot,
		genCmd.Flags().Lookup(nameTemplateRoot),
	)
	if err1 != nil {
		fmt.Fprintln(os.Stderr, err1)
	}

	genCmd.Flags().
		String(nameDestRoot, ".", "root path for writing the output of the templates")
	err2 := viper.BindPFlag(
		pathDestRoot,
		genCmd.Flags().Lookup(nameDestRoot),
	)
	if err2 != nil {
		fmt.Fprintln(os.Stderr, err2)
	}

	genCmd.Flags().StringSlice(
		nameInclude,
		[]string{},
		"globs of config file templates to include",
	)
	err3 := viper.BindPFlag(pathInclude, genCmd.Flags().Lookup(nameInclude))
	if err3 != nil {
		fmt.Fprintln(os.Stderr, err3)
	}

	genCmd.Flags().StringSlice(
		nameExclude,
		[]string{},
		"globs of config file templates to exclude",
	)
	err4 := viper.BindPFlag(pathExclude, genCmd.Flags().Lookup(nameExclude))
	if err4 != nil {
		fmt.Fprintln(os.Stderr, err4)
	}

	rootCmd.AddCommand(genCmd)
}
