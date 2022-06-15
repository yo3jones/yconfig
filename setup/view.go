package setup

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	state             *SetupState
	valueModels       []*ValueModel
	statusLineHeights int
	windowHeight      int
	windowWidth       int
}

type Msg string

const (
	MsgDone    Msg = "done"
	MsgRefresh Msg = "refresh"

	minViewportHeight = 5
	maxViewportHeight = 200
)

func InitModel() tea.Model {
	return &Model{}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case *SetupState:
		m.state = msg
		m.initialzeValueModels()
		return m, nil
	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height
		m.windowWidth = msg.Width
		m.updateViewport()
		return m, nil
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case Msg:
		if msg == MsgDone {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *Model) View() string {
	if m.valueModels == nil {
		return ""
	}

	sb := &strings.Builder{}

	for _, valueModel := range m.valueModels {
		fmt.Fprintf(sb, "%s\n", valueModel.View())
	}

	return sb.String()
}

func (m *Model) initialzeValueModels() {
	m.updateComplete()
	m.updateViewport()

	if m.state == nil || m.valueModels != nil {
		return
	}

	valueModels := make([]*ValueModel, len(m.state.EntryStates))
	statusLineHeights := 0
	for i, entryState := range m.state.EntryStates {
		valueModel := InitValueModel(entryState)
		valueModels[i] = valueModel
		statusLineHeights += valueModel.GetStatusLineHeight()
	}
	m.valueModels = valueModels
	m.statusLineHeights = statusLineHeights

	m.updateComplete()
	m.updateViewport()
}

func (m *Model) updateViewport() {
	if m.valueModels == nil || m.windowHeight <= 0 {
		return
	}

	viewportHeight := m.windowHeight - m.statusLineHeights - 2

	complete := m.state.Status.IsCompleted()

	if complete {
		viewportHeight -= 2
	}

	if complete && m.state.ErroredCount > 1 {
		viewportHeight = viewportHeight / m.state.ErroredCount
	}

	if viewportHeight < minViewportHeight {
		viewportHeight = minViewportHeight
	} else if viewportHeight > maxViewportHeight {
		viewportHeight = maxViewportHeight
	}

	for _, valueModel := range m.valueModels {
		valueModel.SetViewportHeight(viewportHeight)
	}
}

func (m *Model) updateComplete() {
	if m.state == nil {
		return
	}

	if m.valueModels == nil {
		return
	}

	for _, valueModel := range m.valueModels {
		valueModel.setupComplete = m.state.Status.IsCompleted()
	}
}
