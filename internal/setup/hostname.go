package setup

import (
	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// SetHostname sets a hostname on the machine
func (s *Setup) SetHostname() error {

	prompt := colors.Yellow("Provide hostname: ")
	for {
		hostname := utils.AskQuestion(prompt, s.Scanner)
		if hostname != "" {
			cmd := utils.Command("hostnamectl", "set-hostname", hostname)
			return cmd.Run()
		}
	}

}
