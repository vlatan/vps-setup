package setup

import (
	"fmt"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// SetTimezone sets a timezone on the machine
func (s *Setup) SetTimezone() error {

	// Check for env var timezone
	if s.Timezone == "" {
		prompt := colors.Yellow("Provide timezone [Continent/City] [UTC]: ")
		s.Timezone = utils.AskQuestion(prompt, s.Scanner)
		if s.Timezone == "" {
			s.Timezone = "UTC"
		}
	}

	fmt.Println("Setting up timezone...")
	cmd := utils.Command("timedatectl", "set-timezone", s.Timezone)
	return cmd.Run()
}
