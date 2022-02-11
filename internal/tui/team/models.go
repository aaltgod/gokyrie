package team

import "github.com/aaltgod/gokyrie/internal/tui/graph"

type Team struct {
	TeamName, TeamIP string
}

type Model struct {
	Team  Team
	Graph graph.Model
}
