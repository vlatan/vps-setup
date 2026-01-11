package settings

import (
	"fmt"
	"os"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// AutoRestart creates a new config file /etc/needrestart/conf.d/no-prompt.conf
// Adds $nrconf{restart} = 'a'; to that file.
// This will set services to automatically restart after update/upgrade.
func AutoRestart(etc *os.Root) error {
	msg := colors.Yellow("Seting up services autorestart...")
	fmt.Println(msg)

	// Write to file
	name := "needrestart/conf.d/no-prompt.conf"
	data := []byte("$nrconf{restart} = 'a';\n")
	if err := utils.WriteFile(etc, name, data); err != nil {
		return err
	}

	cmd := utils.Command("apt-get", "update")
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
