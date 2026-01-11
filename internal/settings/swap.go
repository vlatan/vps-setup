package settings

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// ChangeSwappiness changes the system swappiness,
// which means changing the free memory treshold (%) at which
// the swap is going to begin to be utilized.
func ChangeSwappiness(scanner *bufio.Scanner, etc *os.Root) error {

	var swappiness string
	prompt := "Provide system swappines value 0-100 [enter for 20]: "
	prompt = colors.Yellow(prompt)

	// Function to check if the swappiness input is valid
	valid := func(s string) bool {
		n, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		return n >= 0 && n <= 100
	}

	// Keep asking the question if swappiness is invalid
	for {
		swappiness = utils.AskQuestion(prompt, scanner)

		if swappiness == "" {
			swappiness = "20"
			break
		}

		if valid(swappiness) {
			break
		}
	}

	// Write to the file
	name := "sysctl.d/99-my-swappiness.conf"
	data := fmt.Appendf([]byte{}, "vm.swappiness = %s\n", swappiness)
	if err := utils.WriteFile(etc, name, data); err != nil {
		return err
	}

	// Load our config
	confPath := filepath.Join(etc.Name(), name)
	cmd := utils.Command("sysctl", "-p", confPath)
	return cmd.Run()
}
