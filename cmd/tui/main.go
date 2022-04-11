package main

import (
	"log"
	"os"

	"github.com/aaltgod/gokyrie/internal/config"
	app "github.com/aaltgod/gokyrie/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/coral"
)

func main() {
	var rootCmd = &coral.Command{
		Use:     "gokyrie",
		Short:   "Traffic monitor",
		Version: "0.1",
		Args:    coral.MaximumNArgs(1),
		Run: func(cmd *coral.Command, args []string) {
			cfg, err := config.GetConfig()
			if err != nil {
				log.Fatal(err)
			}

			m := app.NewBubble(cfg)

			if err := tea.NewProgram(m, tea.WithAltScreen()).Start(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
