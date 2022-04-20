package stat

import (
	"time"

	trafficmonitor "github.com/aaltgod/gokyrie/internal/traffic-monitor"
	"github.com/aaltgod/gokyrie/internal/tui/bubbles/team/stats/graph"
	"github.com/aaltgod/gokyrie/internal/tui/style"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Bubble struct {
	dataCh       chan trafficmonitor.Data
	serviceName  string
	servicePort  string
	style        *style.Styles
	width        int
	widthMargin  int
	height       int
	heightMargin int
	viewport     *viewport.Model
}

func NewBubble(dataCh chan trafficmonitor.Data, serviceName, servicePort string, styles *style.Styles, width, wm, height, hw int) *Bubble {
	b := &Bubble{
		dataCh:       dataCh,
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

	b.setSize(width, height)

	return b
}

func (b *Bubble) Init() tea.Cmd {
	return nil
}

func (b *Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.setSize(msg.Width, msg.Height)
	case DataMsg:
		b.viewport.SetContent(graph.Plot() + msg.IP + msg.Text)
	}

	if cmd := b.getData(); cmd != nil {
		cmds = append(cmds, cmd)
	}

	return b, tea.Batch(cmds...)
}

func (b Bubble) View() string {
	b.viewport.GotoTop()
	return b.viewport.View()
}

func (b *Bubble) setSize(width, height int) {
	b.width = width
	b.height = height
	b.viewport.Width = width - b.widthMargin
	b.viewport.Height = height - b.heightMargin
}

type DataMsg struct {
	IP   string
	Text string
}

func (b *Bubble) getData() tea.Cmd {
	return tea.Tick(time.Millisecond*1000, func(t time.Time) tea.Msg {
		for {
			select {
			case d := <-b.dataCh:
				return DataMsg{
					IP:   d.IP,
					Text: d.Text,
				}
			default:
				return nil
			}
		}
	})
}
