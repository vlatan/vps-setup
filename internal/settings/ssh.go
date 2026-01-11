package settings

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// HardenSSH hardens the SSH config
func HardenSSH(target *string, scanner *bufio.Scanner, etc *os.Root) error {

	var sshPort string
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
		sshPort = utils.AskQuestion(prompt, scanner)

		if sshPort == "" {
			sshPort = "22"
			break
		}

		if valid(sshPort) {
			break
		}
	}

	// // Write to file
	name := "ssh/sshd_config.d/harden.conf"
	data := fmt.Sprintf("Port %s\nAddressFamily inet\nPermitRootLogin no\n", sshPort)
	if err := utils.WriteFile(etc, name, []byte(data)); err != nil {
		return err
	}

	// Restart SSH
	cmd := utils.Command("systemctl", "restart", "ssh")
	if err := cmd.Run(); err != nil {
		return err
	}

	// Set SSH port to target
	*target = sshPort

	return nil
}
