package styles

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var width = 80
var paddingWidth = 40

func init() {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	width = physicalWidth
	BaseStyle.Width(width - 5)
}

func Width() int {
	return width
}

var BaseStyle = lipgloss.NewStyle().
	// Bold(true).
	BorderStyle(lipgloss.RoundedBorder()).
	// Background(lipgloss.Color("#e0e0e0")).
	// PaddingLeft(1).
	Padding(1).
	// PaddingBottom(2).
	Width(60)

var ContentStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("0")).
	// Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "248"}).
	// Background(lipgloss.Color("#e0e0e0")).
	Padding(0, 1, 0).Width(width)

var ChatContentStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	Foreground(lipgloss.Color("0")).
	Align(lipgloss.Left).
	// Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "248"}).
	Background(lipgloss.Color("#e9adff")).
	Padding(1).
	Margin(1).MaxWidth(width - paddingWidth)

var ChatUserContentStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	Foreground(lipgloss.Color("0")).
	Align(lipgloss.Right).
	// Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "248"}).
	Background(lipgloss.Color("#a53bcc")).
	Padding(1).
	Margin(1).MaxWidth(width - paddingWidth)

var ContentInfoStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("7")).
	// Background(lipgloss.Color("3")).
	Padding(0, 1, 0).Width(width)

var ContentErrorStyle = lipgloss.NewStyle().
	// Foreground(lipgloss.Color("7")).
	Foreground(lipgloss.Color("1")).
	// Background(lipgloss.Color("#e0e0e0")).
	Padding(0, 1, 0).Width(width)

var ChatPadding = lipgloss.NewStyle().
	Foreground(lipgloss.Color("0")).
	Align(lipgloss.Center).
	// Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "248"}).
	// Background(lipgloss.Color("#e0e0e0")).
	Padding(0).
	Width(paddingWidth)

var PromptContentStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("0")).
	Align(lipgloss.Center).
	// Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "248"}).
	Background(lipgloss.Color("#e0e0e0")).
	Padding(0, 1, 0, 1).
	Margin(0, 1, 0, 0)

var HelpStyle = lipgloss.NewStyle().
	Align(lipgloss.Left).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("0")).
	Bold(true).
	// fun mode
	Background(lipgloss.Color("#282a35")).
	Foreground(lipgloss.Color("#FF79C6")).
	// regular mode
	// Background(lipgloss.Color("#202124")).
	// Foreground(lipgloss.Color("#e8e8e8"))
	PaddingLeft(2).
	PaddingTop(1).
	PaddingRight(2).
	PaddingBottom(1)

	// misc
	// #b536ff
	// #ca03fc
	// #f403fc
