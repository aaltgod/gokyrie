package style

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	ActiveBorderColor          lipgloss.Color
	InactiveBorderColor        lipgloss.Color
	ActiveServiceStatusColor   lipgloss.Color
	InactiveServiceStatusColor lipgloss.Color

	App    lipgloss.Style
	Header lipgloss.Style
	Footer lipgloss.Style
	Branch lipgloss.Style

	Menu             lipgloss.Style
	MenuItem         lipgloss.Style
	MenuCursor       lipgloss.Style
	SelectedMenuItem lipgloss.Style

	Home             lipgloss.Style
	HomeItem         lipgloss.Style
	HomeCursor       lipgloss.Style
	SelectedHomeItem lipgloss.Style

	HelpKey   lipgloss.Style
	HelpValue lipgloss.Style

	ServiceStatus lipgloss.Style

	TeamBodyBorder lipgloss.Border

	TeamBody lipgloss.Style
}

func DefaultStyles() *Styles {

	s := new(Styles)

	s.ActiveBorderColor = lipgloss.Color("#F98C17")
	s.InactiveBorderColor = lipgloss.Color("#FBB162")
	s.ActiveServiceStatusColor = lipgloss.Color("#EF1010")
	s.InactiveServiceStatusColor = lipgloss.Color("#3d1")

	s.App = lipgloss.NewStyle().
		Margin(1, 2)

	s.Header = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F98C17")).
		Align(lipgloss.Right).
		Bold(true)

	s.Footer = lipgloss.NewStyle().
		MarginTop(1)

	s.Branch = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3d1")).
		Background(lipgloss.Color("#14222b"))

	s.Menu = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(s.InactiveBorderColor).
		Padding(1, 2).
		MarginRight(1).
		Width(24)

	s.MenuCursor = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F98C17")).
		SetString(">")

	s.MenuItem = lipgloss.NewStyle().
		PaddingLeft(3)

	s.SelectedMenuItem = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FBB162")).
		PaddingLeft(2)

	s.Home = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(s.ActiveBorderColor).
		Padding(1, 3)

	s.HomeCursor = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F98C17")).
		SetString(">")

	s.HomeItem = lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(s.InactiveBorderColor)

	s.SelectedHomeItem = lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(s.InactiveBorderColor)

	s.HelpKey = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	s.HelpValue = lipgloss.NewStyle().
		Foreground(lipgloss.Color("239"))

	s.ServiceStatus = lipgloss.NewStyle().
		Foreground(s.InactiveServiceStatusColor)

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
		BorderForeground(s.InactiveBorderColor).
		PaddingRight(1)

	return s
}
