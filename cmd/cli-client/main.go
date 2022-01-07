package main

import (
	"os"

	"github.com/aaltgod/gokyrie/internal/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	teams := []app.Team{
		{
			TeamName: "team1",
			TeamIP:   "154.12.32.11",
		},
		{
			TeamName: "team2",
			TeamIP:   "154.12.33.11",
		},
		{
			TeamName: "team3",
			TeamIP:   "154.12.34.11",
		},
	}

	m := app.NewModel(teams...)

	if err := tea.NewProgram(m, tea.WithAltScreen()).Start(); err != nil {
		os.Exit(1)
	}

}
