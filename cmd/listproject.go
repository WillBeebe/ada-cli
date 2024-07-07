package cmd

import (
	"fmt"
	"os"

	"github.com/container-labs/ada/internal"
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

		selectedProject, err := internal.SelectProject(projects)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error selecting project: %v\n", err)
			os.Exit(1)
		}

		if selectedProject != nil {
			fmt.Printf("You selected project: %s (ID: %d)\n", selectedProject.Name, selectedProject.ID)
			// Here you can add further actions with the selected project
		} else {
			fmt.Println("No project selected")
		}
	},
}

func init() {
	rootCmd.AddCommand(listProjectsCmd)
}
