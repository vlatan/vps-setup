package utils

import (
	"os"
	"os/exec"
)

// Command runs a terminal command
func Command(name string, args ...string) *exec.Cmd {

	// Init the command
	cmd := exec.Command(name, args...)

	// Point the command Stdin, Stdout and Stderr streams
	// to the user OS corresponding streams
	cmd.Stdin = os.Stdin   // Allow interactive input
	cmd.Stdout = os.Stdout // Show the output
	cmd.Stderr = os.Stderr // Show errors directly to the terminal

	return cmd
}
