package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/container-labs/ada/internal/common"
)

var logger = common.Logger()

type CommandOptions struct {
	Command   string
	Directory string
}

// possibly move to something like this to stream command output
// https://stackoverflow.com/questions/37091316/how-to-get-the-realtime-output-for-a-shell-command-in-golang
func Execute(opts *CommandOptions) ([]byte, error) {
	logger.Info(fmt.Sprintf("executing `%s`", opts.Command))
	cmd := exec.Command("bash", "-c", opts.Command)

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Stdout = io.MultiWriter(os.Stdout, &outb)
	cmd.Stderr = io.MultiWriter(os.Stderr, &errb)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	if opts.Directory != "" {
		cmd.Dir = opts.Directory
	}

	err := cmd.Run()

	// output, err := cmd.Output()
	// output, err := cmd.CombinedOutput()
	if err != nil {
		// fmt.Println("error error error")
		// fmt.Println(outb.String())
		// return errb.Bytes(), err

		// return append(errb.Bytes(), outb.Bytes()...), err
		return errb.Bytes(), fmt.Errorf(outb.String())
	}

	return outb.Bytes(), nil
}
