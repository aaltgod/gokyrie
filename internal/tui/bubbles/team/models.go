package team

import "github.com/aaltgod/gokyrie/internal/tui/bubbles/graph"

type Team struct {
	Name, IP string
}

type Model struct {
	Team  Team
	Graph graph.Model
}
