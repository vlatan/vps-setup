package utils

import (
	"os"
	"path/filepath"
)

// WriteFile writes data to file, creates all dirs in path
func WriteFile(root *os.Root, name string, data []byte) error {
	if err := root.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}

	return root.WriteFile(name, data, 0644)
}
