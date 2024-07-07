package golang

import (
	"github.com/container-labs/ada/internal/cmd"
)

func Install() error {
	// install project dependencies
	_, err := cmd.StyledExecute(&cmd.CommandOptions{
		Command: "go get",
	})

	return err
}
