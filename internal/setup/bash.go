package setup

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/vlatan/vps-setup/internal/colors"
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

func (b *Bash) CreatePrompt() error {

	// Backslashes and quotes escaped
	promptContent := []string{
		"# Using tput for colors and formatting. Reset colors.",
		"tput sgr0",
		"",
		"# Set color variables",
		"bold=\"\\[$(tput bold)\\]\"",
		"reset=\"\\[$(tput sgr0)\\]\"",
		"blue=\"\\[$(tput setaf 153)\\]\"",
		"steel_blue=\"\\[$(tput setaf 67)\\]\"",
		"green=\"\\[$(tput setaf 71)\\]\"",
		"orange=\"\\[$(tput setaf 166)\\]\"",
		"red=\"\\[$(tput setaf 167)\\]\"",
		"white=\"\\[$(tput setaf 15)\\]\"",
		"yellow=\"\\[$(tput setaf 228)\\]\"",
		"",
		"# Highlight the user name when logged in as root.",
		"[[ \"${USER}\" == \"root\" ]] && userColor=\"${red}\" || userColor=\"${green}\"",
		"",
		"# Highlight the hostname when connected via SSH.",
		"[[ \"${SSH_TTY}\" ]] && hostColor=\"${red}\" || hostColor=\"${orange}\"",
		"",
		"# Set the default interaction prompt",
		"PS1=\"${bold}\"							                    # bold",
		"PS1+=\"${white}[${yellow}\\A${white}]\"				        # [time]",
		"PS1+=\"[${userColor}\\u${white}@${hostColor}\\h${white}]\"		# [user@host]",
		"PS1+=\"[${steel_blue}\\w${white}]\"				            # [pwd]",
		"PS1+=\"\\n$ ${reset}\"                                       	# new line and reset",
		"export PS1",
		"",
		"# Set the continuation interactive prompt",
		"PS2=\"${yellow} → ${reset}\"",
		"export PS2",
	}

	data := []byte(strings.Join(promptContent, "\n") + "\n")
	if err := utils.WriteFile(b.root, b.prompt, data); err != nil {
		return err
	}

	return nil
}

// FormatBashrc appends content to .bashrc file.
// Enables aliases and the custom prompt.
func (b *Bash) FormatBashrc() error {

	bashrcContent := []string{
		"# use aliases from file if any",
		"if [ -f ~/.bash_aliases ]; then",
		"    . ~/.bash_aliases",
		"fi",
		"",
		"# use custom prompt from file if any",
		"if [ -f ~/.bash_prompt ]; then",
		"    . ~/.bash_prompt",
		"fi",
		"",
		"# increase bash history",
		"export HISTSIZE=10000000",
		"export HISTFILESIZE=10000000",
		"",
		"# set more restrictive user mask",
		"umask 022",
	}

	// Append the content to .bashrc
	flag := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	file, err := b.root.OpenFile(b.bashrc, flag, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data := "\n" + strings.Join(bashrcContent, "\n") + "\n"
	if _, err := file.WriteString(data); err != nil {
		return err
	}

	return nil
}

// FormatBash configures the bash experience
// by creating custom aliases and prompt.
func FormatBash(scanner *bufio.Scanner, home *os.Root) error {

	prompt := "Do you want to format bash? [y/N]: "
	prompt = colors.Yellow(prompt)
	start := strings.ToLower(utils.AskQuestion(prompt, scanner))
	if !slices.Contains([]string{"yes", "y"}, start) {
		return nil
	}

	bash := NewBash(home)
	callables := []func() error{
		bash.CreateAliases,
		bash.CreatePrompt,
		bash.FormatBashrc,
	}

	for _, Callable := range callables {
		if err := Callable(); err != nil {
			return err
		}
	}

	return nil
}
