package setup

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// HardenSSH hardens the SSH config
func HardenSSH(target *string, username string, scanner *bufio.Scanner, etc *os.Root) error {

	var sshPort string
	prompt := "On which PORT do you want to connect via SSH [22]: "
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

	// Only the port is not commented
	hardenContent := []string{
		"Port " + sshPort,
		"# UsePAM no",
		"# AddressFamily inet",
		"# PermitRootLogin no",
		"# PermitEmptyPasswords no",
		"# PasswordAuthentication no",
		"# KbdInteractiveAuthentication no",
		"# AuthenticationMethods publickey",
		"# AllowUsers " + username,
	}

	// // Write to file
	name := "ssh/sshd_config.d/harden.conf"
	data := []byte(strings.Join(hardenContent, "\n") + "\n")
	if err := utils.WriteFile(etc, name, data); err != nil {
		return err
	}

	// Restart SSH, If SSH port is changed we need daemon-reload too
	cmds := [][]string{
		{"systemctl", "daemon-reload"},
		{"systemctl", "restart", "ssh"},
	}

	for _, cmdArgs := range cmds {
		cmd := utils.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	// Set SSH port to target
	*target = sshPort

	return nil
}
