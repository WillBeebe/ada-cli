package python

import (
	"fmt"

	"github.com/container-labs/ada/internal/cmd"
)

func AddDependency(dep string) error {
	_, err := cmd.Execute(&cmd.CommandOptions{
		Command: fmt.Sprintf("poetry add %s", dep),
	})

	return err
}
