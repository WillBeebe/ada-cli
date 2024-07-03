package cmd

import (
	"os"

	"github.com/container-labs/ada/internal/ada"
	"github.com/container-labs/ada/internal/golang"
	"github.com/container-labs/ada/internal/nodejs"
	"github.com/container-labs/ada/internal/python"
	"github.com/container-labs/ada/internal/terraform"
	"github.com/spf13/cobra"
)

var dependencyCmd = &cobra.Command{
	Use:   "dep",
	Short: "manage app dependencies",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var dependencyAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add a dependency to the app",
	Run: func(cmd *cobra.Command, args []string) {
		adaFile := ada.Load()

		dep := args[0]

		var err error
		switch adaFile.Type {
		case "python":
			err = python.AddDependency(dep)
		case "terraform":
			err = terraform.AddDependency(dep)
		case "nodejs":
			err = nodejs.AddDependency(dep)
		case "go":
			err = golang.AddDependency(dep)
		default:
			logger.Info("not implemented")
		}

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
