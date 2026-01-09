package utils

import (
	"os"
	"os/exec"
)

// RunCommand runs a terminal command
func RunCommand(name string, args ...string) error {

	// Init the command
	cmd := exec.Command(name, args...)

	// Point Stdin and Stderr streams to terminal
	cmd.Stdin = os.Stdin   // Allow interactive input
	cmd.Stderr = os.Stderr // Show errors directly to the terminal

	// Run the command
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
