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
	token := utils.AskQuestion(prompt, scanner)
	if token == "" {
		return nil
	}

	cmds := [][]string{
		{"apt", "install", "-y", "ubuntu-advantage-tools"},
		{"pro", "attach", token},
	}

	for _, cmd := range cmds {
		if err := utils.RunCommand(cmd[0], cmd[1:]...); err != nil {
			return err
		}
	}

	return nil
}
