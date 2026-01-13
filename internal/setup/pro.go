package setup

import (
	"fmt"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// AttachUbuntuPro attaches Ubuntu Pro subscription to the machine
func (s *Setup) AttachUbuntuPro() error {

	if err := s.setUbuntuToken(); err != nil {
		return err
	}

	fmt.Println(colors.Yellow("Attaching Ubuntu Pro..."))

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

// setToken sets
func (s *Setup) setUbuntuToken() error {

	if s.UbuntuProToken != "" {
		return nil
	}

	prompt := "Provide Ubuntu Pro token [optional]: "
	prompt = colors.Yellow(prompt)
	token, err := utils.AskPassword(prompt)

	if err != nil {
		return err
	}

	if token == "" {
		return nil
	}

	s.UbuntuProToken = token
	return nil
}
