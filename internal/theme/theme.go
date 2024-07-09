package theme

import (
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	// Color palette
	PrimaryColor   lipgloss.Color
	SecondaryColor lipgloss.Color
	AccentColor    lipgloss.Color
	TextColor      lipgloss.AdaptiveColor
	ErrorColor     lipgloss.AdaptiveColor

	// Text styles
	NormalText lipgloss.Style
	BoldText   lipgloss.Style
	ItalicText lipgloss.Style

	// UI element styles
	Button       lipgloss.Style
	Input        lipgloss.Style
	List         lipgloss.Style
	ListItem     lipgloss.Style
	SelectedItem lipgloss.Style

	// Message styles
	UserMessage lipgloss.Style
	AIMessage   lipgloss.Style

	// Markdown styles
	Heading lipgloss.Style
	Code    lipgloss.Style
	Link    lipgloss.Style
}

var DefaultTheme = Theme{
	PrimaryColor:   lipgloss.Color("#1E90FF"),
	SecondaryColor: lipgloss.Color("#FF69B4"),
	AccentColor:    lipgloss.Color("#FFA500"),
	TextColor:      lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"},
	ErrorColor:     lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF6347"},

	NormalText: lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}),
	BoldText:   lipgloss.NewStyle().Bold(true),
	ItalicText: lipgloss.NewStyle().Italic(true),

	Button: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#1E90FF")).
		Padding(0, 1),

	Input: lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF69B4")).
		Padding(0, 1),

	List: lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#1E90FF")).
		Padding(1),

	ListItem: lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}),

	SelectedItem: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#FF69B4")).
		Bold(true),

	UserMessage: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#E0E0E0")).
		Padding(0, 1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#A0A0A0")),

	AIMessage: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#4A4A4A")).
		Padding(0, 1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#808080")),

	Heading: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFA500")).
		Bold(true),

	Code: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Background(lipgloss.Color("#333333")),

	Link: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1E90FF")).
		Underline(true),
}

var CurrentTheme = DefaultTheme

func SetTheme(theme Theme) {
	CurrentTheme = theme
}
