package setup

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// ChangeSwappiness changes the system swappiness,
// which means changing the free memory treshold (%) at which
// the swap is going to begin to be utilized.
func (s *Setup) ChangeSwappiness() error {

	// Function to check if the swappiness input is valid
	valid := func(s string) bool {
		n, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		return n >= 0 && n <= 100
	}

	// If swappiness is not valid ask the user
	if !valid(s.Swappiness) {
		prompt := "Provide system swappines value 0-100 [20]: "
		prompt = colors.Yellow(prompt)

		// Keep asking the question if swappiness is invalid
		for {
			s.Swappiness = utils.AskQuestion(prompt, s.Scanner)

			if s.Swappiness == "" {
				s.Swappiness = "20"
				break
			}

			if valid(s.Swappiness) {
				break
			}
		}
	}

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
