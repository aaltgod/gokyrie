package cliclient

import (
	"os"
	"time"

	"github.com/aaltgod/gokyrie/internal/tui/graph"
	"github.com/aaltgod/gokyrie/internal/tui/team"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding()

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#AD40FF")).
			Padding(0, 12)

	viewportStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#AD40FF")).
			PaddingTop(1).
			PaddingRight(80).
			PaddingBottom(5).
			PaddingLeft(1)

	graphStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#AD40FF"))

	layout = "Jan 2, 2006 at 3:04:05pm"
)

func NewModel(logger *os.File, teams ...team.Team) Model {

	items := make([]list.Item, len(teams))
	for i, t := range teams {
		items[i] = t
	}

	viewPort := viewport.Model{
		Height: 40,
		Width:  30,
	}

	delegate := NewItemDelegate(NewDelegateKeyMap())
	teamList := list.NewModel(items, delegate, 10, 15)
	teamList.Title = "Gokyrie"
	teamList.Styles.Title = titleStyle

	return Model{
		Logger:   logger,
		List:     teamList,
		ViewPort: viewPort,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		topGap, rightGap, bottomGap, leftGap := appStyle.GetPadding()
		m.List.SetSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctr+c":
			return m, tea.Quit
		case "s":
			cmd = m.List.ToggleSpinner()
			return m, cmd
		default:
			m.Logger.WriteString(time.Now().Format(layout) + "\tEnter " + msg.String() + "\n")
		}
	case graph.GraphMsg:
		cmd = m.Graph.ToggleGraph()
		return m, cmd
	case graph.FrameMsg:
		m.Graph, cmd = m.Graph.Update(msg)
		return m, cmd
	}

	newListModel, cmd := m.List.Update(msg)
	m.List = newListModel
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var graphs string

	if m.Graph.ShowGraph {
		graphs = lipgloss.JoinVertical(
			lipgloss.Left,
			graphStyle.Render(m.Graph.View()),
			graphStyle.Render(m.Graph.View()),
			graphStyle.Render(m.Graph.View()),
			graphStyle.Render(m.Graph.View()),
			graphStyle.Render(m.Graph.View()),
			graphStyle.Render(m.Graph.View()),
			graphStyle.Render(m.Graph.View()),
			graphStyle.Render(m.Graph.View()),
			graphStyle.Render(m.Graph.View()),
			graphStyle.Render(m.Graph.View()),
		)
		m.ViewPort.SetContent(graphs)
	}

	viewport := viewportStyle.Render(m.ViewPort.View())

	return appStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, m.List.View(), viewport))
}
