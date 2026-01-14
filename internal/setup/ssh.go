package setup

import (
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// HardenSSH hardens the SSH config
// This method will modify the setup state,
// namely assign s.SSHPort and optionally add an SSH public key
// for the user.
func (s *Setup) HardenSSH() error {

	var sshPort string
	prompt := "On which PORT do you want to connect via SSH [22]: "
	prompt = colors.Yellow(prompt)

	// Function to check if the port input is valid
	valid := func(s string) bool {
		n, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		return n >= 0 && n <= 65535
	}

	// Keep asking the question if port is invalid
	for {
		sshPort = utils.AskQuestion(prompt, s.Scanner)

		if sshPort == "" {
			sshPort = "22"
			break
		}

		if valid(sshPort) {
			break
		}
	}

	// Only the port is not commented
	hardenContent := []string{
		"Port " + sshPort,
		"# UsePAM no",
		"# AddressFamily inet",
		"# PermitRootLogin no",
		"# PermitEmptyPasswords no",
		"# PasswordAuthentication no",
		"# KbdInteractiveAuthentication no",
		"# AuthenticationMethods publickey",
		"# AllowUsers " + s.Username,
	}

	// // Write to file
	name := "ssh/sshd_config.d/harden.conf"
	data := []byte(strings.Join(hardenContent, "\n") + "\n")
	if err := utils.WriteFile(s.Etc, name, data); err != nil {
		return err
	}

	// Restart SSH, If SSH port is changed we need daemon-reload too
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

	// Set SSH port to setup struct
	s.SSHPort = sshPort

	return nil
}

// addSSHKey writes an SSH public key to user's authorized_keys file
func (s *Setup) addSSHKey(pubKey string) error {

	// Get user's UID and GID
	u, err := user.Lookup(s.Username)
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return err
	}

	sshDir := ".ssh"
	authKeysFile := filepath.Join(sshDir, "authorized_keys")

	// Create .ssh directory
	if err := s.Home.MkdirAll(sshDir, 0700); err != nil {
		return err
	}

	// Write the public key to authorized_keys
	data := []byte(strings.TrimSpace(pubKey) + "\n")
	if err := s.Home.WriteFile(authKeysFile, data, 0600); err != nil {
		return err
	}

	// Change ownership of the .ssh directory
	if err := s.Home.Chown(sshDir, uid, gid); err != nil {
		return err
	}

	// Change ownership of the authorized_keys
	return s.Home.Chown(authKeysFile, uid, gid)
}
