package team

import tea "github.com/charmbracelet/bubbletea"

func (t Team) Title() string       { return t.TeamName }
func (t Team) Description() string { return t.TeamIP }
func (t Team) FilterValue() string { return t.TeamName }

func (m Model) Update() (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	output := m.Team.TeamName + "\n"
	output += m.Graph.View()
	return output
}
