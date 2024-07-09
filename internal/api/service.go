package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Service struct {
	baseURL string
}

func NewService(baseURL string) *Service {
	return &Service{
		baseURL: baseURL,
	}
}

type Project struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Path          string `json:"path"`
	Provider      string `json:"provider"`
	ProviderModel string `json:"provider_model"`
}

func (s *Service) ListProjects(ctx context.Context, pageToken string) ([]Project, error) {
	client, err := NewClient(s.baseURL)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	resp, err := client.ListProjectsGet(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing projects: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var projects []Project
	err = decodeBody(resp, &projects)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return projects, nil
}

func (s *Service) GetProject(ctx context.Context, projectID int) (*Project, error) {
	client, err := NewClient(s.baseURL)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	resp, err := client.ReadProjectsProjectIdGet(ctx, strconv.Itoa(projectID))
	if err != nil {
		return nil, fmt.Errorf("error getting project: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var project Project
	err = decodeBody(resp, &project)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &project, nil
}

func (s *Service) CreateProject(ctx context.Context, name, path, provider, providerModel string) (*Project, error) {
	createProjectBody := CreateProject{
		Name:          name,
		Path:          path,
		Provider:      provider,
		ProviderModel: providerModel,
	}

	client, err := NewClient(s.baseURL)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	resp, err := client.CreateProjectsPost(ctx, createProjectBody)
	if err != nil {
		return nil, fmt.Errorf("error creating project: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var project Project
	err = decodeBody(resp, &project)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &project, nil
}

func decodeBody(resp *http.Response, v interface{}) error {
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(v)
}

// You can add more methods here as needed, such as UpdateProject, DeleteProject, etc.
