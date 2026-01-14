package setup

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/vlatan/vps-setup/internal/colors"
	"github.com/vlatan/vps-setup/internal/utils"
)

// SetupPostfix installs Postfix and sets up SMTP relay
func (s *Setup) SetupPostfix() error {

	prompt := "Do you want to setup Postfix SMTP relay? [y/N]: "
	prompt = colors.Yellow(prompt)
	start := strings.ToLower(utils.AskQuestion(prompt, s.Scanner))
	if !slices.Contains([]string{"yes", "y"}, start) {
		return nil
	}

	for {
		if s.PostfixMailname != "" {
			break
		}
		s.PostfixMailname = utils.AskQuestion(
			"MAILNAME (domain): ",
			s.Scanner,
		)
	}

	for {
		if s.SMTPHost != "" {
			break
		}
		s.SMTPHost = utils.AskQuestion("SMTP_HOST: ", s.Scanner)
	}

	for {
		if s.SMTPPort != "" {
			break
		}
		s.SMTPPort = utils.AskQuestion("SMTP_PORT: ", s.Scanner)
	}

	for {
		if s.SMTPUsername != "" {
			break
		}
		var err error
		s.SMTPUsername, err = utils.AskPassword("SMTP_USERNAME: ")
		if err != nil {
			return err
		}
	}

	for {
		if s.SMTPPassword != "" {
			break
		}
		var err error
		s.SMTPPassword, err = utils.AskPassword("SMTP_PASSWORD: ")
		if err != nil {
			return err
		}
	}

	fmt.Println("Setting up Postfix SMTP relay...")

	stdins := []string{
		fmt.Sprintf("postfix postfix/mailname string %s\n", s.PostfixMailname),
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

	// Write to the file
	name := "postfix/sasl/sasl_passwd"
	data := fmt.Sprintf(
		"[%s]:%s %s:%s",
		s.SMTPHost, s.SMTPPort, s.SMTPUsername, s.SMTPPassword,
	)
	if err := utils.WriteFile(s.Etc, name, []byte(data)); err != nil {
		return err
	}

	saslPaswdFile := filepath.Join(s.Etc.Name(), name)
	cmd = utils.Command("postmap", saslPaswdFile)
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = utils.Command("chmod", "0600", saslPaswdFile, saslPaswdFile+".db")
	if err := cmd.Run(); err != nil {
		return err
	}

	postConfs := []string{
		fmt.Sprintf("myhostname = %s", s.PostfixMailname),
		fmt.Sprintf("relayhost = [%s]:%s", s.SMTPHost, s.SMTPPort),
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
