package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yo3jones/yconfig/install"
)

type installModel struct {
	inst              *install.Install
	viewport          *viewport.Model
	height            int
	width             int
	groupCount        int
	commandCount      int
	viewportMaxHeight int
}

type installMsg string

const (
	progressMsg installMsg = "progress"
	completeMsg installMsg = "complete"
)

var (
	groupLineStyle = lipgloss.NewStyle().
			PaddingTop(1)

	commandLineStyle = lipgloss.NewStyle().
				PaddingLeft(3)

	viewportStyle = lipgloss.NewStyle().
			MarginLeft(5).
			PaddingLeft(2).
			Border(lipgloss.NormalBorder(), false, false, false, true)
)

var installCmd = &cobra.Command{
	Use:   "install [GROUPS]",
	Short: "run configured install scripts",
	Long:  "run configured install scripts",
	Args:  cobra.ArbitraryArgs,
	Run: func(_ *cobra.Command, args []string) {
		program := tea.NewProgram(initModel())

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
				program.Send(progressMsg)
				program.Send(completeMsg)
				os.Exit(1)
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

func initModel() tea.Model {
	viewport := viewport.New(100, 5)

	return &installModel{
		viewport: &viewport,
	}
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
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case *install.Install:
		m.inst = msg
		m.updateCounts()
		return m, nil
	}

	return m, nil
}

func (m *installModel) updateCounts() {
	m.groupCount = 0
	m.commandCount = 0
	for i := range m.inst.Groups {
		group := &m.inst.Groups[i]
		if group.Status == install.Skipped {
			continue
		}
		m.groupCount++
		for i := range group.Commands {
			command := &group.Commands[i]
			if command.Status == install.Skipped {
				continue
			}
			m.commandCount++
		}
	}

	groupHeight := lipgloss.Height(m.groupLineView(&install.Group{Name: "foo"}))
	groupsHeight := m.groupCount * groupHeight
	commandHeight := lipgloss.Height(
		m.commandLineView(&install.Command{Command: "foo"}),
	)
	commandsHeight := m.commandCount * commandHeight

	viewportMaxHeight := m.height - groupsHeight - commandsHeight - 2

	if viewportMaxHeight > 0 {
		m.viewport.Height = viewportMaxHeight
		m.viewportMaxHeight = viewportMaxHeight
	}
}

func (m *installModel) View() string {
	if m.inst == nil {
		return ""
	}
	return m.groupsView(m.inst.Groups)
}

func (m *installModel) groupsView(groups []install.Group) string {
	sb := &strings.Builder{}
	for i := range groups {
		group := &groups[i]
		m.groupView(sb, group)
	}
	return sb.String()
}

func (m *installModel) groupView(
	sb *strings.Builder,
	group *install.Group,
) string {
	if group.Status == install.Skipped {
		return ""
	}

	fmt.Fprintf(sb, "%s\n", m.groupLineView(group))
	m.commandsView(sb, group.Commands)

	return sb.String()
}

func (m *installModel) groupLineView(group *install.Group) string {
	return groupLineStyle.Render(group.Name)
}

func (m *installModel) commandsView(
	sb *strings.Builder,
	commands []install.Command,
) {
	for i := range commands {
		command := &commands[i]
		m.commandView(sb, command)
	}
}

func (m *installModel) commandView(
	sb *strings.Builder,
	command *install.Command,
) {
	if command.Status == install.Skipped {
		return
	}

	fmt.Fprintf(
		sb,
		"%s\n",
		m.commandLineView(command),
	)

	if command.Status != install.Error && command.Status != install.Running {
		return
	}

	outLineCount := getOutLineCount(command.Out)
	if outLineCount < m.viewportMaxHeight {
		m.viewport.Height = outLineCount
	} else {
		m.viewport.Height = m.viewportMaxHeight
	}

	m.viewport.SetContent(string(getContentForViewPort(command.Out)))
	m.viewport.GotoBottom()

	fmt.Fprintf(
		sb,
		"%s\n",
		viewportStyle.Render(m.viewport.View()),
	)
}

func getContentForViewPort(b []byte) []byte {
	bLen := len(b)

	if bLen > 0 && b[bLen-1] == '\n' {
		return b[:bLen-1]
	}

	return b
}

func getOutLineCount(out []byte) int {
	count := 0
	for _, b := range out {
		if b == '\n' {
			count++
		}
	}
	return count
}

func (m *installModel) commandLineView(command *install.Command) string {
	return commandLineStyle.Render(
		fmt.Sprintf(
			"[%s] %s",
			command.Status,
			truncCommand(command.Command),
		),
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
