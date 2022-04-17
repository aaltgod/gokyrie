package cliclient

import (
	"log"
	"os"
	"strings"

	"github.com/aaltgod/gokyrie/internal/config"
	trafficmonitor "github.com/aaltgod/gokyrie/internal/traffic-monitor"
	"github.com/aaltgod/gokyrie/internal/tui/bubbles/selection"
	"github.com/aaltgod/gokyrie/internal/tui/bubbles/team"
	"github.com/aaltgod/gokyrie/internal/tui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/crypto/ssh/terminal"
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

type appState int

const (
	startState appState = iota
	loadState
	errorState
	quittingState
	quitState
)

type TeamEntry struct {
	Name   string
	IP     string
	bubble *team.Bubble
}

type Bubble struct {
	tm         *trafficmonitor.PcapWrapper
	width      int
	height     int
	lastResize tea.WindowSizeMsg
	state      appState
	boxes      []tea.Model
	teamEntry  []TeamEntry
	activeBox  int
	config     *config.Config
	style      *style.Styles
}

func NewBubble(cfg *config.Config) *Bubble {
	b := &Bubble{
		boxes:  make([]tea.Model, 2),
		config: cfg,
		style:  style.DefaultStyles(),
	}

	width, height, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatal("GetSize of terminal error: ", err)
	}

	b.width = width
	b.height = height
	b.state = startState

	return b
}

func (b *Bubble) Init() tea.Cmd {

	tm := trafficmonitor.NewPcapWrapper(b.config)
	// tm.StartListeners()
	b.tm = tm

	menu := selection.NewBubble(b.config.Teams, b.style)
	b.boxes[0] = menu
	b.activeBox = 0

	boxLeftWidth := b.style.Menu.GetWidth() +
		b.style.Menu.GetHorizontalFrameSize()

	// TODO: also send this along with a tea.WindowSizeMsg
	var heightMargin = lipgloss.Height(b.headerView()) +
		lipgloss.Height(b.footerView()) +
		b.style.TeamBody.GetVerticalFrameSize() +
		b.style.App.GetVerticalMargins() -
		lipgloss.Height(b.headerView()) + b.style.TeamBody.GetVerticalBorderSize()

	for _, t := range b.config.Teams {
		b.teamEntry = append(b.teamEntry, TeamEntry{
			t.Name, t.IP,
			team.NewBubble(
				t.Name, t.IP, b.style,
				b.width, boxLeftWidth,
				b.height, heightMargin,
			),
		})
	}

	b.boxes[1] = b.teamEntry[0].bubble

	b.state = loadState

	return nil
}

func (b *Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return b, tea.Quit
		case "tab":
			b.activeBox = (b.activeBox + 1) % 2
		}
	case tea.WindowSizeMsg:
		b.lastResize = msg
		b.width = msg.Width
		b.height = msg.Height

		if b.state == loadState {
			for i, bx := range b.boxes {
				m, cmd := bx.Update(msg)
				b.boxes[i] = m
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
		}
	case selection.SelectedMsg:
		b.activeBox = 1
		team := b.teamEntry[msg.Index].bubble
		b.boxes[1] = team
	case selection.ActiveMsg:
		b.boxes[1] = b.teamEntry[msg.Index].bubble
		cmds = append(cmds, func() tea.Msg {
			return b.lastResize
		})
	}

	if b.state == loadState {
		updatedBubble, cmd := b.boxes[b.activeBox].Update(msg)
		b.boxes[b.activeBox] = updatedBubble
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return b, tea.Batch(cmds...)
}

func (b Bubble) View() string {
	s := strings.Builder{}
	s.WriteString(b.headerView())
	s.WriteRune('\n')

	switch b.state {
	case loadState:
		selection := b.boxView(0)
		teamBody := b.boxView(1)
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, selection, teamBody))
	}

	s.WriteRune('\n')
	s.WriteString(b.footerView())
	return b.style.App.Render(s.String())
}

func (b Bubble) headerView() string {
	width := b.width - b.style.App.GetHorizontalFrameSize()
	return b.style.Header.Copy().Width(width).Render("gokyrie")
}

func (b Bubble) footerView() string {
	branch := b.style.Branch.Render("master")
	gap := lipgloss.NewStyle().
		Width(b.width -
			lipgloss.Width(branch) -
			b.style.App.GetHorizontalFrameSize()).
		Render("")
	footer := lipgloss.JoinHorizontal(lipgloss.Top, gap, branch)
	return b.style.Footer.Render(footer)
}
func (b *Bubble) boxView(boxIdx int) string {

	isActive := boxIdx == b.activeBox

	switch box := b.boxes[boxIdx].(type) {
	case *selection.Bubble:
		s := b.style.Menu
		if isActive {
			s = s.Copy().BorderForeground(b.style.ActiveBorderColor)
		}
		return s.Render(box.View())
	case *team.Bubble:
		box.Active = isActive
		return box.View()
	default:
		//TODO: need to add an handling of an unknown bubble
	}
	//TODO: no need to return anything here
	return ""
}
