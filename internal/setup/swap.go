package setup

import (
	"fmt"
	"path/filepath"

	"github.com/vlatan/vps-setup/internal/utils"
)

// ChangeSwappiness changes the system swappiness,
// which means changing the free memory treshold (%) at which
// the swap is going to begin to be utilized.
func (s *Setup) ChangeSwappiness() error {

	// Check if the swappiness is valid
	if s.Swappiness < 0 || s.Swappiness > 100 {
		return fmt.Errorf("invalid swappiness: %d", s.Swappiness)
	}

	fmt.Println("Setting up the swappiness...")

	// Make parent directories
	name := "sysctl.d/99-my-swappiness.conf"
	if err := s.Etc.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}

	// Write to the file
	data := fmt.Appendf([]byte{}, "vm.swappiness = %d\n", s.Swappiness)
	if err := s.Etc.WriteFile(name, data, 0644); err != nil {
		return err
	}

	// Load our config
	confPath := filepath.Join(s.Etc.Name(), name)
	cmd := utils.Command("sysctl", "-p", confPath)
	return cmd.Run()
}
