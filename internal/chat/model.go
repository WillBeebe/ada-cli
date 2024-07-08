package chat

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/container-labs/ada/internal/api"
)

type message struct {
	content string
	isUser  bool
	time    time.Time
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

type AIService interface {
	StartSession(ctx context.Context) error
	SendMessage(ctx context.Context, prompt string) (string, error)
	GetChatHistory(ctx context.Context) ([]api.ChatMessage, error)
}

// type ChatMessage struct {
// 	ID            int    `json:"id"`
// 	Role          string `json:"role"`
// 	Content       string `json:"content"`
// 	IsContextFile bool   `json:"is_context_file"`
// 	Model         string `json:"model"`
// 	Tokens        int    `json:"tokens"`
// 	IsToolMessage bool   `json:"is_tool_message"`
// }
