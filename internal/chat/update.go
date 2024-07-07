package chat

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.inputMode {
			return m.handleInputMode(msg)
		}
		return m.handleViewMode(msg)

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 6 // Reserve space for input
		if !m.inputMode {
			m.viewport.SetContent(m.formatMessages())
		}
		return m, nil

	case spinner.TickMsg:
		if m.isWaiting {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case errorMsg:
		logger.Error(fmt.Sprintf("Error: %v", msg.err))
		return m, tea.Quit

	case aiResponseMsg:
		logger.Debug("Received AI response in Bubble Tea model")
		m.isWaiting = false
		m.messages = append(m.messages, message{content: string(msg), isUser: false})
		m.viewport.SetContent(m.formatMessages())
		m.viewport.GotoBottom()
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) handleInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		if m.lines[0] == "bye" {
			logger.Debug("User entered 'bye', quitting")
			return m, tea.Quit
		}
		logger.Debug("User finished input, sending message")
		m.inputMode = false
		m.isWaiting = true
		userMessage := strings.Join(m.lines, "\n")
		m.messages = append(m.messages, message{content: userMessage, isUser: true})
		m.viewport.SetContent(m.formatMessages())
		m.viewport.GotoBottom()
		m.lines = []string{""}
		m.cursorX = 0
		m.cursorY = 0
		return m, tea.Batch(sendMessageCmd(userMessage), spinner.Tick)

	case tea.KeyEnter:
		m.lines = append(m.lines, "")
		m.cursorY++
		m.cursorX = 0

	case tea.KeyBackspace:
		if m.cursorX > 0 {
			m.lines[m.cursorY] = m.lines[m.cursorY][:m.cursorX-1] + m.lines[m.cursorY][m.cursorX:]
			m.cursorX--
		} else if m.cursorY > 0 {
			m.cursorY--
			m.cursorX = len(m.lines[m.cursorY])
			m.lines[m.cursorY] += m.lines[m.cursorY+1]
			m.lines = append(m.lines[:m.cursorY+1], m.lines[m.cursorY+2:]...)
		}

	case tea.KeySpace:
		m.lines[m.cursorY] = m.lines[m.cursorY][:m.cursorX] + " " + m.lines[m.cursorY][m.cursorX:]
		m.cursorX++

	case tea.KeyRunes:
		m.lines[m.cursorY] = m.lines[m.cursorY][:m.cursorX] + string(msg.Runes) + m.lines[m.cursorY][m.cursorX:]
		m.cursorX += len(msg.Runes)

	case tea.KeyLeft:
		if m.cursorX > 0 {
			m.cursorX--
		} else if m.cursorY > 0 {
			m.cursorY--
			m.cursorX = len(m.lines[m.cursorY])
		}

	case tea.KeyRight:
		if m.cursorX < len(m.lines[m.cursorY]) {
			m.cursorX++
		} else if m.cursorY < len(m.lines)-1 {
			m.cursorY++
			m.cursorX = 0
		}

	case tea.KeyUp:
		if m.cursorY > 0 {
			m.cursorY--
			if m.cursorX > len(m.lines[m.cursorY]) {
				m.cursorX = len(m.lines[m.cursorY])
			}
		}

	case tea.KeyDown:
		if m.cursorY < len(m.lines)-1 {
			m.cursorY++
			if m.cursorX > len(m.lines[m.cursorY]) {
				m.cursorX = len(m.lines[m.cursorY])
			}
		}
	}

	return m, nil
}

func (m model) handleViewMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		m.inputMode = true
	default:
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

func sendMessageCmd(message string) tea.Cmd {
	return func() tea.Msg {
		logger.Debug(fmt.Sprintf("Sending message to channel: %s", message))
		messageChan <- message
		return nil
	}
}
