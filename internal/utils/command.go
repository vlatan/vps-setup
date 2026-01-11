package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Cmd struct {
	*exec.Cmd
}

// Command initializes a terminal command object
func Command(name string, args ...string) *Cmd {

	// Init the command
	cmd := exec.Command(name, args...)

	// Point the command Stdin, Stdout and Stderr streams
	// to the user OS corresponding streams
	cmd.Stdin = os.Stdin   // Allow interactive input
	cmd.Stdout = os.Stdout // Show the output
	cmd.Stderr = os.Stderr // Show errors directly to the terminal

	return &Cmd{Cmd: cmd}
}

// Run executes the command and captures stderr.
// If the command fails, returns an error with the captured stderr message.
func (c *Cmd) Run() error {
	var stderr bytes.Buffer
	c.Cmd.Stderr = &stderr

	err := c.Cmd.Run()
	if err == nil {
		return nil
	}

	if stderr.Len() == 0 {
		return err
	}

	return fmt.Errorf("%s\n%w", strings.TrimSpace(stderr.String()), err)
}
