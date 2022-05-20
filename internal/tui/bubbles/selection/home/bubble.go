package home

import (
	"strings"

	"github.com/aaltgod/gokyrie/internal/config"
	"github.com/aaltgod/gokyrie/internal/tui/style"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

var symbols = []string{
	"①", "②", "③",
	"④", "⑤", "⑥",
	"⑦", "⑧", "⑨",
	"⑩",
}

type Bubble struct {
	Services     []config.Service
	Teams        []config.Team
	SelectedTeam int
	width        int
	widthMargin  int
	height       int
	heightMargin int
	styles       *style.Styles
	viewport     *viewport.Model
}

func NewBubble(services []config.Service, teams []config.Team, width, widthMargin, height, heightMargin int, styles *style.Styles) *Bubble {
	b := &Bubble{
		Services:     services,
		Teams:        teams,
		styles:       styles,
		widthMargin:  widthMargin,
		heightMargin: heightMargin,
		viewport: &viewport.Model{
			Width:  width,
			Height: height,
		},
	}
	b.setSize(width, height)
	return b
}

func (b *Bubble) Init() tea.Cmd {
	return nil
}

func (b Bubble) View() string {
	var (
		homeMaxItemWidth = 15
		items            = make([]string, len(b.Teams))
	)

	for i, item := range b.Teams {
		teamName := truncate.StringWithTail(item.Name, uint(homeMaxItemWidth), "…")
		teamAddress := truncate.StringWithTail(item.IP, uint(homeMaxItemWidth), "…")
		if i == b.SelectedTeam {
			items[i] = b.styles.SelectedHomeItem.Render(
				teamName + "\n" + teamAddress + "\n" + b.serviceStatuses())
		} else {
			items[i] = b.styles.HomeItem.Render(
				teamName + "\n" + teamAddress + "\n" + b.serviceStatuses())
		}
	}

	homeWitdh := b.width - b.widthMargin - 10
	columnAmount := homeWitdh / homeMaxItemWidth
	remainigItems := len(b.Teams) - columnAmount
	lineAmount := 1

	for {
		if remainigItems > 0 {
			remainigItems -= columnAmount
			lineAmount++
		} else {
			break
		}
	}

	itemLines := make([]string, lineAmount)
	offset := 0
	for i := 0; i < lineAmount; i++ {
		lastIdx := offset + columnAmount
		// FIXME: refactore it!!!
		if lastIdx > len(items)-1 {
			lastIdx = len(items)
		}
		itemLines[i] = lipgloss.JoinHorizontal(lipgloss.Center, items[offset:lastIdx]...)
		offset += columnAmount
	}

	home := b.styles.Home.Copy()

	body := home.Width(b.width - b.widthMargin - 10).
		Height(b.height - b.heightMargin - lipgloss.Height(b.headerView()) - 1).
		Render(lipgloss.JoinVertical(lipgloss.Top, itemLines...))

	b.viewport.SetContent(body)
	return b.viewport.View()
}

func (b *Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.setSize(msg.Width, msg.Height)
	}

	return b, tea.Batch(cmds...)
}

func (b *Bubble) setSize(w, h int) {
	b.width = w
	b.height = h
	b.viewport.Width = b.width - b.widthMargin
	b.viewport.Height = b.height - b.heightMargin
}

func (b Bubble) headerView() string {
	return ""
}

func (b Bubble) serviceStatuses() string {
	s := strings.Builder{}
	for i := range b.Services {
		s.WriteString(
			b.styles.ServiceStatus.Copy().
				Foreground(b.styles.ActiveServiceStatusColor).Render(symbols[i]),
		)
		s.WriteString(" ")
	}
	return s.String()
}
