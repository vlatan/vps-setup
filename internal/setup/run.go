package setup

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// Run runs all the jobs in the setup
func (s *Setup) Run() {

	msg := "WARNING: This script will modify the machine:"
	fmt.Println(colors.Red(msg))

	// Print all the jobs infos
	jobs := s.GetJobs()
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

	// Update the system
	cmd := utils.Command("apt-get", "update")
	if err := cmd.Run(); err != nil {
		utils.Exit(err)
	}

	// Execute the jobs
	if err := s.ProcessJobs(jobs); err != nil {
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
}
