package app

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#AD40FF")).
			Padding(0, 12)
)

func (t Team) Title() string       { return t.TeamName }
func (t Team) Description() string { return t.TeamIP }
func (t Team) FilterValue() string { return t.TeamName }

func NewModel(teams ...Team) Model {

	items := make([]list.Item, len(teams))
	for i, team := range teams {
		items[i] = team
	}

	delegate := list.NewDefaultDelegate()
	teamList := list.NewModel(items, delegate, 0, 0)
	teamList.Title = "Gokyrie"
	teamList.Styles.Title = titleStyle

	return Model{
		List: teamList,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		topGap, rightGap, bottomGap, leftGap := appStyle.GetPadding()
		m.List.SetSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctr+c":
			return m, tea.Quit
		case "s":
			cmd := m.List.ToggleSpinner()
			return m, cmd
		}
	}

	newListModel, cmd := m.List.Update(msg)
	m.List = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return appStyle.Render(m.List.View())
}
