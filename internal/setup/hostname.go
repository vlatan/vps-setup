package setup

import (
	"fmt"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// SetHostname sets a hostname on the machine
func (s *Setup) SetHostname() error {

	prompt := colors.Yellow("Provide hostname: ")
	for {
		// Check for env var hostname first
		if s.Hostname != "" {
			break
		}
		s.Hostname = utils.AskQuestion(prompt, s.Scanner)
	}

	fmt.Println(colors.Yellow("Seting up hostname..."))
	cmd := utils.Command("hostnamectl", "set-hostname", s.Hostname)
	return cmd.Run()
}
