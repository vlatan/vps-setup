package settings

import (
	"bufio"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// AddUser adds new user and makes that user sudoer
func AddUser(target *string, scanner *bufio.Scanner) error {

	var username string
	prompt := colors.Yellow("Provide username: ")

	for {
		username = utils.AskQuestion(prompt, scanner)
		if username != "" {
			break
		}
	}

	cmds := [][]string{
		{"adduser", username},
		{"adduser", username, "sudo"},
	}

	for _, cmd := range cmds {
		if err := utils.RunCommand(cmd[0], cmd[1:]...); err != nil {
			return err
		}
	}

	// Set username to target
	*target = username

	return nil
}
