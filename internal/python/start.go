package python

import "github.com/container-labs/ada/internal/cmd"

func Start() error {
	// install project dependencies
	_, err := cmd.StyledExecute(&cmd.CommandOptions{
		Command: "poetry run python src/main.py",
	})

	return err
}
