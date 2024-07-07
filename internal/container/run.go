package container

import (
	"fmt"

	"github.com/container-labs/ada/internal/ada"
	"github.com/container-labs/ada/internal/cmd"
)

// .PHONY: docker-run-bash
// docker-run-bash: generate-local-env ## run the app in a docker container but drop into a bash shell
//
//  docker run -it \
//    -p 8000:8000 \
//    --env-file=.env \
//    -v$${HOME}/.config/gcloud:/root/.config/gcloud \
//    $(APP_NAME):$(VERSION) bash
func Run(AdaFile *ada.AdaFile, runBash bool) error {
	// err := k8s.GenerateLocal(AdaFile)
	// if err != nil {
	// 	return err
	// }
	versionTag := ada.GetGitVersionTag()
	bashCmd := ""
	if runBash {
		bashCmd = "bash"
	}
	command := fmt.Sprintf("docker run -it -p %d:%d --env-file=.env -v$${HOME}/.config/gcloud:/root/.config/gcloud %s:%s %s",
		AdaFile.Port, AdaFile.Port, AdaFile.Name, versionTag, bashCmd)
	_, err := cmd.Execute(&cmd.CommandOptions{
		Command: command,
	})
	return err
}
