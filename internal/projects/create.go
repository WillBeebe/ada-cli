package projects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
)

func (c *APIClient) CreateProject(data map[string]string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/projects", c.BaseURL)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type Project struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Path          string `json:"path"`
	Provider      string `json:"provider"`
	ProviderModel string `json:"provider_model"`
}

type CreateResponse struct {
	Response Project `json:"response"`
}

func Create(dryrun bool) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	client := &APIClient{BaseURL: "http://0.0.0.0:8000"}
	data := map[string]string{
		"name": fmt.Sprintf("cli-%s", randSeq(4)),
		"path": currentDir,
		// "provider":       "anthropic",
		// "provider_model": "claude-3-opus-20240229",
		"provider":       "groq",
		"provider_model": "llama3-8b-8192",
	}
	projectID := 74
	if dryrun {
		fmt.Println("Dry Run Create Project:", data)
	} else {
		response, err := client.CreateProject(data)
		if err != nil {
			fmt.Println("Error creating:", err)
		} else {
			var projectResponse CreateResponse
			bytes, _ := json.Marshal(response)
			json.Unmarshal(bytes, &projectResponse)
			fmt.Printf("created project: %s\n", projectResponse.Response.Name)
			projectID = projectResponse.Response.ID
		}
	}

	err = Sync(currentDir, projectID, dryrun)
	if err != nil {
		fmt.Println("Error creating:", err)
	}
	return nil
}
