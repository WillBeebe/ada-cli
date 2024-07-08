package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/container-labs/ada/internal/ada"
	"github.com/container-labs/ada/internal/common"
)

var logger = common.Logger()

type ChatService struct {
	client *ClientWithResponses
	config *ada.Config
}

type ResponseBody struct {
	Content string `json:"content"`
}

type PromptResponse struct {
	Response ResponseBody `json:"response"`
}

type ChatMessage struct {
	ID            int    `json:"id"`
	Role          string `json:"role"`
	Content       string `json:"content"`
	IsContextFile bool   `json:"is_context_file"`
	Model         string `json:"model"`
	Tokens        int    `json:"tokens"`
	IsToolMessage bool   `json:"is_tool_message"`
}

func NewChatService(baseURL string) (*ChatService, error) {
	client, err := NewClientWithResponses(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err)
	}

	config := ada.LoadConfig()
	if config == nil {
		return nil, fmt.Errorf("error loading global config")
	}

	return &ChatService{
		client: client,
		config: config,
	}, nil
}
func (s *ChatService) StartSession(ctx context.Context) error {
	// No session initialization needed for this API
	return nil
}

func (s *ChatService) SendMessage(ctx context.Context, prompt string) (string, error) {
	body := PostQuery{
		Prompt:    prompt,
		ProjectId: fmt.Sprintf("%d", s.config.CurrentProjectID),
	}

	resp, err := s.client.PromptPromptPostWithBody(ctx, "application/json", createJSONReader(body))
	if err != nil {
		return "", fmt.Errorf("error sending prompt: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var result PromptResponse
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}
	logger.Info(fmt.Sprintf("+%v", result))
	return result.Response.Content, nil

	// return "", fmt.Errorf("response not found in API result")
}

func (s *ChatService) GetChatHistory(ctx context.Context) ([]ChatMessage, error) {
	projectID := s.config.CurrentProjectID
	resp, err := s.client.ReadHistoryProjectsProjectIdHistoryGetWithResponse(ctx, fmt.Sprintf("%d", projectID))
	if err != nil {
		return nil, fmt.Errorf("error getting chat history: %v", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var chatHistory []ChatMessage
	err = json.Unmarshal(resp.Body, &chatHistory)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling chat history: %v", err)
	}

	return chatHistory, nil
}

func createJSONReader(v interface{}) io.Reader {
	jsonBytes, _ := json.Marshal(v)
	return strings.NewReader(string(jsonBytes))
}
