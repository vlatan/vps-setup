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

// InstallPostfix installs Postfix and sets up SMTP relay
func InstallPostfix(scanner *bufio.Scanner, etc *os.Root) error {

	prompt := "Do you want to install Postfix and setup SMTP relay? [y/n]: "
	prompt = colors.Yellow(prompt)
	start := strings.ToLower(utils.AskQuestion(prompt, scanner))
	if !slices.Contains([]string{"yes", "y"}, start) {
		return nil
	}

	var domain string
	for {
		domain = utils.AskQuestion("MAILNAME (domain): ", scanner)
		if domain != "" {
			break
		}
	}

	var smtpHost string
	for {
		smtpHost = utils.AskQuestion("SMTP_HOST: ", scanner)
		if smtpHost != "" {
			break
		}
	}

	var smtpPort string
	for {
		smtpPort = utils.AskQuestion("SMTP_PORT: ", scanner)
		if smtpPort != "" {
			break
		}
	}

	var smtpUsername string
	for {
		smtpUsername, err := utils.AskPassword("SMTP_USERNAME: ")
		if err != nil {
			return err
		}

		if smtpUsername != "" {
			break
		}
	}

	var smtpPassword string
	for {
		smtpPassword, err := utils.AskPassword("SMTP_PASSWORD: ")
		if err != nil {
			return err
		}

		if smtpPassword != "" {
			break
		}
	}

	stdins := []string{
		fmt.Sprintf("postfix postfix/mailname string %s\n", domain),
		"postfix postfix/main_mailer_type string 'Internet Site'\n",
	}

	for _, stdin := range stdins {
		cmd := utils.Command("debconf-set-selections")
		cmd.Stdin = strings.NewReader(stdin)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	cmd := utils.Command(
		"apt-get",
		"install",
		"-y",
		"mailutils",
		"postfix",
		"ca-certificates",
		"libsasl2-modules",
	)

	if err := cmd.Run(); err != nil {
		return err
	}

	// Create dirs that do not exist in the file path
	name := "postfix/sasl/sasl_passwd"
	if err := etc.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}

	// Write to the file
	data := fmt.Sprintf("[%s]:%s %s:%s", smtpHost, smtpPort, smtpUsername, smtpPassword)
	if err := etc.WriteFile(name, []byte(data), 0644); err != nil {
		return err
	}

	saslPaswdFile := filepath.Join(etc.Name(), name)
	cmd = utils.Command("postmap", saslPaswdFile)
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = utils.Command("chmod", "0600", saslPaswdFile, saslPaswdFile+".db")
	if err := cmd.Run(); err != nil {
		return err
	}

	postConfs := []string{
		fmt.Sprintf("myhostname = %s", domain),
		fmt.Sprintf("relayhost = [%s]:%s", smtpHost, smtpPort),
		"smtp_sasl_auth_enable = yes",
		"smtp_sasl_security_options = noanonymous",
		fmt.Sprintf("smtp_sasl_password_maps = hash:%s", saslPaswdFile),
		"smtp_use_tls = yes",
		"smtp_tls_security_level = encrypt",
		"smtp_tls_note_starttls_offer = yes",
		"smtp_tls_CAfile = /etc/ssl/certs/ca-certificates.crt",
	}

	for _, conf := range postConfs {
		cmd = utils.Command("postconf", "-e", fmt.Sprintf("%q", conf))
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	cmd = utils.Command("systemctl", "restart", "postfix")
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
