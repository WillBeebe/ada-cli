package cmd

import (
	"github.com/container-labs/ada/internal/projects"
	"github.com/spf13/cobra"
)

var graphCmd = &cobra.Command{
	Use: "graph",
	// Short: "",
	// Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projects.Graph(args[0])
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
}
