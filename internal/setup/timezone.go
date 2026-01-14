package setup

import (
	"fmt"

	"github.com/vlatan/vps-setup/internal/utils"
)

// SetTimezone sets a timezone on the machine
func (s *Setup) SetTimezone() error {
	fmt.Println("Setting up timezone...")
	cmd := utils.Command("timedatectl", "set-timezone", s.Timezone)
	return cmd.Run()
}
