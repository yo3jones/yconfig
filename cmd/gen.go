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
	nameDestRoot     string = "destination-root"
	pathDestRoot     string = "generate.destinationRoot"
	nameInclude      string = "include"
	pathInclude      string = "generate.include"
	nameExclude      string = "exclude"
	pathExclude      string = "generate.exclude"
	nameLink         string = "link"
	pathLink         string = "generate.link"
	nameDelay        string = "delay"
	pathDelay        string = "generate.delay"
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
		link := viper.GetBool(pathLink)
		delay := viper.GetInt(pathDelay)

		program := tea.NewProgram(model{})

		go func() {
			err := generate.New().
				TemplateRoot(templateRoot).
				DesinationRoot(destinationRoot).
				Include(include).
				Exclude(exclude).
				Link(link).
				Delay(delay).
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
	var err error

	genCmd.Flags().
		String(nameTemplateRoot, ".", "root path for the template files")
	err = viper.BindPFlag(
		pathTempalteRoot,
		genCmd.Flags().Lookup(nameTemplateRoot),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	genCmd.Flags().
		String(nameDestRoot, ".", "root path for writing the output of the templates")
	err = viper.BindPFlag(
		pathDestRoot,
		genCmd.Flags().Lookup(nameDestRoot),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	genCmd.Flags().StringSlice(
		nameInclude,
		[]string{},
		"globs of config file templates to include",
	)
	err = viper.BindPFlag(pathInclude, genCmd.Flags().Lookup(nameInclude))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	genCmd.Flags().StringSlice(
		nameExclude,
		[]string{},
		"globs of config file templates to exclude",
	)
	err = viper.BindPFlag(pathExclude, genCmd.Flags().Lookup(nameExclude))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	genCmd.Flags().Bool(
		nameLink,
		true,
		"whether to link the generated config files to the users home dir",
	)
	err = viper.BindPFlag(pathLink, genCmd.Flags().Lookup(nameLink))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	genCmd.Flags().
		Int(nameDelay, 0, "add a delay in mils to see cool animations")
	err = viper.BindPFlag(pathDelay, genCmd.Flags().Lookup(nameDelay))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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

	maxPathWidth := 0
	for _, p := range m.progress.TemplatesProgress {
		maxPathWidth = maxWidth(maxPathWidth, lipgloss.Width(p.Path))
	}

	sb := strings.Builder{}
	bracketStyle := lipgloss.NewStyle().
		Bold(true)

	for _, p := range m.progress.TemplatesProgress {
		var symbol string
		symbolStyle := lipgloss.NewStyle().Bold(true)
		switch p.Status {
		case generate.Waiting:
			symbolStyle.Foreground(lipgloss.Color("13"))
			symbol = ""
		case generate.Generating, generate.Linking:
			symbolStyle.Foreground(lipgloss.Color("12"))
			symbol = "◯"
		case generate.Complete:
			symbolStyle.Foreground(lipgloss.Color("10"))
			symbol = "✔"
		case generate.Error:
			symbolStyle.Foreground(lipgloss.Color("9"))
			symbol = "⚠"
		default:
			symbol = ""
		}
		sb.WriteString(
			fmt.Sprintf(
				"   %s %s %s %s   %s\n",
				bracketStyle.Render("["),
				symbolStyle.Copy().
					Width(10).
					Align(lipgloss.Right).
					Render(p.Status.String()),
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
