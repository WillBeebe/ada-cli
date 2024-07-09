package theme

// import (
// 	"github.com/charmbracelet/bubbles/list"
// 	"github.com/charmbracelet/bubbles/textinput"
// 	"github.com/charmbracelet/lipgloss"
// 	"github.com/container-labs/ada/internal/theme"
// )

// type Button struct {
// 	Text    string
// 	OnClick func()
// }

// func (b Button) View(focused bool) string {
// 	style := theme.CurrentTheme.Button
// 	if focused {
// 		style = style.Copy().Background(theme.CurrentTheme.AccentColor)
// 	}
// 	return style.Render(b.Text)
// }

// type Input struct {
// 	textinput.Model
// }

// func NewInput(placeholder string) Input {
// 	ti := textinput.New()
// 	ti.Placeholder = placeholder
// 	ti.PromptStyle = theme.CurrentTheme.NormalText
// 	ti.TextStyle = theme.CurrentTheme.NormalText
// 	return Input{Model: ti}
// }

// type List struct {
// 	list.Model
// }

// func NewList(title string, items []list.Item) List {
// 	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
// 	l.Title = title
// 	l.Styles.Title = theme.CurrentTheme.BoldText.Copy().Foreground(theme.CurrentTheme.AccentColor)
// 	return List{Model: l}
// }

// type Message struct {
// 	Content string
// 	IsUser  bool
// }

// func (m Message) View() string {
// 	var style lipgloss.Style
// 	if m.IsUser {
// 		style = theme.CurrentTheme.UserMessage
// 	} else {
// 		style = theme.CurrentTheme.AIMessage
// 	}
// 	return style.Render(m.Content)
// }

// Add more components as needed
