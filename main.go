package main

import (
	"bufio"
	"os"

	"github.com/vlatan/vps-setup/internal/settings"
)

// Establish root for securely opening files
// const root = "/"

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	settings.ChangeSwappiness(scanner)
}
