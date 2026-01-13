package setup

import (
	"fmt"

	"github.com/vlatan/vps-setup/internal/colors"
)

type Job struct {
	Info     string
	Callable func() error
}

// GetJobs returns a slice of all the jobs to be done.
// The order of some methods is important
func (s *Setup) GetJobs() []Job {
	return []Job{
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
}

// ProcessJobs processes the jobs in a given slice
func (s *Setup) ProcessJobs(jobs []Job) error {

	for _, job := range jobs {

		if job.Callable == nil {
			if job.Info != "" {
				fmt.Printf("Skipping: %s\n", job.Info)
			}
			continue
		}

		if err := job.Callable(); err != nil {
			msg := colors.Red(fmt.Sprintf("Failed: %s", job.Info))
			return fmt.Errorf("%s\n%w", msg, err)
		}
	}

	return nil
}
