package main

import (
	"log"
	"os"

	app "github.com/aaltgod/gokyrie/internal/tui"
	"github.com/aaltgod/gokyrie/internal/tui/team"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	logger, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	teams := []team.Team{
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
		{
			TeamName: "team4",
			TeamIP:   "154.12.34.11",
		},
		{
			TeamName: "team5",
			TeamIP:   "154.12.34.11",
		},
		{
			TeamName: "team6",
			TeamIP:   "154.12.34.11",
		},
		{
			TeamName: "team7",
			TeamIP:   "154.12.34.11",
		},
		{
			TeamName: "team8",
			TeamIP:   "154.12.34.11",
		},
		{
			TeamName: "team9",
			TeamIP:   "154.12.34.11",
		},
		{
			TeamName: "team10",
			TeamIP:   "154.12.34.11",
		},
	}

	m := app.NewModel(logger, teams...)

	if err := tea.NewProgram(m, tea.WithAltScreen()).Start(); err != nil {
		m.Logger.WriteString("App doesn't start: " + err.Error())
		os.Exit(1)
	}

}
