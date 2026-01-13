package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/setup"
	"github.com/vlatan/vps-setup/internal/utils"
)

func main() {

	msg := "WARNING: This script will modify the machine:"
	fmt.Println(colors.Red(msg))

	s := setup.New()
	defer s.Close()

	var username, sshPort string
	var home *os.Root
	scanner := bufio.NewScanner(os.Stdin)

	// Open /etc as root
	etc, err := os.OpenRoot("/etc")
	if err != nil {
		utils.Exit(err)
	}
	defer etc.Close()

	primaryJobs := []utils.Job{
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
			Callable: func() error { return setup.HardenSSH(&sshPort, username, scanner, etc) },
		},
		{
			Info:     "Setup Postfix SMTP relay",
			Callable: func() error { return setup.SetupPostfix(scanner, etc) },
		},
		{
			Info:     "Setup ufw (uncomplicated firewall)",
			Callable: func() error { return setup.SetupFirewall(sshPort, scanner, etc) },
		},
		{
			Info:     "Install and configure Fail2Ban",
			Callable: func() error { return setup.InstallFail2Ban(sshPort, scanner, etc) },
		},
		{
			Info:     "Install and configure Docker",
			Callable: func() error { return setup.InstallDocker(username, scanner, etc) },
		},
	}

	// These jobs require the home user root
	secondaryJobs := []utils.Job{
		{
			Info:     "Format the bash prompt",
			Callable: func() error { return setup.FormatBash(scanner, home) },
		},
		{
			Info:     "Create bare git repository",
			Callable: func() error { return setup.SetupGitRepo(username, scanner, home) },
		},
	}

	// Print all the jobs infos
	allJobs := append(primaryJobs, secondaryJobs...)
	for _, job := range allJobs {
		msg := colors.Yellow(fmt.Sprintf("* %s", job.Info))
		fmt.Println(msg)
	}

	// Check if the user wants to continue
	prompt := "Continue? [y/N]: "
	startScript := strings.ToLower(utils.AskQuestion(prompt, scanner))
	if !slices.Contains([]string{"yes", "y"}, startScript) {
		return
	}

	startTime := time.Now()

	// Execute the first batch of jobs
	if err := utils.ProcessJobs(scanner, primaryJobs); err != nil {
		utils.Exit(err)
	}

	// Open /home/xxx as root
	userDir := filepath.Join("/home", username)
	home, err = os.OpenRoot(userDir)
	if err != nil {
		utils.Exit(err)
	}
	defer home.Close()

	// Execute the second batch of jobs
	if err := utils.ProcessJobs(scanner, secondaryJobs); err != nil {
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
		colors.Yellow(sshPort),
		colors.Green("with user"),
		colors.Yellow(username)+colors.Green("."),
	)
	msg = "Make sure you complete the setup according to the documentaion."
	fmt.Println(colors.Green(msg))
}
