package setup

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	progress          []*Progress
	valueModels       []*ValueModel
	statusLineHeights int
	windowHeight      int
	windowWidth       int
}

type Msg string

const (
	MsgDone    Msg = "done"
	MsgRefresh Msg = "refresh"
)

func InitModel() tea.Model {
	return &Model{}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []*Progress:
		m.progress = msg
		m.initialzeValueModels()
		return m, nil
	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height
		m.windowWidth = msg.Width
		m.initialzeViewport()
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

	hasError := false
	allComplete := true
	for i, valueModel := range m.valueModels {
		progress := m.progress[i]
		if progress.Status == StatusError {
			hasError = true
		}
		if progress.Status != StatusComplete {
			allComplete = false
		}
		if i > 0 {
			fmt.Fprintln(sb)
		}
		fmt.Fprint(sb, valueModel.View())
	}
	if hasError || allComplete {
		fmt.Fprintln(sb)
	}

	return sb.String()
}

func (m *Model) initialzeValueModels() {
	if m.progress == nil || m.valueModels != nil {
		return
	}

	valueModels := make([]*ValueModel, len(m.progress))
	statusLineHeights := 0
	for i, p := range m.progress {
		valueModel := InitValueModel(p)
		valueModels[i] = valueModel
		statusLineHeights += valueModel.GetStatusLineHeight()
	}
	m.valueModels = valueModels
	m.statusLineHeights = statusLineHeights

	m.initialzeViewport()
}

func (m *Model) initialzeViewport() {
	if m.valueModels == nil || m.windowHeight <= 0 {
		return
	}

	viewportHeight := m.windowHeight - m.statusLineHeights - 1
	if viewportHeight < 5 {
		viewportHeight = 5
	} else if viewportHeight > 200 {
		viewportHeight = 200
	}

	for _, valueModel := range m.valueModels {
		valueModel.SetViewportHeight(viewportHeight)
	}
}
