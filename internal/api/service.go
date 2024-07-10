package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Service struct {
	baseURL  string
	adaToken string
	client   *Client
}

type Project struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Path          string `json:"path"`
	Provider      string `json:"provider"`
	ProviderModel string `json:"provider_model"`
}

func NewService(baseURL, adaToken string) (*Service, error) {
	client, err := NewClient(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	return &Service{
		baseURL:  baseURL,
		adaToken: adaToken,
		client:   client,
	}, nil
}

func (s *Service) addCustomHeader(ctx context.Context, req *http.Request) error {
	req.Header.Set("x-ada-token", s.adaToken)
	return nil
}

func (s *Service) ListProjects(ctx context.Context, pageToken string) ([]Project, error) {
	s.client.RequestEditors = []RequestEditorFn{s.addCustomHeader}

	resp, err := s.client.ListProjectsGet(ctx)
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
	s.client.RequestEditors = []RequestEditorFn{s.addCustomHeader}

	resp, err := s.client.ReadProjectsProjectIdGet(ctx, projectID)
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

	s.client.RequestEditors = []RequestEditorFn{s.addCustomHeader}

	resp, err := s.client.CreateProjectsPost(ctx, createProjectBody)
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
