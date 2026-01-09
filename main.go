package main

import (
	"bufio"
	"os"

	"github.com/vlatan/vps-setup/internal/settings"
)

// Establish root for securely opening files
// const root = "/"

func main() {

	etc, err := os.OpenRoot("/etc")
	if err != nil {
		panic(err)
	}
	defer etc.Close()

	scanner := bufio.NewScanner(os.Stdin)

	callables := []func() error{
		func() error { return settings.SetNeedRestart(etc) },
		func() error { return settings.ChangeSwappiness(scanner, etc) },
	}

	for _, callable := range callables {
		if err = callable(); err != nil {
			panic(err)
		}
	}
}
