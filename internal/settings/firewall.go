package settings

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// SetupFirewall sets ap an uncomplicated firewall (ufw) on the machine
func SetupFirewall(sshPort string, scanner *bufio.Scanner, etc *os.Root) error {

	msg := colors.Yellow("Setting up firewall (ufw)...")
	fmt.Println(msg)

	msg = "Note: Exposed ports using Docker will bypass the ufw rules"
	fmt.Println(msg)

	url := "https://docs.docker.com/engine/install/ubuntu/#firewall-limitations"
	fmt.Println("Source: " + url)

	// Install ufw
	cmd := utils.Command("apt-get", "install", "-y", "ufw")
	if err := cmd.Run(); err != nil {
		return err
	}

	name := "default/ufw"
	data, err := etc.ReadFile(name)
	if err != nil {
		return err
	}

	// Enforce IPV6 firewall rules
	data = bytes.ReplaceAll(data, []byte("IPV6=no"), []byte("IPV6=yes"))
	if err := etc.WriteFile(name, data, 0644); err != nil {
		return err
	}

	cmds := [][]string{
		{"ufw", "default", "allow", "outgoing"},
		{"ufw", "default", "deny", "incoming"},
		{"ufw", "allow", fmt.Sprintf("%s/tcp", sshPort)},
		{"ufw", "allow", "http/tcp"},
		{"ufw", "allow", "https/tcp"},
		{"ufw", "--force", "enable"},
		{"ufw", "logging", "on"},
	}

	for _, cmdArgs := range cmds {
		cmd := utils.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
