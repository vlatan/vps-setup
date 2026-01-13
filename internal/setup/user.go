package setup

import (
	"os"
	"path/filepath"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// AddUser adds new user and makes that user sudoer.
// This method will modify the setup state,
// namely s.Username and s.Home.
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

	// Provide the username to the setup
	s.Username = username

	userDir := filepath.Join("/home", s.Username)
	home, err := os.OpenRoot(userDir)
	if err != nil {
		return err
	}

	// Provide the user home open root to the setup
	s.Home = home

	return nil
}
