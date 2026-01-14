package setup

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// AddUser adds new user and makes that user sudoer.
// This method will modify the setup state,
// namely s.Username and s.Home.
func (s *Setup) AddUser() error {

	prompt := colors.Yellow("Provide username: ")
	for {
		// Check for env var username first
		if s.Username != "" {
			break
		}
		s.Username = utils.AskQuestion(prompt, s.Scanner)
	}

	var cmds []*exec.Cmd
	if s.Password != "" {
		cmd1 := exec.Command("adduser", "--gecos", "", "--disabled-password", s.Username)
		cmd2 := exec.Command("chpasswd")
		cmd2.Stdin = strings.NewReader(fmt.Sprintf("%s:%s", s.Username, s.Password))
		cmds = append(cmds, cmd1, cmd2)
	} else {
		cmds = append(cmds, exec.Command("adduser", "--gecos", "", s.Username))
	}

	// Add user to sudo group, make them sudoer
	cmds = append(cmds, exec.Command("adduser", s.Username, "sudo"))

	fmt.Println("Adding new user...")
	for _, cmd := range cmds {
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	// Get and set user's UID and GID
	u, err := user.Lookup(s.Username)
	if err != nil {
		return err
	}

	s.uid, err = strconv.Atoi(u.Uid)
	if err != nil {
		return err
	}

	s.Gid, err = strconv.Atoi(u.Gid)
	if err != nil {
		return err
	}

	// Open the user dir as root
	userDir := filepath.Join("/home", s.Username)
	home, err := os.OpenRoot(userDir)
	if err != nil {
		return err
	}

	// Provide the user opened root to the setup
	s.Home = home

	return nil
}
