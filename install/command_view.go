package install

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type commandModel struct {
	command           *Command
	viewport          *viewport.Model
	viewportMaxHeight int
}

var (
	commandLineStyle = lipgloss.NewStyle().
				PaddingLeft(2)

	bracketStyle = lipgloss.NewStyle().
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Width(10).
			Align(lipgloss.Right).
			Foreground(lipgloss.Color("10"))

	viewportStyle = lipgloss.NewStyle().
			MarginLeft(5).
			PaddingLeft(2).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("12"))
)

func NewCommandModel(command *Command) *commandModel {
	viewport := viewport.New(100, 5)

	return &commandModel{
		command:  command,
		viewport: &viewport,
	}
}

func (m *commandModel) GetCollapsedHeight() int {
	return lipgloss.Height(
		m.commandLineView(&Command{Command: "foo"}),
	)
}

func (m *commandModel) SetViewportMaxHeight(viewportMaxHeight int) {
	m.viewportMaxHeight = viewportMaxHeight
	m.viewport.Height = viewportMaxHeight
}

func (*commandModel) Init() tea.Cmd {
	return nil
}

func (m *commandModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *commandModel) View() string {
	sb := &strings.Builder{}

	fmt.Fprintf(
		sb,
		"%s\n",
		m.commandLineView(m.command),
	)

	if m.command.Status != Error && m.command.Status != Running {
		return sb.String()
	}

	outLineCount := getOutLineCount(m.command.Out)
	if outLineCount < m.viewportMaxHeight {
		m.viewport.Height = outLineCount
	} else {
		m.viewport.Height = m.viewportMaxHeight
	}

	m.viewport.SetContent(string(getContentForViewPort(m.command.Out)))
	m.viewport.GotoBottom()

	fmt.Fprintf(
		sb,
		"%s\n",
		viewportStyle.Render(m.viewport.View()),
	)

	return sb.String()
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

var commandLineTemplate string = fmt.Sprintf(
	"%s%%s%s %%s",
	bracketStyle.Render("["),
	bracketStyle.Render("]"),
)

func (m *commandModel) commandLineView(command *Command) string {
	return commandLineStyle.Render(
		fmt.Sprintf(
			commandLineTemplate,
			getStatusStyle(command).Render(command.Status.String()),
			truncCommand(command.Command),
		),
	)
}

func getStatusStyle(command *Command) *lipgloss.Style {
	style := statusStyle.Copy()

	switch command.Status {
	case Skipped:
	case Waiting:
		style.Foreground(lipgloss.Color("13"))
	case Running:
		style.Foreground(lipgloss.Color("12"))
	case Complete:
		style.Foreground(lipgloss.Color("10"))
	case Error:
		style.Foreground(lipgloss.Color("9"))
	}

	return &style
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
