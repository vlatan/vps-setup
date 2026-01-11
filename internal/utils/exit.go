package utils

import (
	"fmt"
	"os"

	"github.com/vlatan/vps-setup/internal/colors"
)

func Exit(err error) {
	msg := colors.Red("Setup was interrupted. The machine is in a 'dirty' state.")
	fmt.Fprintln(os.Stderr, msg)
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
