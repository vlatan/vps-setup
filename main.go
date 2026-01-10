package main

import (
	"bufio"
	"os"

	"github.com/vlatan/vps-setup/internal/settings"
)

func main() {

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

		// "Set timezone",
		// "Add new user",
		// "Harden SSH access",
		// "Setup ufw (uncomplicated firewall)",
		// "Install and configure Docker",
		// "Install and configure Postfix",
		// "Install and configure Fail2Ban",
		// "Format the bash prompt",
		// "Create bare git repository",
	}

	// Check whether to start
	if err := settings.Start(scanner, jobs); err != nil {
		panic(err)
	}
}
