package setup

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/vlatan/vps-setup/internal/utils"
)

// SetupPostfix installs Postfix and sets up SMTP relay
func (s *Setup) SetupPostfix() error {

	needed := map[string]string{
		"postfix mailname": s.PostfixMailname,
		"smtp host":        s.SMTPHost,
		"smtp port":        s.SMTPPort,
		"smtp username":    s.SMTPUsername,
		"smtp password":    s.SMTPPassword,
	}

	for key, value := range needed {
		if value == "" {
			return fmt.Errorf("%s not found", key)
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
		cmd = utils.Command("postconf", "-e", conf)
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
