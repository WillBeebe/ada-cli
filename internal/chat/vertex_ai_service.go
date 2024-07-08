package chat

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/vertexai/genai"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type VertexAIService struct {
	projectID   string
	region      string
	modelName   string
	client      *genai.Client
	chatSession *genai.ChatSession
}

func NewVertexAIService(projectID, region, modelName string) *VertexAIService {
	return &VertexAIService{
		projectID: projectID,
		region:    region,
		modelName: modelName,
	}
}

func (s *VertexAIService) StartSession(ctx context.Context) error {
	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return fmt.Errorf("error loading default credentials: %v", err)
	}

	opts := []option.ClientOption{
		option.WithCredentials(creds),
	}

	client, err := genai.NewClient(ctx, s.projectID, s.region, opts...)
	if err != nil {
		return fmt.Errorf("error creating client: %v", err)
	}

	s.client = client
	gemini := client.GenerativeModel(s.modelName)
	s.chatSession = gemini.StartChat()

	return nil
}

func (s *VertexAIService) SendMessage(ctx context.Context, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	response, err := s.chatSession.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("error sending message: %v", err)
	}

	if len(response.Candidates) == 0 {
		return "I am unable to help with that, please try again", nil
	}

	responseString := ""
	for _, part := range response.Candidates[0].Content.Parts {
		responseString += fmt.Sprintf("%v", part)
	}

	return responseString, nil
}
