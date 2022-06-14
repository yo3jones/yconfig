package setup

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	statusLineStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	statusBracketStyle = lipgloss.NewStyle().
				Bold(true)

	statusStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Width(10).
			Align(lipgloss.Right)

	viewportStyle = lipgloss.NewStyle().
			MarginLeft(5).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			PaddingLeft(2).
			BorderForeground(lipgloss.Color("12"))
)

type ValueModel struct {
	progress       *Progress
	viewport       *viewport.Model
	viewportHeight int
}

func InitValueModel(progress *Progress) *ValueModel {
	return &ValueModel{
		progress: progress,
		viewport: &viewport.Model{},
	}
}

func (m *ValueModel) Init() tea.Cmd {
	return nil
}

func (m *ValueModel) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *ValueModel) GetStatusLineHeight() int {
	sb := &strings.Builder{}
	m.renderStatusLine(sb)
	return lipgloss.Height(sb.String())
}

func (m *ValueModel) SetViewportHeight(viewportHeight int) {
	m.viewportHeight = viewportHeight
}

func (m *ValueModel) View() string {
	sb := &strings.Builder{}

	m.renderStatusLine(sb)

	status := m.progress.Status

	if status != StatusRunning && status != StatusError {
		return sb.String()
	}

	fmt.Fprintln(sb)

	out := m.progress.Out
	outLen := len(out)
	// TODO trim all trailing whitespace
	if outLen > 0 && out[outLen-1] == '\n' {
		out = out[:outLen-1]
	}

	viewportHeight := m.viewportHeight
	lineCount := getLineCount(out)
	if lineCount < viewportHeight {
		viewportHeight = lineCount
	}
	m.viewport.Height = viewportHeight
	m.viewport.SetContent(string(out))
	m.viewport.GotoBottom()

	fmt.Fprintf(sb, "%s", viewportStyle.Render(m.viewport.View()))

	return sb.String()
}

func (m *ValueModel) renderStatusLine(sb *strings.Builder) {
	statusLine := fmt.Sprintf(
		"%s%s%s %s",
		statusBracketStyle.Render("["),
		getStatusStyle(m.progress.Status).Render(m.progress.Status.String()),
		statusBracketStyle.Render("]"),
		m.progress.Value.GetName(),
	)
	fmt.Fprintf(
		sb,
		"%s",
		statusLineStyle.Render(statusLine),
	)
}

func getStatusStyle(status Status) lipgloss.Style {
	style := statusStyle.Copy()
	switch status {
	case StatusWaiting:
		style.Foreground(lipgloss.Color("13"))
	case StatusRunning:
		style.Foreground(lipgloss.Color("12"))
	case StatusComplete:
		style.Foreground(lipgloss.Color("10"))
	case StatusError:
		style.Foreground(lipgloss.Color("9"))
	}

	return style
}

func getLineCount(out []byte) int {
	count := 0
	for i := range out {
		if out[i] == '\n' {
			count++
		}
	}
	return count + 1
}
