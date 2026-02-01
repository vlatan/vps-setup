package setup

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vlatan/vps-setup/internal/utils"
)

// HardenSSH hardens the SSH config
// This method will modify the setup state,
// namely assign s.SSHPort and optionally s.SSHPubKey
func (s *Setup) HardenSSH() error {

	if !validPort(s.SSHPort) || !validKey(s.SSHPubKey) {
		return errors.New("invalid SSH port and/or public key")
	}

	fmt.Println("Hardening SSH access...")
	if err := s.addSSHPubKey(); err != nil {
		return err
	}

	hardenContent := []string{
		"Port " + s.SSHPort,
		"AddressFamily inet",
		"PermitRootLogin no",
		"PermitEmptyPasswords no",
		"AllowUsers " + s.Username,
		"UsePAM no",
		"PasswordAuthentication no",
		"KbdInteractiveAuthentication no",
		"AuthenticationMethods publickey",
	}

	// Make parent directories
	name := "ssh/sshd_config.d/harden.conf"
	if err := s.Etc.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}

	// Write to file
	data := []byte(strings.Join(hardenContent, "\n") + "\n")
	if err := s.Etc.WriteFile(name, data, 0644); err != nil {
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

	// Change the ownership
	for _, path := range []string{authKeysFile, sshDir} {
		if err := s.Home.Chown(path, s.Uid, s.Gid); err != nil {
			return err
		}
	}

	return nil
}

func validPort(port string) bool {
	n, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	return n >= 0 && n <= 65535
}

func validKey(key string) bool {
	parts := strings.Split(key, " ")
	return len(parts) >= 2
}
