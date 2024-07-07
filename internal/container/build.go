package container

import (
	"fmt"

	"github.com/container-labs/ada/internal/ada"
	"github.com/container-labs/ada/internal/cmd"
	"github.com/container-labs/ada/internal/common"
)

var logger = common.Logger()

func Build(adaFile *ada.AdaFile) error {
	// right now this is only used by some frontend apps in dpe
	// so it's an opt-in feature, not on by default
	// if adaFile.Container.Build.GenerateEnv {
	// 	// err := k8s.GenerateLocal(adaFile)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// }
	versionTag := ada.GetGitShortSHA()
	command := fmt.Sprintf("docker build -t %s:%s .", adaFile.Name, versionTag)
	_, err := cmd.Execute(&cmd.CommandOptions{
		Command: command,
	})
	return err
}
