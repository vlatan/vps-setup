package settings

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// AutoRestart creates a new config file /etc/needrestart/conf.d/no-prompt.conf
// Adds $nrconf{restart} = 'a'; to that file.
// This will set services to automatically restart after update/upgrade.
func AutoRestart(etc *os.Root) error {
	msg := colors.Yellow("Seting up services autorestart...")
	fmt.Println(msg)

	// Create dirs that do not exist in the file path
	name := "needrestart/conf.d/no-prompt.conf"
	if err := etc.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}

	data := "$nrconf{restart} = 'a';\n"
	if err := etc.WriteFile(name, []byte(data), 0644); err != nil {
		return err
	}

	if err := utils.RunCommand("apt-get", "update"); err != nil {
		return err
	}

	return nil
}
