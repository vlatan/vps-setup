package setup

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/vlatan/vps-setup/internal/utils"
)

// InstallDocker installs and configures Docker
func (s *Setup) InstallDocker() error {

	// Install using the APT repository
	// https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository
	fmt.Println("Setting up Docker...")

	// Add Docker's official GPG key
	cmds := [][]string{
		{"apt-get", "install", "-y", "ca-certificates", "curl"},
		{"install", "-m", "0755", "-d", "/etc/apt/keyrings"},
		{
			"curl",
			"-fsSL",
			"https://download.docker.com/linux/ubuntu/gpg",
			"-o",
			"/etc/apt/keyrings/docker.asc",
		},
		{"chmod", "a+r", "/etc/apt/keyrings/docker.asc"},
	}

	for _, cmdArgs := range cmds {
		cmd := utils.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	// Add the repository to APT sources
	heredoc := []string{
		"sudo tee /etc/apt/sources.list.d/docker.sources <<EOF",
		"Types: deb",
		"URIs: https://download.docker.com/linux/ubuntu",
		`Suites: $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}")`,
		"Components: stable",
		"Signed-By: /etc/apt/keyrings/docker.asc",
		"EOF",
	}

	// We need the shell for the heredoc to work
	cmd := utils.Command("/bin/bash", "-c", strings.Join(heredoc, "\n"))
	if err := cmd.Run(); err != nil {
		return err
	}

	// Refresh the package lists to include the packages from the new Docker repository
	if err := utils.Command("apt-get", "update").Run(); err != nil {
		return err
	}

	// Install docker packages
	cmdArgs := []string{
		"apt-get",
		"install",
		"-y",
		"docker-ce",
		"docker-ce-cli",
		"containerd.io",
		"docker-buildx-plugin",
		"docker-compose-plugin",
	}

	cmd = utils.Command(cmdArgs[0], cmdArgs[1:]...)
	if err := cmd.Run(); err != nil {
		return err
	}

	// Manage Docker as a non-root user
	// https://docs.docker.com/engine/install/linux-postinstall/

	// Check if the docker group exists first
	cmd = utils.Command("getent", "group", "docker")
	if err := cmd.Run(); err != nil {
		cmd = utils.Command("groupadd", "docker")
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	// Add user to docker group
	cmd = utils.Command("usermod", "-aG", "docker", s.Username)
	if err := cmd.Run(); err != nil {
		return err
	}

	// Configure logging
	// https://docs.docker.com/config/containers/logging/local/

	// Configure Docker log driver
	dockerLogDriver := map[string]any{
		"log-driver": "local",
		"log-opts": map[string]string{
			"max-size": "10m",
		},
	}

	// Marshal to JSON with indentation
	jsonData, err := json.MarshalIndent(dockerLogDriver, "", "    ")
	if err != nil {
		return err
	}

	// Make parent directories
	name := "docker/daemon.json"
	if err := s.Etc.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}

	// Write to the file
	data := append(jsonData, '\n')
	if err := s.Etc.WriteFile(name, data, 0644); err != nil {
		return err
	}

	// Configure rsyslog to capture container logs
	// Note: Only containers with syslog driver and tags starting with "docker-" will be captured
	// Example in compose.yml:
	//   logging:
	//     driver: syslog
	//     options:
	//       tag: docker-servicename
	// https://docs.docker.com/config/containers/logging/syslog/

	rsyslogConf := []string{
		"# Create a template for the target log file",
		"$template CUSTOM_LOGS,\"/var/log/containers/%programname%.log\"",
		"",
		"if $programname startswith  'docker-' then ?CUSTOM_LOGS",
		"& stop",
	}

	// Make parent directories
	name = "rsyslog.d/40-docker.conf"
	if err := s.Etc.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}

	// Write to the file
	data = []byte(strings.Join(rsyslogConf, "\n") + "\n")
	if err := s.Etc.WriteFile(name, data, 0644); err != nil {
		return err
	}

	// Configure logrotate
	logrotateConf := []string{
		"# consult 'man logrotate' for explanation",
		"/var/log/containers/*.log {",
		"    daily",
		"    rotate 20",
		"    missingok",
		"    notifempty",
		"    compress",
		"    delaycompress",
		"    postrotate",
		"        /usr/lib/rsyslog/rsyslog-rotate",
		"    endscript",
		"}",
	}

	// Make parent directories
	name = "logrotate.d/docker"
	if err := s.Etc.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}

	// Write to the file
	data = []byte(strings.Join(logrotateConf, "\n") + "\n")
	if err := s.Etc.WriteFile(name, data, 0644); err != nil {
		return err
	}

	// Restart services in this order to avoid race conditions
	cmds = [][]string{
		{"systemctl", "restart", "rsyslog"},
		{"systemctl", "restart", "logrotate"},
		{"systemctl", "restart", "docker"},
	}

	for _, cmdArgs := range cmds {
		cmd := utils.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
