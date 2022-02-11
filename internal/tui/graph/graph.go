package graph

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ShowGraph bool
}

type GraphMsg struct{}
type FrameMsg struct{}

func Frame() tea.Cmd {
	return tea.Tick(time.Duration(float64(time.Second)/fps), func(time.Time) tea.Msg {
		return FrameMsg{}
	})
}

func Graph() tea.Cmd {
	return tea.Tick(time.Duration(float64(time.Second)/fps), func(time.Time) tea.Msg {
		return GraphMsg{}
	})
}

func (m *Model) ToggleGraph() tea.Cmd {
	if !m.ShowGraph {
		return m.StartGraph()
	}
	m.StopGraph()
	return nil
}

func (m *Model) StartGraph() tea.Cmd {
	m.ShowGraph = true
	return Frame()
}

func (m *Model) StopGraph() {
	m.ShowGraph = false
}

func (m Model) View() string {

	graph := Plot()

	return graph
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, Frame()
}
