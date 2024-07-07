package chat

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
)

type message struct {
	content string
	isUser  bool
}

type model struct {
	lines       []string
	cursorX     int
	cursorY     int
	err         error
	viewport    viewport.Model
	messages    []message
	spinner     spinner.Model
	isWaiting   bool
	inputMode   bool
	currentLine int
}

// Add any model-specific methods here
