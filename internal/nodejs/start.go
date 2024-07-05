package nodejs

import "github.com/container-labs/ada/internal/cmd"

func Start() error {
	_, err := cmd.Execute(&cmd.CommandOptions{
		Command: "npm start",
	})

	return err
}
