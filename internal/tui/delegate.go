package cliclient

import (
	"github.com/aaltgod/gokyrie/internal/tui/bubbles/graph"
	"github.com/aaltgod/gokyrie/internal/tui/bubbles/team"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var statusMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
	Render

func NewItemDelegate(keys *DelegateKeyMap) list.DefaultDelegate {

	var d = list.NewDefaultDelegate()

	d.SetSpacing(2)

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var teamName string

		if i, ok := m.SelectedItem().(team.Team); ok {
			teamName = i.Name
		} else {
			return nil
		}
		_ = teamName

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.Enter):
				return graph.Graph()
			case key.Matches(msg, keys.BackSpace):
				return m.NewStatusMessage("BACK")
			}
		}

		return nil
	}

	return d
}

func NewDelegateKeyMap() *DelegateKeyMap {
	return &DelegateKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		BackSpace: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("backspace", "back"),
		),
	}
}
