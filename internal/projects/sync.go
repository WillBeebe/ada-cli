package projects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/unidoc/unidoc/pdf/extractor"
	unidoc "github.com/unidoc/unidoc/pdf/model"
)

type FileObject struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func (c *APIClient) AddFileToProject(projectID int, data map[string]string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/projects/%d/files", c.BaseURL, projectID)
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

func readFile(filePath string) (*FileObject, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %s", filePath)
	}

	if strings.ToLower(filepath.Ext(filePath)) == ".pdf" {
		f, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return nil, err
		}
		defer f.Close()
		content := []byte{}

		pdfReader, err := unidoc.NewPdfReader(f)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return nil, err
		}
		pages, _ := pdfReader.GetNumPages()

		for i := 0; i < pages; i++ {
			page, err := pdfReader.GetPage(i + 1)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			// 	page, err := pdfReader.GetPage(pageNum)
			// if err != nil {
			// 	return err
			// }

			ex, err := extractor.New(page)
			if err != nil {
				return nil, err
			}

			text, err := ex.ExtractText()
			if err != nil {
				return nil, err
			}
			pageBytes := []byte(text)
			content = append(content, pageBytes...)

			fmt.Printf("Page %d: %s\n", i+1, text)
		}

		return &FileObject{
			Path:    filePath,
			Content: string(content),
		}, nil
	}

	return &FileObject{
		Path:    filePath,
		Content: string(content),
	}, nil
}

func isAllowedFile(filePath string, allowedExtensions []string) bool {
	ext := filepath.Ext(filePath)
	for _, allowed := range allowedExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}

func readFilesRecursively(directory string, allowedExtensions, blacklist []string) []*FileObject {
	var fileObjects []*FileObject

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && contains(blacklist, info.Name()) {
			return filepath.SkipDir
		}

		if !info.IsDir() && !contains(blacklist, info.Name()) && isAllowedFile(path, allowedExtensions) {
			fileObject, err := readFile(path)
			if err == nil {
				fileObjects = append(fileObjects, fileObject)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	return fileObjects
}

func addFilesToProject(projectID int, fileObjects []*FileObject, client *APIClient, dryRun bool) {
	for _, fileObject := range fileObjects {
		fmt.Printf("file content is: %s \n", fileObject.Content)
		data := map[string]string{
			"name":    filepath.Base(fileObject.Path),
			"content": fileObject.Content,
			"path":    fileObject.Path,
			// "is_context_file": "true",
		}

		if dryRun {
			fmt.Printf("Dry Run: File '%s' would be added to the project.\n", data["path"])
		} else {
			response, err := client.AddFileToProject(projectID, data)
			if err != nil {
				fmt.Println("Error creating:", err)
			} else {
				fmt.Println(response)
				fmt.Printf("File '%s' added to the project.\n", data["path"])
			}
		}
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func Sync(directory string, projectID int, dryRun bool) error {
	// todo: move client to factory
	client := &APIClient{BaseURL: "http://0.0.0.0:8000"}
	allowedExtensions := []string{".go", ".py", ".yaml", ".js", ".json", ".pdf"}
	blacklist := []string{".git", ".gitkeep", "bin", "node_modules", ".venv"}

	fileObjects := readFilesRecursively(directory, allowedExtensions, blacklist)
	addFilesToProject(projectID, fileObjects, client, dryRun)

	return nil
}
