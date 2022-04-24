package team

import (
	trafficmonitor "github.com/aaltgod/gokyrie/internal/traffic-monitor"
	"github.com/aaltgod/gokyrie/internal/tui/bubbles/team/flags"
	stat "github.com/aaltgod/gokyrie/internal/tui/bubbles/team/stats"
	"github.com/aaltgod/gokyrie/internal/tui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	statState state = iota
	flagState
	configState
)

type Bubble struct {
	dataCh       chan trafficmonitor.Data
	state        state
	serviceName  string
	servicePort  string
	graph        string
	width        int
	widthMargin  int
	height       int
	heightMargin int
	boxes        []tea.Model

	style *style.Styles

	Active bool
}

func NewBubble(dataCh chan trafficmonitor.Data, serviceName, servicePort string, styles *style.Styles, width, wm, height, hm int) *Bubble {
	b := &Bubble{
		dataCh:       dataCh,
		serviceName:  serviceName,
		servicePort:  servicePort,
		state:        statState,
		style:        styles,
		width:        width,
		widthMargin:  wm,
		height:       height,
		heightMargin: hm,
		boxes:        make([]tea.Model, 2),
	}

	b.boxes[statState] = stat.NewBubble(
		dataCh, serviceName, servicePort,
		b.style, width, wm+b.style.TeamBody.GetHorizontalBorderSize(),
		height, hm+lipgloss.Height(b.headerView())+styles.TeamBody.GetVerticalBorderSize(),
	)
	b.boxes[flagState] = flags.NewBubble(
		serviceName, servicePort,
		b.style, width, wm+b.style.TeamBody.GetHorizontalBorderSize(),
		height, hm+lipgloss.Height(b.headerView())+styles.TeamBody.GetVerticalBorderSize(),
	)

	return b
}

func (b *Bubble) Init() tea.Cmd {
	return nil
}

func (b *Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.height = msg.Height
		b.width = msg.Width
		for i, bx := range b.boxes {
			m, cmd := bx.Update(msg)
			b.boxes[i] = m
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	updatedBubble, cmd := b.boxes[b.state].Update(msg)
	b.boxes[b.state] = updatedBubble
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return b, tea.Batch(cmds...)
}

func (b Bubble) View() string {

	header := b.headerView()
	ts := b.style.TeamBody.Copy()
	if b.Active {
		ts = ts.BorderForeground(b.style.ActiveBorderColor)
	}

	body := ts.Width(b.width - b.widthMargin - b.style.TeamBody.GetVerticalFrameSize()).
		Height(b.height - b.heightMargin - lipgloss.Height(header)).
		Render(b.boxes[b.state].View())

	return body
}

func (b Bubble) headerView() string {
	return ""
}
