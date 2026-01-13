package settings

import (
	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/config"
	"github.com/vlatan/vps-setup/internal/utils"
)

// SetHostname sets a hostname on the machine
func SetHostname(cfg *config.Config) error {

	prompt := colors.Yellow("Provide hostname: ")
	for {
		hostname := utils.AskQuestion(prompt, cfg.Scanner)
		if hostname != "" {
			cmd := utils.Command("hostnamectl", "set-hostname", hostname)
			return cmd.Run()
		}
	}

}
