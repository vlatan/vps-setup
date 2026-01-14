package setup

import (
	"fmt"
	"strings"

	"github.com/vlatan/vps-setup/internal/utils"
)

// InstallFail2Ban installs and configures Fail2Ban
func (s *Setup) InstallFail2Ban() error {

	fmt.Println("Setting up Fail2Ban...")
	cmd := utils.Command("apt-get", "install", "-y", "fail2ban")
	if err := cmd.Run(); err != nil {
		return err
	}

	// Config content
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
		fmt.Sprintf("port = %s", s.SSHPort),
	}

	// Write to the file
	name := "fail2ban/jail.local"
	data := []byte(strings.Join(content, "\n") + "\n")
	if err := utils.WriteFile(s.Etc, name, data); err != nil {
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
