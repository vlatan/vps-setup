package utils

import (
	"os"
	"os/exec"
)

// RunCommand runs a terminal command
func RunCommand(name string, args ...string) error {

	// Init the command
	cmd := exec.Command(name, args...)

	// Point the command Stdin, Stdout and Stderr streams
	// to the user OS corresponding streams
	cmd.Stdin = os.Stdin   // Allow interactive input
	cmd.Stdout = os.Stdout // Show the output
	cmd.Stderr = os.Stderr // Show errors directly to the terminal

	return cmd.Run()
}
