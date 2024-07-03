package nodejs

import "github.com/container-labs/ada/internal/cmd"

func Install() error {
	_, err := cmd.Execute(&cmd.CommandOptions{
		Command: "npm install",
	})

	return err
}
