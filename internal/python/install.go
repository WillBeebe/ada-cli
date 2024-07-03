package python

import (
	"github.com/container-labs/ada/internal/cmd"
)

func Install() error {
	// install project dependencies
	_, err := cmd.Execute(&cmd.CommandOptions{
		Command: "poetry install --all-extras",
	})

	return err
}

func InstallDeps() error {
	// install external dependencies
	_, err := cmd.Execute(&cmd.CommandOptions{
		Command: "pip3 install poetry pre-commit",
	})

	if err != nil {
		return err
	}

	// install pre-commit hooks
	_, err = cmd.Execute(&cmd.CommandOptions{
		Command: "poetry run pre-commit install",
	})

	return err
}
