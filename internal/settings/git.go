package settings

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// SetupGitRepo sets up a bare git repo
func SetupGitRepo(username string, scanner *bufio.Scanner, home *os.Root) error {

	prompt := "Provide Git repo name [optiona]: "
	prompt = colors.Yellow(prompt)
	checkoutDirName := utils.AskQuestion(prompt, scanner)

	if checkoutDirName == "" {
		return nil
	}

	// Create the checkout dir
	if err := home.MkdirAll(checkoutDirName, 0755); err != nil {
		return err
	}

	// Define repo and checkout dirs absolute paths
	repoDirName := checkoutDirName + ".git"
	repoDirAbsPath := filepath.Join(home.Name(), repoDirName)
	checkoutDirAbsPath := filepath.Join(home.Name(), checkoutDirName)

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
	if err := utils.WriteFile(home, hookFile, data); err != nil {
		return err
	}

	// Asign permissions
	hookFileAbsPath := filepath.Join(home.Name(), hookFile)
	cmds := [][]string{
		{"chmod", "+x", hookFileAbsPath},
		{"chown", "-R", username + ":" + username, checkoutDirAbsPath, repoDirAbsPath},
	}

	for _, cmdArgs := range cmds {
		cmd := utils.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
