package theme

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

func NewStyledTextInput(placeholder string, charLimit int) textinput.Model {
	t := textinput.New()
	t.Placeholder = placeholder
	t.CharLimit = charLimit
	t.Cursor.Style = CurrentTheme.NormalText.Copy().Foreground(CurrentTheme.AccentColor)
	t.PromptStyle = CurrentTheme.NormalText.Copy().Foreground(CurrentTheme.AccentColor)
	t.TextStyle = CurrentTheme.NormalText
	return t
}

func NewStyledList(title string, items []list.Item) list.Model {
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = CurrentTheme.BoldText.Copy().Foreground(CurrentTheme.AccentColor)
	return l
}

func StyledButton(text string, focused bool) string {
	if focused {
		return CurrentTheme.Button.Copy().Background(CurrentTheme.AccentColor).Render(text)
	}
	return CurrentTheme.Button.Render(text)
}
