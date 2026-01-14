package setup

import (
	"fmt"

	"github.com/vlatan/vps-setup/internal/utils"
)

// AttachUbuntuPro attaches Ubuntu Pro subscription to the machine
func (s *Setup) AttachUbuntuPro() error {

	if s.UbuntuProToken == "" {
		fmt.Println("Skipping attaching Ubuntu Pro...")
		return nil
	}

	fmt.Println("Attaching Ubuntu Pro...")
	cmds := [][]string{
		{"apt-get", "install", "-y", "ubuntu-advantage-tools"},
		{"pro", "attach", s.UbuntuProToken},
	}

	for _, cmdArgs := range cmds {
		cmd := utils.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
