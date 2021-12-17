package main

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	TeamNames, TeamIPs []string
}

func newModel() model {
	return model{
		TeamNames: []string{"team1", "team2", "team3"},
		TeamIPs:   []string{"154.12.32.11", "154.12.33.11", "154.12.34.11"},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctr+c":
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m model) View() string {

	var b strings.Builder
	b.WriteString("Teams\tIPS\n\n")
	for i, item := range m.TeamNames {
		b.WriteString(item)
		b.WriteString("\t")
		b.WriteString(m.TeamIPs[i])
		b.WriteString("\n")
	}

	return b.String()
}

func main() {

	m := newModel()

	if err := tea.NewProgram(m, tea.WithAltScreen()).Start(); err != nil {
		os.Exit(1)
	}

}
