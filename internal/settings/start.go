package settings

import (
	"bufio"
	"fmt"
	"slices"
	"strings"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

type Job struct {
	Info     string
	Callable func() error
}

// Start enumerates and prints the jobs and asks
// the user whether to continue.
func ProcessJobs(scanner *bufio.Scanner, jobs []Job) error {

	// Print all jobs
	for _, job := range jobs {
		msg := colors.Yellow(fmt.Sprintf("* %s", job.Info))
		fmt.Println(msg)
	}

	prompt := "Continue? [y/n]: "
	start := strings.ToLower(utils.AskQuestion(prompt, scanner))
	if !slices.Contains([]string{"yes", "y"}, start) {
		return nil
	}

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
