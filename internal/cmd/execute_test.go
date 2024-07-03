package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name    string
		command string
		wantErr bool
		wantOut string
	}{
		{"Invalid Command", "abcd1234xyz", true, ""},
		{"Valid Command", "echo hello", false, "hello\n"},
		{"Valid Command with double quotes", "echo -n \"hello\"", false, "hello"},
		{"Valid Command with single quotes", "echo -n 'hello'", false, "hello"},
		{"Valid Command with pipe", "echo 'foo\nbar' | grep 'ar'", false, "bar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &CommandOptions{Command: tt.command}

			// Capture standard output and error
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			rStdout, wStdout, _ := os.Pipe()
			rStderr, wStderr, _ := os.Pipe()
			os.Stdout = wStdout
			os.Stderr = wStderr

			defer func() {
				os.Stdout = oldStdout
				os.Stderr = oldStderr
				wStdout.Close()
				wStderr.Close()
			}()

			output, err := Execute(opts)
			wStdout.Close()
			wStderr.Close()

			var bufStdout, bufStderr bytes.Buffer
			if _, err := bufStdout.ReadFrom(rStdout); err != nil {
				t.Fatal(err)
			}
			if _, err := bufStderr.ReadFrom(rStderr); err != nil {
				t.Fatal(err)
			}

			// Check for error consistency
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute()\nGot error:%v\nWant error:%v", err, tt.wantErr)
			}

			// // Check for correct stdout output
			// // TODO: logger.Info() interferes with this test
			// if strings.TrimSpace(bufStdout.String()) != strings.TrimSpace(tt.wantOut) {
			// 	t.Errorf("Unexpected stdout output.\nGot:%v\nWant:%v", bufStdout.String(), tt.wantOut)
			// }

			// Check if the output matches the expected output
			if !tt.wantErr && strings.TrimSpace(string(output)) != strings.TrimSpace(tt.wantOut) {
				t.Errorf("Execute()\nGot:%v\nWant:%v", string(output), tt.wantOut)
			}
		})
	}
}
