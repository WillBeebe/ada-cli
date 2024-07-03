package cmd

import (
	"github.com/container-labs/ada/internal/projects"
	"github.com/spf13/cobra"
)

var (
	projectID int
	// allowedExtensions []string
	// blacklist         []string
	dryRun bool
	// baseURL           string
)

// type FileObject struct {
// 	Path    string `json:"path"`
// 	Content string `json:"content"`
// }

// type APIClient struct {
// 	BaseURL string
// }

// func (c *APIClient) AddFileToProject(projectID int, data map[string]string) (map[string]interface{}, error) {
// 	url := fmt.Sprintf("%s/projects/%d/files", c.BaseURL, projectID)
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var result map[string]interface{}
// 	err = json.NewDecoder(resp.Body).Decode(&result)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

// func readFile(filePath string) (*FileObject, error) {
// 	content, err := ioutil.ReadFile(filePath)
// 	if err != nil {
// 		return nil, fmt.Errorf("Error reading file: %s", filePath)
// 	}

// 	return &FileObject{
// 		Path:    filePath,
// 		Content: string(content),
// 	}, nil
// }

// func isAllowedFile(filePath string, allowedExtensions []string) bool {
// 	ext := filepath.Ext(filePath)
// 	for _, allowed := range allowedExtensions {
// 		if ext == allowed {
// 			return true
// 		}
// 	}
// 	return false
// }

// func readFilesRecursively(directory string, allowedExtensions, blacklist []string) []*FileObject {
// 	var fileObjects []*FileObject

// 	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		if info.IsDir() && contains(blacklist, info.Name()) {
// 			return filepath.SkipDir
// 		}

// 		if !info.IsDir() && !contains(blacklist, info.Name()) && isAllowedFile(path, allowedExtensions) {
// 			fileObject, err := readFile(path)
// 			if err == nil {
// 				fileObjects = append(fileObjects, fileObject)
// 			}
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		os.Exit(1)
// 	}

// 	return fileObjects
// }

// func addFilesToProject(projectID int, fileObjects []*FileObject, client *APIClient, dryRun bool) {
// 	for _, fileObject := range fileObjects {
// 		data := map[string]string{
// 			"name":    filepath.Base(fileObject.Path),
// 			"content": fileObject.Content,
// 			"path":    fileObject.Path,
// 			// "is_context_file": "true",
// 		}

// 		if dryRun {
// 			fmt.Printf("Dry Run: File '%s' would be added to the project.\n", data["path"])
// 		} else {
// 			response, err := client.AddFileToProject(projectID, data)
// 			if err != nil {
// 				fmt.Println("Error creating:", err)
// 			} else {
// 				fmt.Println(response)
// 			}
// 		}
// 	}
// }

// func contains(slice []string, item string) bool {
// 	for _, s := range slice {
// 		if s == item {
// 			return true
// 		}
// 	}
// 	return false
// }

var syncCmd = &cobra.Command{
	Use:   "sync [directory]",
	Short: "Recursively list files in a directory and add them to a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := args[0]
		// client := &APIClient{BaseURL: baseURL}

		// fileObjects := readFilesRecursively(directory, allowedExtensions, blacklist)
		// addFilesToProject(projectID, fileObjects, client, dryRun)
		projects.Sync(directory, projectID, dryRun)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().IntVarP(&projectID, "project-id", "p", 0, "Project ID")
	syncCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Perform a dry run without adding files")

	syncCmd.MarkFlagRequired("project-id")
}
