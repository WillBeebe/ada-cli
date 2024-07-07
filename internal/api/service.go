package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Service struct {
	client  *http.Client
	baseURL string
}

func NewService(baseURL string) *Service {
	return &Service{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// Add other relevant fields
}

// type ListProjectsResponse struct {
// 	Projects      []Project `json:"projects"`
// 	NextPageToken string    `json:"nextPageToken,omitempty"`
// }

func (s *Service) ListProjects(ctx context.Context, pageToken string) ([]Project, error) {
	url := fmt.Sprintf("%s/projects?pageToken=%s", s.baseURL, pageToken)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return projects, nil
}

func (s *Service) GetProject(ctx context.Context, projectID int) (*Project, error) {
	url := fmt.Sprintf("%s/projects/%d", s.baseURL, projectID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &project, nil
}

// You can add more methods here as needed, such as CreateProject, UpdateProject, DeleteProject, etc.
