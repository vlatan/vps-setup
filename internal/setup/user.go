package setup

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vlatan/vps-setup/internal/utils"
)

// AddUser adds new user and makes that user sudoer.
// This method will modify the setup state,
// namely s.Username and s.Home.
func (s *Setup) AddUser() error {

	if s.Username == "" || s.Password == "" {
		return errors.New("no username and/or password found")
	}

	cmd1 := utils.Command("adduser", "--gecos", "", "--disabled-password", s.Username)
	cmd2 := utils.Command("chpasswd")
	cmd2.Stdin = strings.NewReader(fmt.Sprintf("%s:%s", s.Username, s.Password))
	cmd3 := utils.Command("adduser", s.Username, "sudo")
	cmd4 := utils.Command("adduser", s.Username, "adm")
	cmds := []*exec.Cmd{cmd1, cmd2, cmd3, cmd4}

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

	s.Uid, err = strconv.Atoi(u.Uid)
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
