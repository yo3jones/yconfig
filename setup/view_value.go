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

	retryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("3"))

	statusStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Width(10)

	viewportStyle = lipgloss.NewStyle().
			MarginLeft(5).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			PaddingLeft(2).
			BorderForeground(lipgloss.Color("4"))
)

type ValueModel struct {
	state          *EntryState
	setupComplete  bool
	viewport       *viewport.Model
	viewportHeight int
}

func InitValueModel(state *EntryState) *ValueModel {
	return &ValueModel{
		state:    state,
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

	status := m.state.Status

	if !m.setupComplete && status != StatusRunning {
		return sb.String()
	} else if m.setupComplete && status != StatusError {
		return sb.String()
	}

	fmt.Fprintln(sb)

	out := m.state.Out
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

	runtimeViewportStyle := viewportStyle
	if m.state.Status == StatusError {
		runtimeViewportStyle = runtimeViewportStyle.Copy().
			BorderForeground(lipgloss.Color("1"))
	}
	fmt.Fprintf(sb, "%s", runtimeViewportStyle.Render(m.viewport.View()))

	return sb.String()
}

func (m *ValueModel) renderStatusLine(sb *strings.Builder) {
	statusLine := fmt.Sprintf(
		"%s%s%s %s %s",
		statusBracketStyle.Render("["),
		m.getStatusStyle().Render(m.state.Status.String()),
		statusBracketStyle.Render("]"),
		m.state.Entry.Name,
		m.renderRetry(),
	)
	fmt.Fprintf(
		sb,
		"%s",
		statusLineStyle.Render(statusLine),
	)
}

func (m *ValueModel) renderRetry() string {
	if !m.state.Retrying {
		return ""
	}

	return retryStyle.Render(
		fmt.Sprintf(
			"( %d of %d )",
			m.state.Tries,
			m.state.Entry.RetryCount+1,
		),
	)
}

func (m *ValueModel) getStatusStyle() lipgloss.Style {
	style := statusStyle.Copy()
	switch m.state.Status {
	case StatusWaiting:
		if m.state.Retrying {
			style.Foreground(lipgloss.Color("3"))
		} else {
			style.Foreground(lipgloss.Color("5"))
		}
	case StatusRunning:
		style.Foreground(lipgloss.Color("4"))
	case StatusComplete:
		if m.state.Retrying {
			style.Foreground(lipgloss.Color("3"))
		} else {
			style.Foreground(lipgloss.Color("2"))
		}
	case StatusError:
		style.Foreground(lipgloss.Color("1"))
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
