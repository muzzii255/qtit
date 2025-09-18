package dashboard

import "github.com/charmbracelet/lipgloss"

// Base colors for dark theme
const (
	ColorBackground = "#282c34"
	ColorForeground = "#abb2bf"
	ColorAccent     = "#61afef"
	ColorSuccess    = "#98c379"
	ColorError      = "#e06c75"
	ColorSubtle     = "#5c6370"
	ColorPanel      = "#1e222a"
)

// Theme colors as lipgloss.Color
var (
	BgColor      = lipgloss.Color(ColorBackground)
	FgColor      = lipgloss.Color(ColorForeground)
	AccentColor  = lipgloss.Color(ColorAccent)
	SuccessColor = lipgloss.Color(ColorSuccess)
	ErrorColor   = lipgloss.Color(ColorError)
	SubtleColor  = lipgloss.Color(ColorSubtle)
	PanelColor   = lipgloss.Color(ColorPanel)
)

// Common styles
var (
	// Base style for general text
	BaseStyle = lipgloss.NewStyle().
		Foreground(FgColor)

	// Header style for titles and headings
	HeaderStyle = lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true)

	// Selected row style
	SelectedStyle = lipgloss.NewStyle().
		Foreground(SuccessColor).
		Bold(true)

	// Subtle text style for controls and less important info
	SubtleStyle = lipgloss.NewStyle().
		Foreground(SubtleColor)

	// Status success style
	StatusSuccess = lipgloss.NewStyle().
		Foreground(SuccessColor)

	// Status error style
	StatusError = lipgloss.NewStyle().
		Foreground(ErrorColor)

	// Command prompt style
	CommandPromptStyle = lipgloss.NewStyle().
		Foreground(AccentColor)

	// Command input text style
	CommandTextStyle = lipgloss.NewStyle().
		Foreground(FgColor)

	// Panel/box background style
	PanelStyle = lipgloss.NewStyle().
		Foreground(FgColor).
		Background(PanelColor).
		Padding(0, 1).
		MarginTop(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(AccentColor)

	// Paginator active dot style
	PaginatorActiveDot = lipgloss.NewStyle().
		Foreground(AccentColor)

	// Paginator inactive dot style
	PaginatorInactiveDot = lipgloss.NewStyle().
		Foreground(SubtleColor)
)