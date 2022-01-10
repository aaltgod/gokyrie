package app

import (
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

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var teamName string

		if i, ok := m.SelectedItem().(Team); ok {
			teamName = i.TeamName
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.Enter):
				return m.NewStatusMessage(statusMessageStyle("You chose " + teamName))
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
