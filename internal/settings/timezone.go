package settings

import (
	"bufio"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// SetTimezone sets a timezone on the machine
func SetTimezone(scanner *bufio.Scanner) error {
	prompt := colors.Yellow("Provide timezone [Continent/City]: ")
	timezone := utils.AskQuestion(prompt, scanner)
	cmd := utils.Command("timedatectl", "set-timezone", timezone)
	return cmd.Run()
}
