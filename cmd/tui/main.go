package main

import (
	"log"
	"os"

	"github.com/aaltgod/gokyrie/internal/config"
	app "github.com/aaltgod/gokyrie/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	logger, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	teams := []config.Team{
		{
			Name: "team1",
			IP:   "154.12.32.11",
		},
		{
			Name: "team2",
			IP:   "154.12.33.11",
		},
		{
			Name: "team3",
			IP:   "154.12.34.11",
		},
		{
			Name: "team4",
			IP:   "154.12.34.11",
		},
		{
			Name: "team5",
			IP:   "154.12.34.11",
		},
		{
			Name: "team6",
			IP:   "154.12.34.11",
		},
		{
			Name: "team7",
			IP:   "154.12.34.11",
		},
		{
			Name: "team8",
			IP:   "154.12.34.11",
		},
		{
			Name: "team9",
			IP:   "154.12.34.11",
		},
		{
			Name: "team10",
			IP:   "154.12.34.11",
		},
	}

	// m := app.NewModel(logger, teams...)

	cfg := config.NewConfig(teams...)
	m := app.NewBubble(cfg)

	if err := tea.NewProgram(m, tea.WithAltScreen()).Start(); err != nil {
		// m.Logger.WriteString("App doesn't start: " + err.Error())
		os.Exit(1)
	}

}
