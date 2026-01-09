package settings

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// ChangeSwappiness changes the system swappiness,
// which means changing the free memory treshold (%) at which
// the swap is going to begin to be utilized.
func ChangeSwappiness(scanner *bufio.Scanner, etc *os.Root) error {

	prompt := "Provide system swappines value 0-100 [enter for 20]: "
	prompt = colors.Yellow(prompt)
	swappiness := utils.AskQuestion(prompt, scanner)

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
		if swappiness == "" {
			swappiness = "20"
			break
		}

		if valid(swappiness) {
			break
		}

		swappiness = utils.AskQuestion(prompt, scanner)
	}

	msg := fmt.Sprintf("Setting up system swappiness to %s...", swappiness)
	fmt.Println(colors.Yellow(msg))

	// Open file in append mode, create if it doesn't exist
	flag := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	file, err := etc.OpenFile("sysctl.conf", flag, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	// Append line to the file
	line := fmt.Sprintf("\nvm.swappiness = %s\n", swappiness)
	if _, err := file.WriteString(line); err != nil {
		return err
	}

	cmds := [][]string{
		{"sysctl", "-p"},
		{"swapoff", "-a"},
		{"swapon", "-a"},
	}

	for _, cmd := range cmds {
		if err = utils.RunCommand(cmd[0], cmd[1]); err != nil {
			return err
		}
	}

	return nil
}
