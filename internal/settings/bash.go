package settings

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vlatan/vps-setup/internal/utils"
)

type Bash struct {
	root    *os.Root
	aliases string
	prompt  string
	bashrc  string
}

func NewBash(root *os.Root) *Bash {
	return &Bash{
		root:    root,
		aliases: ".bash_aliases",
		prompt:  ".bash_prompt",
		bashrc:  ".bashrc",
	}
}

// CreateAliases creates new custom bash aliases
func (b *Bash) CreateAliases() error {

	updateContent := []string{
		"sudo apt update",
		"sudo apt upgrade",
		"sudo apt autoremove",
		"sudo apt autoclean",
	}

	downContent := []string{
		"docker compose down --remove-orphans",
		"docker system prune --force",
	}

	update := strings.Join(updateContent, " && ")
	build := "docker compose up --pull=always --build --detach"
	down := strings.Join(downContent, " && ")

	aliasesContent := []string{
		"# update the repos and upgrade",
		fmt.Sprintf("alias sysupdate=%q", update),
		"",
		"# list files/folders",
		"alias ll='ls -lha'",
		"",
		"# pull images, build and run the containers in background",
		fmt.Sprintf("alias build=%q", build),
		"",
		"# bring down the running containers",
		"# remove dangling images",
		"# remove orphan containers",
		fmt.Sprintf("alias down=%q", down),
	}

	data := []byte(strings.Join(aliasesContent, "\n") + "\n")
	if err := utils.WriteFile(b.root, b.aliases, data); err != nil {
		return err
	}

	return nil
}

// FormatBash configures the bash experience
// by creating custom aliases and prompt.
func FormatBash(username string) error {

	// Open the user dir as root
	userDir := filepath.Join("/home", username)
	home, err := os.OpenRoot(userDir)
	if err != nil {
		return err
	}
	defer home.Close()

	bash := NewBash(home)

	if err = bash.CreateAliases(); err != nil {
		return err
	}

	return nil
}
