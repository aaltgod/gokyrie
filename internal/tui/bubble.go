package cliclient

import (
	"strings"

	"github.com/aaltgod/gokyrie/internal/config"
	"github.com/aaltgod/gokyrie/internal/tui/bubbles/selection"
	"github.com/aaltgod/gokyrie/internal/tui/style"
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

type appState int

const (
	startState appState = iota
	errorState
	quittingState
	quitState
)

type Bubble struct {
	width, height int
	lastResize    tea.WindowSizeMsg
	state         appState
	boxes         []tea.Model
	activeBox     int
	config        *config.Config
	styles        *style.Styles
}

func NewBubble(cfg *config.Config) *Bubble {
	b := &Bubble{
		boxes:  make([]tea.Model, 2),
		config: cfg,
		styles: style.DefaultStyles(),
	}
	b.state = startState
	return b
}

func (b *Bubble) Init() tea.Cmd {
	menu := selection.NewBubble(b.config.Teams, b.styles)
	b.boxes[0] = menu
	b.activeBox = 0

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
	case selection.SelectedMsg:

	case selection.ActiveMsg:

	}

	if b.state == startState {
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
	case startState:
		selection := b.boxView(0)
		teamBody := b.boxView(1)
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, selection, teamBody))
	default:
		break
	}
	return b.styles.App.Render(s.String())
}

func (b Bubble) headerView() string {
	width := b.width - b.styles.App.GetHorizontalFrameSize()
	return b.styles.Header.Copy().Width(width).Render("gokyrie")
}

func (b *Bubble) boxView(boxIdx int) string {
	isActive := boxIdx == b.activeBox

	switch box := b.boxes[boxIdx].(type) {
	case *selection.Bubble:
		s := b.styles.Menu
		if isActive {
			s = s.Copy().BorderForeground(b.styles.ActiveBorderColor)
		}
		return s.Render(box.View())
	default:
		break
	}

	return ""
}
