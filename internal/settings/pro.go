package settings

import (
	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/config"
	"github.com/vlatan/vps-setup/internal/utils"
)

// AttachUbuntuPro attaches Ubuntu Pro subscription to the machine
func AttachUbuntuPro(cfg *config.Config) error {

	prompt := "Provide Ubuntu Pro token [optional]: "
	prompt = colors.Yellow(prompt)
	token, err := utils.AskPassword(prompt)

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
