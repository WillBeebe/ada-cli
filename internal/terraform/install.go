package terraform

import (
	"github.com/container-labs/ada/internal/cmd"
)

func Install() error {
	_, err := cmd.StyledExecute(&cmd.CommandOptions{
		Command: "terraform init",
	})

	return err
}
