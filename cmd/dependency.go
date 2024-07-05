package cmd

import (
	"os"

	"github.com/container-labs/ada/internal"
	"github.com/container-labs/ada/internal/ada"
	"github.com/spf13/cobra"
)

var dependencyCmd = &cobra.Command{
	Use:   "dep",
	Short: "manage app dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		// Command logic here
	},
}

var dependencyAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add a dependency to the app",
	Run: func(cmd *cobra.Command, args []string) {
		adaFile := ada.Load()

		dep := args[0]

		strategy, err := internal.LanguageFactory(adaFile)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		err = strategy.AddDependency(dep)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(dependencyCmd)
	dependencyCmd.AddCommand(dependencyAddCmd)
}
