package ada

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const fileName = "ada.yaml"

type Config struct {
	Name             string `yaml:"name"`
	Version          string `yaml:"version"`
	CurrentProject   string `yaml:"current_project"`
	CurrentProjectID int    `yaml:"current_project_id"`
}

func LoadConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return nil
	}

	config := Config{
		Name:             "Ada",
		Version:          "1.0.0",
		CurrentProject:   "",
		CurrentProjectID: 0,
	}

	filePath := filepath.Join(homeDir, fileName)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File does not exist, create it with the default configuration
		data, err := yaml.Marshal(&config)
		if err != nil {
			fmt.Println("Error marshaling YAML:", err)
			return nil
		}

		err = ioutil.WriteFile(filePath, data, 0644)
		if err != nil {
			fmt.Println("Error writing file:", err)
			return nil
		}
		fmt.Println("File created:", filePath)
	} else {
		// File exists, read it
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return nil
		}

		err = yaml.Unmarshal(data, &config)
		if err != nil {
			fmt.Println("Error unmarshaling YAML:", err)
			return nil
		}
	}

	return &config
}

func SaveConfig(config *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %v", err)
	}

	filePath := filepath.Join(homeDir, fileName)

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshaling YAML: %v", err)
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}
