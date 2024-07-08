package cmd

import (
	"fmt"
	"os"

	"github.com/container-labs/ada/internal"
	"github.com/container-labs/ada/internal/ada"
	"github.com/container-labs/ada/internal/api"
	"github.com/spf13/cobra"
)

var listProjectsCmd = &cobra.Command{
	Use:   "list-projects",
	Short: "List all projects and select one",
	Run: func(cmd *cobra.Command, args []string) {
		service := api.NewService("http://localhost:8000")

		projects, err := service.ListProjects(cmd.Context(), "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing projects: %v\n", err)
			os.Exit(1)
		}

		// Load the current config
		config := ada.LoadConfig()
		if config == nil {
			fmt.Println("Error loading config")
			os.Exit(1)
		}

		selectedProject, err := internal.SelectProject(projects, config.CurrentProjectID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error selecting project: %v\n", err)
			os.Exit(1)
		}

		if selectedProject != nil {
			fmt.Printf("You selected project: %s (ID: %d)\n", selectedProject.Name, selectedProject.ID)

			// Update the config with the selected project
			config.CurrentProject = selectedProject.Name
			config.CurrentProjectID = selectedProject.ID

			// Save the updated config
			err = ada.SaveConfig(config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Project ID saved to config file")
		} else {
			fmt.Println("No project selected")
		}
	},
}

func init() {
	rootCmd.AddCommand(listProjectsCmd)
}
