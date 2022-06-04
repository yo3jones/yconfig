package cmd

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

type model struct {
	progress *generate.Progress
}

type ProgressMsg struct {
	progress *generate.Progress
}

type doneMsg string

var genCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate config files from templates",
	Long:  "generate config files from templates",
	Run: func(_ *cobra.Command, _ []string) {
		templateRoot := viper.GetString(pathTempalteRoot)
		destinationRoot := viper.GetString(pathDestRoot)
		include := viper.GetStringSlice(pathInclude)
		exclude := viper.GetStringSlice(pathExclude)

		program := tea.NewProgram(model{})

		go func() {
			err := generate.New().
				TemplateRoot(templateRoot).
				DesinationRoot(destinationRoot).
				Include(include).
				Exclude(exclude).
				OnProgress(func(progress *generate.Progress) {
					program.Send(ProgressMsg{progress})
				}).
				Generate()
			if err != nil {
				panic(err)
			}

			var dm doneMsg = "done"
			program.Send(dm)
		}()

		if err := program.Start(); err != nil {
			panic(err)
		}
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

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case ProgressMsg:
		m.progress = msg.progress
		return m, nil
	case doneMsg:
		return m, tea.Quit
	}
	return m, nil
}

func maxWidth(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (m model) View() string {
	if m.progress == nil {
		return "starting"
	}

	var maxPathWidth int = 0
	for _, p := range m.progress.TemplatesProgress {
		maxPathWidth = maxWidth(maxPathWidth, lipgloss.Width(p.Path))
	}

	sb := strings.Builder{}
	bracketStyle := lipgloss.NewStyle().
		Bold(true)

	for _, p := range m.progress.TemplatesProgress {
		var symbol string
		symbolStyle := lipgloss.NewStyle().Bold(true)
		switch p.Progress {
		case "waiting":
			symbolStyle.Foreground(lipgloss.Color("13"))
			symbol = ""
		case "generating":
			symbolStyle.Foreground(lipgloss.Color("12"))
			symbol = "◯"
		case "complete":
			symbolStyle.Foreground(lipgloss.Color("10"))
			symbol = "✔"
		case "error":
			symbolStyle.Foreground(lipgloss.Color("9"))
			symbol = "⚠"
		default:
			symbol = ""
		}
		sb.WriteString(
			fmt.Sprintf(
				"   %s %s %s %s %s\n",
				bracketStyle.Render("["),
				symbolStyle.Copy().
					Width(10).
					Align(lipgloss.Right).
					Render(p.Progress),
				bracketStyle.Render("]"),
				lipgloss.NewStyle().
					Width(maxPathWidth).
					Render(p.Path),
				symbolStyle.Render(symbol),
			),
		)
	}

	return sb.String()
}
