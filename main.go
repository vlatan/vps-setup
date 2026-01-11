package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/settings"
	"github.com/vlatan/vps-setup/internal/utils"
)

func main() {

	msg := "WARNING: This script will modify the machine"
	fmt.Println(colors.Red(msg))

	var username, sshPort string
	scanner := bufio.NewScanner(os.Stdin)

	// Open /etc as root
	etc, err := os.OpenRoot("/etc")
	if err != nil {
		utils.Exit(err)
	}
	defer etc.Close()

	jobs := []settings.Job{
		{
			Info:     "Enable services autorestart",
			Callable: func() error { return settings.AutoRestart(etc) },
		},
		{
			Info:     "Change the swappiness",
			Callable: func() error { return settings.ChangeSwappiness(scanner, etc) },
		},
		{
			Info:     "Attach to Ubuntu Pro",
			Callable: func() error { return settings.AttachUbuntuPro(scanner) },
		},
		{
			Info:     "Set hostname",
			Callable: func() error { return settings.SetHostname(scanner) },
		},
		{
			Info:     "Set timezone",
			Callable: func() error { return settings.SetTimezone(scanner) },
		},
		{
			Info:     "Add new user",
			Callable: func() error { return settings.AddUser(&username, scanner) },
		},
		{
			Info:     "Harden SSH access",
			Callable: func() error { return settings.HardenSSH(&sshPort, scanner, etc) },
		},
		{
			Info:     "Install and configure Postfix",
			Callable: func() error { return settings.InstallPostfix(scanner, etc) },
		},
	}

	// Execute jobs
	if err := settings.ProcessJobs(scanner, jobs); err != nil {
		utils.Exit(err)
	}

	// Open /home/xxx as root
	userDir := filepath.Join("/home", username)
	home, err := os.OpenRoot(userDir)
	if err != nil {
		utils.Exit(err)
	}
	defer home.Close()

	// These jobs require sshPort, username and home
	jobs = []settings.Job{
		{
			Info:     "Setup ufw (uncomplicated firewall)",
			Callable: func() error { return settings.SetupFirewall(sshPort, scanner, etc) },
		},
		{
			Info:     "Install and configure Fail2Ban",
			Callable: func() error { return settings.InstallFail2Ban(sshPort, scanner, etc) },
		},
		{
			Info:     "Install and configure Docker",
			Callable: func() error { return settings.InstallDocker(username, scanner, etc) },
		},
		{
			Info:     "Format the bash prompt",
			Callable: func() error { return settings.FormatBash(scanner, home) },
		},
		{
			Info: "Create bare git repository",
		},
	}

	// Execute jobs
	if err := settings.ProcessJobs(scanner, jobs); err != nil {
		utils.Exit(err)
	}
}
