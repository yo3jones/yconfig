package cmd

import (
	"fmt"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yo3jones/yconfig/install"
)

type installModel struct {
	inst *install.Install
}

type installMsg string

const (
	progressMsg installMsg = "progress"
	completeMsg installMsg = "complete"
)

var installCmd = &cobra.Command{
	Use:   "install [GROUPS]",
	Short: "run configured install scripts",
	Long:  "run configured install scripts",
	Args:  cobra.ArbitraryArgs,
	Run: func(_ *cobra.Command, args []string) {
		program := tea.NewProgram(&installModel{})

		go func() {
			config := viper.Get("install")
			installer, err := install.New(&config)
			if err != nil {
				panic(err)
			}

			err = installer.
				Groups(args).
				OnProgress(func(inst *install.Install) {
					program.Send(inst)
				}).
				Install()
			if err != nil {
				panic(err)
			}

			program.Send(progressMsg)
			program.Send(completeMsg)
		}()

		if err := program.Start(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func (m *installModel) Init() tea.Cmd {
	return nil
}

func (m *installModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case installMsg:
		if msg == completeMsg {
			return m, tea.Quit
		}
	case *install.Install:
		m.inst = msg
		return m, nil
	}

	return m, nil
}

func (m *installModel) View() string {
	if m.inst == nil {
		return ""
	}
	return groupsView(m.inst.Groups)
}

func groupsView(groups []install.Group) string {
	sb := &strings.Builder{}
	for i := range groups {
		group := &groups[i]
		groupView(sb, group)
	}
	return sb.String()
}

func groupView(sb *strings.Builder, group *install.Group) string {
	if group.Status == install.Skipped {
		return ""
	}

	fmt.Fprintf(sb, "\n%s\n", group.Name)
	commandsView(sb, group.Commands)

	return sb.String()
}

func commandsView(sb *strings.Builder, commands []install.Command) {
	for i := range commands {
		command := &commands[i]
		commandView(sb, command)
	}
}

func commandView(sb *strings.Builder, command *install.Command) {
	if command.Status == install.Skipped {
		return
	}

	fmt.Fprintf(
		sb,
		"  [%s] %s\n",
		command.Status,
		truncCommand(command.Command),
	)

	if command.Status != install.Error && command.Status != install.Running {
		return
	}

	fmt.Fprintf(
		sb,
		"%s\n",
		lipgloss.NewStyle().
			MarginLeft(6).
			Render(string(command.Out)),
	)
}

var (
	m1 *regexp.Regexp = regexp.MustCompile("\\\\\n")
	m2                = regexp.MustCompile("[[:space:]]+")
)

func truncCommand(command string) string {
	str := command
	str = m1.ReplaceAllString(str, " ")
	str = m2.ReplaceAllString(str, " ")

	// TODO trim long commands

	return str
}
