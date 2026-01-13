package main

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/setup"
	"github.com/vlatan/vps-setup/internal/utils"
)

func main() {

	s := setup.New()
	defer s.Close()

	jobs := []utils.Job{
		{
			Info:     "Enable services autorestart",
			Callable: func() error { return s.AutoRestart() },
		},
		{
			Info:     "Change the swappiness",
			Callable: func() error { return s.ChangeSwappiness() },
		},
		{
			Info:     "Attach to Ubuntu Pro",
			Callable: func() error { return s.AttachUbuntuPro() },
		},
		{
			Info:     "Set hostname",
			Callable: func() error { return s.SetHostname() },
		},
		{
			Info:     "Set timezone",
			Callable: func() error { return s.SetTimezone() },
		},
		{
			Info:     "Add new user",
			Callable: func() error { return s.AddUser() },
		},
		{
			Info:     "Harden SSH access (add commented rules)",
			Callable: func() error { return s.HardenSSH() },
		},
		{
			Info:     "Setup Postfix SMTP relay",
			Callable: func() error { return s.SetupPostfix() },
		},
		{
			Info:     "Setup ufw (uncomplicated firewall)",
			Callable: func() error { return s.SetupFirewall() },
		},
		{
			Info:     "Install and configure Fail2Ban",
			Callable: func() error { return s.InstallFail2Ban() },
		},
		{
			Info:     "Install and configure Docker",
			Callable: func() error { return s.InstallDocker() },
		},
		{
			Info:     "Format the bash prompt",
			Callable: func() error { return s.FormatBash() },
		},
		{
			Info:     "Create bare git repository",
			Callable: func() error { return s.SetupGitRepo() },
		},
	}

	msg := "WARNING: This script will modify the machine:"
	fmt.Println(colors.Red(msg))

	// Print all the jobs infos
	for _, job := range jobs {
		msg := colors.Yellow(fmt.Sprintf("* %s", job.Info))
		fmt.Println(msg)
	}

	// Check if the user wants to continue
	prompt := "Continue? [y/N]: "
	startScript := strings.ToLower(utils.AskQuestion(prompt, s.Scanner))
	if !slices.Contains([]string{"yes", "y"}, startScript) {
		return
	}

	startTime := time.Now()

	// Execute the jobs
	if err := utils.ProcessJobs(s.Scanner, jobs); err != nil {
		utils.Exit(err)
	}

	timeTook := time.Since(startTime)

	fmt.Println(
		colors.Green("Installation done. Time took:"),
		fmt.Sprintf("%.2f", timeTook.Seconds()),
		colors.Green("seconds."),
	)

	fmt.Println(
		colors.Green("Log out and log back in on port"),
		colors.Yellow(s.SSHPort),
		colors.Green("with user"),
		colors.Yellow(s.Username)+colors.Green("."),
	)
	msg = "Make sure you complete the setup according to the documentaion."
	fmt.Println(colors.Green(msg))
}
