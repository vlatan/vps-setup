package setup

import (
	"errors"
	"fmt"

	"github.com/vlatan/vps-setup/internal/utils"
)

// SetHostname sets a hostname on the machine
func (s *Setup) SetHostname() error {

	if s.Hostname == "" {
		return errors.New("no hostname found")
	}

	fmt.Println("Setting up hostname...")
	cmd := utils.Command("hostnamectl", "set-hostname", s.Hostname)
	return cmd.Run()
}
