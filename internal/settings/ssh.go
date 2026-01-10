package settings

import (
	"bufio"
	"os"
	"strconv"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// HardenSSH hardens the SSH config
func HardenSSH(scanner *bufio.Scanner, etc *os.Root) error {

	var port string
	prompt := "On which PORT do you want to connect via SSH [enter for 22]: "
	prompt = colors.Yellow(prompt)

	// Function to check if the port input is valid
	valid := func(s string) bool {
		n, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		return n >= 0 && n <= 65535
	}

	// Keep asking the question if port is invalid
	for {
		port = utils.AskQuestion(prompt, scanner)

		if port == "" {
			port = "22"
			break
		}

		if valid(port) {
			break
		}
	}

	return nil
}
