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
// The order of the methods is important,
// because some of them change the state
// and others need that new state.
func (s *Setup) GetJobs() []Job {
	return []Job{
		{
			Info:     "Enable services autorestart",
			Callable: s.AutoRestart,
		},
		{
			Info:     "Change the swappiness",
			Callable: s.ChangeSwappiness,
		},
		{
			Info:     "Attach to Ubuntu Pro",
			Callable: s.AttachUbuntuPro,
		},
		{
			Info:     "Set hostname",
			Callable: s.SetHostname,
		},
		{
			Info:     "Set timezone",
			Callable: s.SetTimezone,
		},
		{
			Info:     "Add new user",
			Callable: s.AddUser,
		},
		{
			Info:     "Harden SSH access",
			Callable: s.HardenSSH,
		},
		{
			Info:     "Setup Postfix SMTP relay",
			Callable: s.SetupPostfix,
		},
		{
			Info:     "Setup ufw (uncomplicated firewall)",
			Callable: s.SetupFirewall,
		},
		{
			Info:     "Install and configure Fail2Ban",
			Callable: s.InstallFail2Ban,
		},
		{
			Info:     "Install and configure Docker",
			Callable: s.InstallDocker,
		},
		{
			Info:     "Format the bash prompt",
			Callable: s.FormatBash,
		},
		{
			Info:     "Create bare git repository",
			Callable: s.SetupGitRepo,
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
