package setup

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/vlatan/vps-setup/internal/utils"
)

// ChangeSwappiness changes the system swappiness,
// which means changing the free memory treshold (%) at which
// the swap is going to begin to be utilized.
func (s *Setup) ChangeSwappiness() error {

	// Helper function to check if the swappiness value is valid
	valid := func(s string) bool {
		n, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		return n >= 0 && n <= 100
	}

	// Check if the swappiness is valid
	if valid(s.Swappiness) {
		return fmt.Errorf("invalid swappiness: %s", s.Swappiness)
	}

	fmt.Println("Setting up the swappiness...")

	// Write to the file
	name := "sysctl.d/99-my-swappiness.conf"
	data := fmt.Appendf([]byte{}, "vm.swappiness = %s\n", s.Swappiness)
	if err := utils.WriteFile(s.Etc, name, data); err != nil {
		return err
	}

	// Load our config
	confPath := filepath.Join(s.Etc.Name(), name)
	cmd := utils.Command("sysctl", "-p", confPath)
	return cmd.Run()
}
