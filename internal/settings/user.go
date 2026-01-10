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
		{"adduser", "--gecos", "", username},
		{"adduser", username, "sudo"},
	}

	for _, cmdArgs := range cmds {
		cmd := utils.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	// Set username to target
	*target = username

	return nil
}
