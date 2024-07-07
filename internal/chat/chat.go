package chat

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/vertexai/genai"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/container-labs/ada/internal/common"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var logger = common.Logger()

type aiResponseMsg string
type errorMsg struct{ err error }

const (
	apiTimeout = 30 * time.Second
)

func Chat() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Info("Starting Chat function")

	err := TestCredentials(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("Credentials test failed: %v", err))
		fmt.Println("Please ensure your Google Cloud credentials are set up correctly:")
		fmt.Println("1. Run 'gcloud auth application-default login' in your terminal")
		fmt.Println("2. Verify that you're logged in to the correct account with 'gcloud auth list'")
		fmt.Println("3. Check your current project with 'gcloud config list project'")
		fmt.Println("4. Ensure the project ID in the code matches your active project")
		os.Exit(1)
	}

	logger.Info("Credentials test passed, starting chat session")

	chatSession, err := startSession(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to start session: %v", err))
		os.Exit(1)
	}

	logger.Info("Chat session started successfully")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	p := tea.NewProgram(initialModel())

	done := make(chan bool, 1)

	go func() {
		logger.Info("Starting message handling goroutine")
		for {
			select {
			case msg := <-messageChan:
				logger.Debug(fmt.Sprintf("Received message from channel: %s", msg))
				response, err := sendMessage(ctx, chatSession, msg)
				if err != nil {
					logger.Error(fmt.Sprintf("Error sending message: %v", err))
					p.Send(errorMsg{err})
				} else {
					logger.Debug(fmt.Sprintf("Received response from AI: %s", response))
					p.Send(aiResponseMsg(response))
				}
			case <-sigs:
				logger.Info("Received interrupt signal, quitting")
				p.Send(tea.Quit())
				return
			case <-done:
				logger.Info("Received done signal, exiting goroutine")
				return
			case <-ctx.Done():
				logger.Info("Context cancelled, exiting goroutine")
				p.Send(errorMsg{ctx.Err()})
				return
			}
		}
	}()

	logger.Info("Starting Bubble Tea program")
	if _, err := p.Run(); err != nil {
		logger.Error(fmt.Sprintf("Error running program: %v", err))
		os.Exit(1)
	}

	logger.Info("Bubble Tea program finished, signaling goroutine to exit")
	done <- true
}

var messageChan = make(chan string)

func startSession(ctx context.Context) (*genai.ChatSession, error) {
	var projectId = "ada-test-1234" // Ensure this matches your actual project ID
	var region = "us-central1"
	var modelName = "gemini-1.5-flash-001"

	// Log the current working directory and GOOGLE_APPLICATION_CREDENTIALS value
	cwd, _ := os.Getwd()
	logger.Debug(fmt.Sprintf("Current working directory: %s", cwd))
	logger.Debug(fmt.Sprintf("GOOGLE_APPLICATION_CREDENTIALS: %s", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))

	// Attempt to load default credentials
	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load default credentials: %v", err))
		return nil, fmt.Errorf("error loading default credentials: %v", err)
	}
	logger.Debug(fmt.Sprintf("Loaded credentials for project: %s", creds.ProjectID))

	// Create client options with default credentials
	opts := []option.ClientOption{
		option.WithCredentials(creds),
	}

	client, err := genai.NewClient(ctx, projectId, region, opts...)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create client: %v", err))
		return nil, fmt.Errorf("error creating client: %v", err)
	}

	gemini := client.GenerativeModel(modelName)
	chat := gemini.StartChat()

	logger.Info("Successfully started chat session")
	return chat, nil
}

const apiCallTimeout = 60 * time.Second

func sendMessage(ctx context.Context, session *genai.ChatSession, prompt string) (string, error) {
	logger.Debug(fmt.Sprintf("Preparing to send prompt: %s", prompt))

	// Create a new context with a timeout
	ctx, cancel := context.WithTimeout(ctx, apiCallTimeout)
	defer cancel()

	logger.Debug("Sending message to API...")
	startTime := time.Now()

	responseChan := make(chan *genai.GenerateContentResponse)
	errChan := make(chan error)

	go func() {
		response, err := session.SendMessage(ctx, genai.Text(prompt))
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	}()

	select {
	case response := <-responseChan:
		logger.Debug(fmt.Sprintf("Received response from API after %v", time.Since(startTime)))
		if len(response.Candidates) == 0 {
			logger.Warn("Received empty response from API")
			return "I am unable to help with that, please try again", nil
		}

		responseString := ""
		for _, part := range response.Candidates[0].Content.Parts {
			responseString += fmt.Sprintf("%v", part)
		}

		logger.Debug(fmt.Sprintf("Processed response: %s", responseString))
		return responseString, nil

	case err := <-errChan:
		logger.Error(fmt.Sprintf("Error sending message after %v: %v", time.Since(startTime), err))
		return "", fmt.Errorf("error sending message: %v", err)

	case <-ctx.Done():
		logger.Error(fmt.Sprintf("Request timed out after %v", apiCallTimeout))
		return "", fmt.Errorf("request timed out after %v", apiCallTimeout)
	}
}

func sendMessageWithTimeout(session *genai.ChatSession, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	responseChan := make(chan string)
	errChan := make(chan error)

	go func() {
		response, err := sendMessage(ctx, session, prompt)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	}()

	select {
	case response := <-responseChan:
		return response, nil
	case err := <-errChan:
		return "", err
	case <-ctx.Done():
		return "", fmt.Errorf("request timed out after %v", apiTimeout)
	}
}

func printModelResponse(response string) {
	maxWidth := width - paddingWidth
	fmt.Println(lipgloss.JoinHorizontal(
		lipgloss.Top,
		ChatContentStyle.Width(maxWidth-5).Align(lipgloss.Left).Render(response)),
		ChatPadding.Width(paddingWidth).Render(""),
	)
}

func printUserPrompt(prompt string) {
	w := lipgloss.Width
	maxWidth := width - paddingWidth
	chatPrompt := ChatUserContentStyle.UnsetWidth().Align(lipgloss.Right).Render(prompt)
	promptWidth := w(chatPrompt)
	if promptWidth > maxWidth {
		chatPrompt = ChatUserContentStyle.Width(maxWidth).Align(lipgloss.Right).Render(prompt)
	}
	paddingNeeded := width - w(chatPrompt)
	fmt.Println(lipgloss.JoinHorizontal(
		lipgloss.Top,
		ChatPadding.Width(paddingNeeded).Render(""),
		chatPrompt),
	)
}

// Add this function to test the credentials separately
func TestCredentials(ctx context.Context) error {
	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return fmt.Errorf("failed to load default credentials: %v", err)
	}
	token, err := creds.TokenSource.Token()
	if err != nil {
		return fmt.Errorf("failed to get token: %v", err)
	}
	logger.Info(fmt.Sprintf("Successfully obtained token for project: %s", creds.ProjectID))
	logger.Debug(fmt.Sprintf("Token expires at: %v", token.Expiry))
	return nil
}
