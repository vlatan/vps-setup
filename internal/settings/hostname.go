package settings

import (
	"bufio"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// SetHostname sets a hostname on the machine
func SetHostname(scanner *bufio.Scanner) error {

	prompt := colors.Yellow("Provide hostname: ")
	for {
		hostname := utils.AskQuestion(prompt, scanner)
		if hostname != "" {
			return utils.RunCommand("hostnamectl", "set-hostname", hostname)
		}
	}

}
