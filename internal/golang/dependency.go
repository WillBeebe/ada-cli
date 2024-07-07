package golang

import (
	"fmt"

	"github.com/container-labs/ada/internal/cmd"
)

func AddDependency(dep string) error {
	_, err := cmd.StyledExecute(&cmd.CommandOptions{
		Command: fmt.Sprintf("go get %s", dep),
	})

	return err
}
