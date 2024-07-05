package cmd

import (
	"fmt"
	"os"

	"github.com/container-labs/ada/internal"
	"github.com/container-labs/ada/internal/ada"
	"github.com/container-labs/ada/internal/common"
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
		logger.Info(fmt.Sprintf("%v", adaFile))

		strategy, err := internal.LanguageFactory(adaFile)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		err = strategy.Install()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
