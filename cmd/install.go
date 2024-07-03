package cmd

import "github.com/spf13/cobra"

var installCmd = &cobra.Command{
	Use: "install",
	// Short: "",
	// Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(installCmd)

}
