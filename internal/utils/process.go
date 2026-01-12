package utils

import (
	"bufio"
	"fmt"

	"github.com/vlatan/vps-setup/internal/colors"
)

type Job struct {
	Info     string
	Callable func() error
}

// Start enumerates and prints the jobs and asks
// the user whether to continue.
func ProcessJobs(scanner *bufio.Scanner, jobs []Job) error {

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
