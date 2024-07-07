package chat

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var width = 90
var paddingWidth = 40

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"})
	promptStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#383838", Dark: "#D9DCCF"})
	textStyle   = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"})

	commandStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F87"))
	boldStyle    = lipgloss.NewStyle().Bold(true)
	italicStyle  = lipgloss.NewStyle().Italic(true)

	userMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#E0E0E0")).
			Padding(0, 1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#A0A0A0"))

	aiMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#4A4A4A")).
			Padding(0, 1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#808080"))

	// Markdown styles
	headingStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Bold(true)
	codeStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Background(lipgloss.Color("#333333"))
	linkStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#1E90FF")).Underline(true)
	listItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4"))
)

const (
	cursorText = "|"
)

func init() {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	width = physicalWidth
	paddingWidth = width * 3 / 10
}

var ChatContentStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	Foreground(lipgloss.Color("0")).
	Align(lipgloss.Left).
	// Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "248"}).
	Background(lipgloss.Color("#323ea8")).
	Padding(1).
	Margin(1)
var ChatUserContentStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	Foreground(lipgloss.Color("0")).
	Align(lipgloss.Right).
	// Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "248"}).
	Background(lipgloss.Color("#a53bcc")).
	Padding(1).
	Margin(1)
var ChatPadding = lipgloss.NewStyle().
	Foreground(lipgloss.Color("0")).
	Align(lipgloss.Center).
	// Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "248"}).
	// Background(lipgloss.Color("#e0e0e0")).
	Padding(0).
	Width(paddingWidth)
