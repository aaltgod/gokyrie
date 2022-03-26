package style

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	ActiveBorderColor   lipgloss.Color
	InactiveBorderColor lipgloss.Color

	App    lipgloss.Style
	Header lipgloss.Style

	Menu             lipgloss.Style
	MenuItem         lipgloss.Style
	MenuCursor       lipgloss.Style
	SelectedMenuItem lipgloss.Style

	TeamBodyBorder lipgloss.Border

	TeamBody lipgloss.Style
}

func DefaultStyles() *Styles {

	s := new(Styles)

	s.ActiveBorderColor = lipgloss.Color("#3d1967")
	s.InactiveBorderColor = lipgloss.Color("#413b47")

	s.App = lipgloss.NewStyle().
		Margin(1, 2)

	s.Header = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#413b47")).
		Align(lipgloss.Right).
		Bold(true)

	s.Menu = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(s.InactiveBorderColor).
		Padding(1, 2).
		MarginRight(1).
		Width(24)

	s.MenuCursor = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#008712")).
		SetString(">")

	s.MenuItem = lipgloss.NewStyle().
		PaddingLeft(2)

	s.SelectedMenuItem = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#908208")).
		PaddingLeft(2)

	s.TeamBodyBorder = lipgloss.Border{
		Top:         "━",
		Bottom:      "━",
		Left:        "┃",
		Right:       "┃",
		TopLeft:     "┏",
		TopRight:    "┓",
		BottomLeft:  "┗",
		BottomRight: "┛",
	}

	s.TeamBody = lipgloss.NewStyle().
		BorderStyle(s.TeamBodyBorder).
		BorderBackground(s.InactiveBorderColor).
		PaddingRight(1)

	return s
}
