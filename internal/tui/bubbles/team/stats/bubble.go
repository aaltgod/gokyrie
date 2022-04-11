package stat

import (
	"github.com/aaltgod/gokyrie/internal/tui/style"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Bubble struct {
	serviceName  string
	servicePort  string
	style        *style.Styles
	width        int
	widthMargin  int
	height       int
	heightMargin int
	viewport     *viewport.Model
}

func NewBubble(serviceName, servicePort string, styles *style.Styles, width, wm, height, hw int) *Bubble {
	b := &Bubble{
		serviceName:  serviceName,
		servicePort:  servicePort,
		style:        styles,
		widthMargin:  wm,
		heightMargin: hw,
		viewport: &viewport.Model{
			Height: height,
			Width:  width,
		},
	}

	b.SetSize(width, height)

	return b
}

func (b *Bubble) Init() tea.Cmd {
	return nil
}

func (b *Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.SetSize(msg.Width, msg.Height)
	}

	return b, tea.Batch(cmds...)
}

func (b Bubble) View() string {
	b.viewport.SetContent(b.serviceName + " " + b.servicePort)
	b.viewport.GotoTop()
	return b.viewport.View()
}

func (b *Bubble) SetSize(width, height int) {
	b.width = width
	b.height = height
	b.viewport.Width = width - b.widthMargin
	b.viewport.Height = height - b.heightMargin
}
