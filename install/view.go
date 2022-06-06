package install

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type installModel struct {
	inst                  *Install
	groupModels           []groupModel
	groupsCollapsedHeight int
	windowHeight          int
}

type installMsg string

const (
	ProgressMsg installMsg = "progress"
	CompleteMsg installMsg = "complete"
	ErrorMsg    installMsg = "error"
)

func InitModel() tea.Model {
	m := &installModel{}

	return m
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
		switch msg {
		case CompleteMsg:
			return m, tea.Quit
		case ErrorMsg:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height
		m.updateViewportMaxHeight()
		return m, nil
	case *Install:
		if m.inst != nil {
			return m, nil
		}

		m.inst = msg
		m.groupModels = []groupModel{}

		for i := range m.inst.Groups {
			group := &m.inst.Groups[i]
			if group.Status == Skipped {
				continue
			}
			groupModel := NewGroupModel(group)
			m.groupModels = append(m.groupModels, *groupModel)
		}

		m.updateGroupsCollapsedHeight()
		m.updateViewportMaxHeight()
		return m, nil
	}

	return m, nil
}

func (m *installModel) updateGroupsCollapsedHeight() {
	groupsCollapsedHeight := 0
	for i := range m.groupModels {
		groupModel := &m.groupModels[i]
		groupsCollapsedHeight += groupModel.GetCollapsedHeight()
	}
	m.groupsCollapsedHeight = groupsCollapsedHeight
}

func (m *installModel) updateViewportMaxHeight() {
	if m.inst == nil {
		return
	}
	if m.windowHeight <= 0 {
		return
	}

	viewportMaxHeight := m.windowHeight - m.groupsCollapsedHeight - 2

	if viewportMaxHeight < 10 {
		viewportMaxHeight = 10
	}

	for i := range m.groupModels {
		groupModel := &m.groupModels[i]
		groupModel.SetViewportMaxHeight(viewportMaxHeight)
	}
}

func (m *installModel) View() string {
	if m.inst == nil {
		return ""
	}

	sb := &strings.Builder{}
	for i := range m.groupModels {
		groupModel := &m.groupModels[i]
		sb.WriteString(groupModel.View())
	}

	return sb.String()
}
