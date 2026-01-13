package setup

import (
	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/config"
	"github.com/vlatan/vps-setup/internal/utils"
)

// SetTimezone sets a timezone on the machine
func SetTimezone(cfg *config.Config) error {
	prompt := colors.Yellow("Provide timezone [Continent/City] [UTC]: ")
	timezone := utils.AskQuestion(prompt, cfg.Scanner)
	if timezone == "" {
		timezone = "UTC"
	}

	cmd := utils.Command("timedatectl", "set-timezone", timezone)
	return cmd.Run()
}
