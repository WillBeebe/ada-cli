package cmd

import (
	"os"

	"github.com/container-labs/ada/internal/ada"
	"github.com/container-labs/ada/internal/container"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(containerCmd)
	containerCmd.AddCommand(containerBuildCmd)
	containerCmd.AddCommand(containerRunCmd)
	containerCmd.AddCommand(containerPushCmd)
	containerPushCmd.Flags().BoolP("artifactory", "a", false, "Push image to artifactory")
	containerRunCmd.Flags().BoolP("bash", "b", false, "Start the container with bash")
}

var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "Build and run Docker containers for your projects",
}
var containerBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the container for the project in your cwd",
	Run: func(cmd *cobra.Command, args []string) {
		adaFile := ada.Load()
		err := container.Build(adaFile)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	},
}
var containerRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the container for the project in your cwd",
	Run: func(cmd *cobra.Command, args []string) {
		adaFile := ada.Load()
		runBash, _ := cmd.Flags().GetBool("bash")
		err := container.Run(adaFile, runBash)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	},
}
var containerPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push the container to GAR. Use --artifactory to push to Artifactory instead",
	Run: func(cmd *cobra.Command, args []string) {
		adaFile := ada.Load()
		pushToArtifactory, _ := cmd.Flags().GetBool("artifactory")
		var err error
		if pushToArtifactory {
			err = container.PushArtifactory(adaFile)
		} else {
			err = container.Push(adaFile)
		}
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	},
}
