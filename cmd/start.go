package cmd

import "github.com/spf13/cobra"

var startCmd = &cobra.Command{
	Use: "start",
	// Short: "",
	// Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
