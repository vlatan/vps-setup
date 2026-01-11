package settings

import (
	"bufio"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// AttachUbuntuPro attaches Ubuntu Pro subscription to the machine
func AttachUbuntuPro(scanner *bufio.Scanner) error {

	prompt := "Provide Ubuntu Pro token [enter to skip]: "
	prompt = colors.Yellow(prompt)
	token, err := utils.AskSensitiveQuestion(prompt)

	if err != nil {
		return err
	}

	if token == "" {
		return nil
	}

	cmds := [][]string{
		{"apt-get", "install", "-y", "ubuntu-advantage-tools"},
		{"pro", "attach", token},
	}

	for _, cmdArgs := range cmds {
		cmd := utils.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
