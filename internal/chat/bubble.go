package chat

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/container-labs/ada/internal/api"
)

func initialModel(history []api.ChatMessage) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	v := viewport.New(80, 20)
	v.KeyMap = viewport.KeyMap{
		PageDown: key.NewBinding(
			key.WithKeys("pgdown"),
			key.WithHelp("pgdn", "page down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("pgup", "page up"),
		),
		HalfPageUp: key.NewBinding(
			key.WithKeys("ctrl+u"),
			key.WithHelp("ctrl+u", "half page up"),
		),
		HalfPageDown: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "half page down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
	}

	return model{
		lines:     []string{""},
		cursorX:   0,
		cursorY:   0,
		viewport:  v,
		messages:  formatChatHistory(history),
		spinner:   s,
		isWaiting: false,
		inputMode: true,
	}
}

func formatChatHistory(history []api.ChatMessage) []message {
	var formattedMessages []message
	for _, msg := range history {
		formattedMsg := message{
			isUser:  msg.Role == "user",
			content: msg.Content,
		}
		formattedMessages = append(formattedMessages, formattedMsg)
	}
	return formattedMessages
}

func (m model) Init() tea.Cmd {
	return tea.Batch(spinner.Tick, viewport.Sync(m.viewport))
}

func BubblePrompt() (string, error) {
	var initial []api.ChatMessage
	p := tea.NewProgram(initialModel(initial))
	m, err := p.Run()
	if err != nil {
		return "", err
	}
	return strings.Join(m.(model).lines, "\n"), nil
}
