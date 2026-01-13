package setup

import (
	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// AddUser adds new user and makes that user sudoer
func (s *Setup) AddUser() error {

	var username string
	prompt := colors.Yellow("Provide username: ")

	for {
		username = utils.AskQuestion(prompt, s.Scanner)
		if username != "" {
			break
		}
	}

	cmds := [][]string{
		{"adduser", "--gecos", "", username},
		{"adduser", username, "sudo"},
	}

	for _, cmdArgs := range cmds {
		cmd := utils.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	// Provide the username to the setup struct
	s.Username = username

	return nil
}
