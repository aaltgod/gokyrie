package main

import (
	"context"
	"log"
	"os"

	"github.com/aaltgod/gokyrie/internal/config"
	trafficmonitor "github.com/aaltgod/gokyrie/internal/traffic-monitor"
	app "github.com/aaltgod/gokyrie/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/coral"
	"golang.org/x/sync/errgroup"
)

func main() {
	var rootCmd = &coral.Command{
		Use:     "gokyrie",
		Short:   "Traffic monitor",
		Version: "0.1",
		Args:    coral.MaximumNArgs(1),
		RunE: func(cmd *coral.Command, args []string) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}

			ctx, cancel := context.WithCancel(context.Background())
			eg, ctx := errgroup.WithContext(ctx)

			dataCh := make(chan trafficmonitor.Data, len(cfg.Teams))

			eg.Go(func() error {
				tm := trafficmonitor.NewPcapWrapper(cfg, dataCh)
				return tm.StartListeners(ctx)
			})

			eg.Go(func() error {
				defer cancel()
				m := app.NewBubble(ctx, cfg, dataCh)
				p := tea.NewProgram(
					m, tea.WithAltScreen(),
					tea.WithMouseAllMotion(),
				)
				return p.Start()
			})

			return eg.Wait()
		},
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	os.Exit(0)
}
