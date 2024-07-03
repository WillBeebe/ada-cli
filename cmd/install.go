package cmd

import (
	"os"

	"github.com/container-labs/ada/internal/ada"
	"github.com/container-labs/ada/internal/common"
	"github.com/container-labs/ada/internal/golang"
	"github.com/container-labs/ada/internal/nodejs"
	"github.com/container-labs/ada/internal/python"
	"github.com/container-labs/ada/internal/terraform"
	"github.com/spf13/cobra"
)

var logger = common.Logger()

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install the app in the cwd on the local machine",
	Run: func(cmd *cobra.Command, args []string) {
		adaFile := ada.Load()

		if adaFile.Install.Disabled {
			logger.Info("Skipping install")
			return
		}

		var err error
		switch adaFile.Type {
		case "python":
			err = python.Install()
		case "terraform":
			err = terraform.Install()
		case "nodejs":
			err = nodejs.Install()
		case "go":
			err = golang.Install()
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
	rootCmd.AddCommand(installCmd)
}
