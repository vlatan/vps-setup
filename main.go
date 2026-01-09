package main

import (
	"bufio"
	"os"

	"github.com/vlatan/vps-setup/internal/settings"
)

// Establish root for securely opening files
// const root = "/"

func main() {

	root, err := os.OpenRoot("/")
	if err != nil {
		panic(err)
	}
	defer root.Close()

	if err = settings.SetNeedRestart(root); err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	settings.ChangeSwappiness(scanner)

}
