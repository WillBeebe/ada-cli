package cmd

import (
	"fmt"

	"github.com/container-labs/ada/internal/ada"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		globalConfig := ada.LoadConfig()
		logger.Debug(fmt.Sprintf("%+v\n", globalConfig))
	},
}
