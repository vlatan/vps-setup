package setup

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/vlatan/vps-setup/internal/utils"
)

// SetupGitRepo sets up a bare git repo
func (s *Setup) SetupGitRepo() error {

	if s.GitRepoName == "" {
		return errors.New("git repo name not found")
	}

	fmt.Println("Setting up Git repo...")
	checkoutDirName := s.GitRepoName

	// Create the checkout dir
	if err := s.Home.MkdirAll(checkoutDirName, 0755); err != nil {
		return err
	}

	// Define repo and checkout dirs absolute paths
	repoDirName := checkoutDirName + ".git"
	repoDirAbsPath := filepath.Join(s.Home.Name(), repoDirName)
	checkoutDirAbsPath := filepath.Join(s.Home.Name(), checkoutDirName)

	// Create the bare repo
	cmd := utils.Command("git", "init", "--bare", "--initial-branch=main", repoDirAbsPath)
	if err := cmd.Run(); err != nil {
		return err
	}

	// Hook content
	hookContent := []string{
		"#!/bin/bash",
		"",
		"# Git is running non-interactively",
		"# so we need to set umask explicitely here",
		"# to preserve the 755 and 644 permissions",
		"umask 022",
		"",
		"# Check out the files in this directory",
		fmt.Sprintf("GIT_WORK_TREE=%s git checkout -f main", checkoutDirAbsPath),
	}

	hookFile := filepath.Join(repoDirName, "hooks", "post-receive")
	data := []byte(strings.Join(hookContent, "\n") + "\n")
	if err := s.Home.WriteFile(hookFile, data, 0644); err != nil {
		return err
	}

	// Asign permissions
	hookFileAbsPath := filepath.Join(s.Home.Name(), hookFile)
	cmds := [][]string{
		{"chmod", "+x", hookFileAbsPath},
		{
			"chown", "-R",
			s.Username + ":" + s.Username,
			checkoutDirAbsPath, repoDirAbsPath,
		},
	}

	for _, cmdArgs := range cmds {
		cmd := utils.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
