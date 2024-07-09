package cmd

import (
	"fmt"
	"os"

	"github.com/container-labs/ada/internal"
	"github.com/container-labs/ada/internal/ada"
	"github.com/container-labs/ada/internal/api"

	"github.com/container-labs/ada/internal/create"
	"github.com/spf13/cobra"
)

var listProjectsCmd = &cobra.Command{
	Use:   "projects",
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

var projectsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Run: func(cmd *cobra.Command, args []string) {
		newProject, err := create.CreateProject()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating project: %v\n", err)
			os.Exit(1)
		}

		service := api.NewService("http://localhost:8000")

		createdProject, err := service.CreateProject(cmd.Context(), newProject.Name, newProject.Path, newProject.Provider, newProject.ProviderModel)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating project on server: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Project created successfully: %s (ID: %d)\n", createdProject.Name, createdProject.ID)

		// Load the current config
		config := ada.LoadConfig()
		if config == nil {
			fmt.Println("Error loading config")
			os.Exit(1)
		}

		// Update the config with the new project
		config.CurrentProject = createdProject.Name
		config.CurrentProjectID = createdProject.ID

		// Save the updated config
		err = ada.SaveConfig(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("New project set as current project in config file")
	},
}

func init() {
	rootCmd.AddCommand(listProjectsCmd)
	listProjectsCmd.AddCommand(projectsCreateCmd)
}
