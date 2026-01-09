package settings

import (
	"bufio"
	"fmt"
	"strconv"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// ChangeSwappiness changes the system swappiness,
// which means changing the free memory treshold (%) at which
// the swap is going to begin to be utilized.
func ChangeSwappiness(scanner *bufio.Scanner) {

	prompt := "Provide system swappines value 0-100 [enter for 20]: "
	prompt = colors.Yellow(prompt)
	swappiness := utils.AskQuestion(prompt, scanner)

	// Function to check if the swappiness input is valid
	valid := func(s string) bool {
		n, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		return n >= 0 && n <= 100
	}

	for {
		if swappiness == "" {
			swappiness = "20"
			break
		}

		if valid(swappiness) {
			break
		}

		swappiness = utils.AskQuestion(prompt, scanner)
	}

	fmt.Println(swappiness)
}
