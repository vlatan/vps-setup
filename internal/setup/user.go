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

	prompt := colors.Yellow("Provide username: ")
	for {
		// Check for env var username first
		if s.Username != "" {
			break
		}
		s.Username = utils.AskQuestion(prompt, s.Scanner)
	}

	cmds := [][]string{
		{"adduser", "--gecos", "", s.Username},
		{"adduser", s.Username, "sudo"},
	}

	for _, cmdArgs := range cmds {
		cmd := utils.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	userDir := filepath.Join("/home", s.Username)
	home, err := os.OpenRoot(userDir)
	if err != nil {
		return err
	}

	// Provide the user home open root to the setup
	s.Home = home

	return nil
}
