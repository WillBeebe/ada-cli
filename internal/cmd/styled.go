package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// var logger = common.Logger()

// Define some styles
var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(subtle).
		String()

	commandStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			PaddingLeft(1)

	outputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#333333")).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)
)

// type CommandOptions struct {
// 	Command   string
// 	Directory string
// }

func StyledExecute(opts *CommandOptions) ([]byte, error) {
	logger.Info(fmt.Sprintf("executing `%s`", opts.Command))

	// Print styled command
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Foreground(special).Render("$"),
		commandStyle.Render(opts.Command),
	))

	cmd := exec.Command("bash", "-c", opts.Command)

	var outb, errb bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &outb)
	cmd.Stderr = io.MultiWriter(os.Stderr, &errb)

	if opts.Directory != "" {
		cmd.Dir = opts.Directory
	}

	err := cmd.Run()

	// Style and print output
	output := strings.TrimSpace(outb.String())
	if output != "" {
		fmt.Println(outputStyle.Render(output))
	}

	if err != nil {
		errorOutput := strings.TrimSpace(errb.String())
		if errorOutput != "" {
			fmt.Println(errorStyle.Render("Error: " + errorOutput))
		}
		return errb.Bytes(), fmt.Errorf(outb.String())
	}

	return outb.Bytes(), nil
}
