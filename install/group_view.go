package install

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type groupModel struct {
	group         *Group
	commandModels []commandModel
}

var groupLineStyle = lipgloss.NewStyle().
	PaddingTop(1).
	Bold(true).
	Underline(true)

func NewGroupModel(group *Group) *groupModel {
	commandModels := []commandModel{}
	for i := range group.Commands {
		command := &group.Commands[i]
		if command.Status == Skipped {
			continue
		}

		commandModel := NewCommandModel(command)
		commandModels = append(commandModels, *commandModel)
	}

	return &groupModel{
		group:         group,
		commandModels: commandModels,
	}
}

func (m *groupModel) GetCollapsedHeight() int {
	groupHeight := lipgloss.Height(m.groupLineView(&Group{Name: "foo"}))

	commandsHeight := 0
	for i := range m.commandModels {
		commandsModel := m.commandModels[i]
		commandsHeight += commandsModel.GetCollapsedHeight()
	}

	return groupHeight + commandsHeight
}

func (m *groupModel) SetViewportMaxHeight(viewportMaxHeight int) {
	for i := range m.commandModels {
		commandModel := &m.commandModels[i]
		commandModel.SetViewportMaxHeight(viewportMaxHeight)
	}
}

func (*groupModel) Init() tea.Cmd {
	return nil
}

func (m *groupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *groupModel) View() string {
	sb := &strings.Builder{}

	fmt.Fprintf(sb, "%s\n", m.groupLineView(m.group))
	m.commandsView(sb)

	return sb.String()
}

func (m *groupModel) groupLineView(group *Group) string {
	return groupLineStyle.Render(group.Name)
}

func (m *groupModel) commandsView(sb *strings.Builder) {
	for i := range m.commandModels {
		commandModel := &m.commandModels[i]
		fmt.Fprintf(sb, "%s", commandModel.View())
	}
}
