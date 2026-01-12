package settings

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// SetupGitRepo sets up a bare git repo
func SetupGitRepo(username string, scanner *bufio.Scanner, home *os.Root) error {

	prompt := "Do you want to create a bare git repo [y/n]: "
	prompt = colors.Yellow(prompt)
	start := strings.ToLower(utils.AskQuestion(prompt, scanner))
	if !slices.Contains([]string{"yes", "y"}, start) {
		return nil
	}

	var dirName string
	for {
		dirName = utils.AskQuestion("Provide repo directory name: ", scanner)
		if dirName != "" {
			break
		}
	}

	// Absolute dir paths
	checkoutDirAbsPath := filepath.Join(home.Name(), dirName)
	repoDirAbsPath := checkoutDirAbsPath + ".git"

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

	hookFile := filepath.Join(dirName+".git", "hooks", "post-receive")
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
