package team

import tea "github.com/charmbracelet/bubbletea"

func (t Team) Title() string       { return t.Name }
func (t Team) Description() string { return t.IP }
func (t Team) FilterValue() string { return t.Name }

func (m Model) Update() (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	output := m.Team.Name + "\n"
	output += m.Graph.View()
	return output
}
