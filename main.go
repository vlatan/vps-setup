package main

import (
	"bufio"
	"os"

	"github.com/vlatan/vps-setup/internal/settings"
)

func main() {

	var username, sshPort string
	scanner := bufio.NewScanner(os.Stdin)

	etc, err := os.OpenRoot("/etc")
	if err != nil {
		panic(err)
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

		// "Setup ufw (uncomplicated firewall)",
		// "Install and configure Docker",
		// "Install and configure Postfix",
		// "Install and configure Fail2Ban",
		// "Format the bash prompt",
		// "Create bare git repository",
	}

	// Start the machine setup
	if err := settings.Start(scanner, jobs); err != nil {
		panic(err)
	}
}
