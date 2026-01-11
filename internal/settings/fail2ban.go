package settings

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// InstallFail2Ban installs and configures Fail2Ban
func InstallFail2Ban(sshPort string, scanner *bufio.Scanner, etc *os.Root) error {

	prompt := "Do you want to install Fail2Ban? [y/n]: "
	prompt = colors.Yellow(prompt)
	start := strings.ToLower(utils.AskQuestion(prompt, scanner))
	if !slices.Contains([]string{"yes", "y"}, start) {
		return nil
	}

	cmd := utils.Command("apt-get", "install", "-y", "fail2ban")
	if err := cmd.Run(); err != nil {
		return err
	}

	// Create dirs that do not exist in the file path
	name := "fail2ban/jail.local"
	if err := etc.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}

	content := []string{
		"[DEFAULT]",
		"bantime  = 30m",
		"banaction = ufw",
		"action = %(action_)s",
		"",
		"# action = %(action_mwl)s",
		"# sender = root@domain.com",
		"# destemail = webmaster@domain.com",
		"",
		"[sshd]",
		"enabled = true",
		fmt.Sprintf("port = %s", sshPort),
	}

	// Write to the file
	data := strings.Join(content, "\n") + "\n"
	if err := etc.WriteFile(name, []byte(data), 0644); err != nil {
		return err
	}

	cmds := [][]string{
		{"systemctl", "enable", "fail2ban"},
		{"systemctl", "start", "fail2ban"},
	}

	for _, cmdArgs := range cmds {
		cmd := utils.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
