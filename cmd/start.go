package cmd

import (
	"os"

	"github.com/container-labs/ada/internal"
	"github.com/container-labs/ada/internal/ada"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the application",
	Run: func(cmd *cobra.Command, args []string) {
		adaFile := ada.Load()

		strategy, err := internal.LanguageFactory(adaFile)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		err = strategy.Start()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
