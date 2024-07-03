package cmd

import (
	"github.com/container-labs/ada/internal/projects"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use: "create",
	// Short: "",
	// Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := projects.Create(dryRun)
		if err != nil {
			// do better
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Perform a dry run without adding files")
}
