package main

import (
	"bufio"
	"os"

	"github.com/vlatan/vps-setup/internal/settings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	settings.ChangeSwappiness(scanner)
}
