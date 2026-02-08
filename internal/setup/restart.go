package setup

import (
	"fmt"
	"path/filepath"
)

// AutoRestart creates a new config file /etc/needrestart/conf.d/no-prompt.conf
// Adds $nrconf{restart} = 'a'; to that file.
// This will set services to automatically restart after update/upgrade.
func (s *Setup) AutoRestart() error {
	fmt.Println("Seting up services autorestart...")

	// Make parent directories
	name := "needrestart/conf.d/no-prompt.conf"
	if err := s.Etc.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}

	// Write to file
	data := []byte("$nrconf{restart} = 'a';\n")
	if err := s.Etc.WriteFile(name, data, 0644); err != nil {
		return err
	}

	return nil
}
