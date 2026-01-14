package setup

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// HardenSSH hardens the SSH config
// This method will modify the setup state,
// namely assign s.SSHPort and optionally s.SSHPubKey
func (s *Setup) HardenSSH() error {

	s.setSSHPort()
	s.setSSHPubKey()
	fmt.Println("Hardening SSH access...")

	hardenContent := []string{
		"Port " + s.SSHPort,
		"AddressFamily inet",
		"PermitRootLogin no",
		"PermitEmptyPasswords no",
		"AllowUsers " + s.Username,
	}

	pound := "# "
	if s.SSHPubKey != "" {
		pound = ""
		if err := s.addSSHPubKey(); err != nil {
			return err
		}
	}

	rest := []string{
		pound + "UsePAM no",
		pound + "PasswordAuthentication no",
		pound + "KbdInteractiveAuthentication no",
		pound + "AuthenticationMethods publickey",
	}

	hardenContent = append(hardenContent, rest...)

	// Write to file
	name := "ssh/sshd_config.d/harden.conf"
	data := []byte(strings.Join(hardenContent, "\n") + "\n")
	if err := utils.WriteFile(s.Etc, name, data); err != nil {
		return err
	}

	// Restart SSH
	// If SSH port is changed we need daemon-reload too
	cmds := [][]string{
		{"systemctl", "daemon-reload"},
		{"systemctl", "restart", "ssh"},
	}

	for _, cmdArgs := range cmds {
		cmd := utils.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Setup) setSSHPort() {

	// Helper function to check if the port input is valid
	valid := func(s string) bool {
		n, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		return n >= 0 && n <= 65535
	}

	// Check if a valid SSH port value is already set
	if valid(s.SSHPort) {
		return
	}

	prompt := colors.Yellow("Provide SSH port [22]: ")
	for { // Keep asking the question if SSH port is invalid
		s.SSHPort = utils.AskQuestion(prompt, s.Scanner)

		// If empty response provide default value
		if s.SSHPort == "" {
			s.SSHPort = "22"
			break
		}

		if valid(s.SSHPort) {
			break
		}
	}
}

// setSSHPubKey sets the SSH public key value to the setup struct
func (s *Setup) setSSHPubKey() {

	// Helper function to check if the SSH public key input is valid
	valid := func(s string) bool {
		parts := strings.Split(s, " ")
		return len(parts) >= 2
	}

	// Check if a valid SSH public key value is already set
	if valid(s.SSHPubKey) {
		return
	}

	prompt := colors.Yellow("Provide SSH public key [optional]: ")
	for { // Keep asking the question if SSH public key is invalid
		s.SSHPubKey = utils.AskQuestion(prompt, s.Scanner)

		// If empty response skip
		if s.SSHPubKey == "" {
			break
		}

		if valid(s.SSHPubKey) {
			break
		}
	}

}

// addSSHPubKey writes an SSH public key to user's authorized_keys file
func (s *Setup) addSSHPubKey() error {

	sshDir := ".ssh"
	authKeysFile := filepath.Join(sshDir, "authorized_keys")

	// Create .ssh directory
	if err := s.Home.MkdirAll(sshDir, 0700); err != nil {
		return err
	}

	// Write the public key to authorized_keys
	data := []byte(strings.TrimSpace(s.SSHPubKey) + "\n")
	if err := s.Home.WriteFile(authKeysFile, data, 0600); err != nil {
		return err
	}

	// Change ownership of the .ssh directory
	if err := s.Home.Chown(sshDir, s.Uid, s.Gid); err != nil {
		return err
	}

	// Change ownership of the authorized_keys
	return s.Home.Chown(authKeysFile, s.Uid, s.Gid)
}
