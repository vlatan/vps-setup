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

	// Use Satellite system - it relays everything and doesn't treat domain as local
	fmt.Println("Setting up Postfix SMTP relay...")
	stdins := []string{
		fmt.Sprintf("postfix postfix/mailname string %s\n", s.PostfixMailname),
		"postfix postfix/main_mailer_type string 'Satellite system'\n",
		fmt.Sprintf("postfix postfix/relayhost string [%s]:%s\n", s.SMTPHost, s.SMTPPort),
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

	// Create mailutils config to use correct email domain
	mailutilsConf := fmt.Sprintf("address { email-domain %s; };", s.PostfixMailname)
	if err := s.Etc.WriteFile("mailutils.conf", []byte(mailutilsConf), 0644); err != nil {
		return err
	}

	// Make parent directories
	name := "postfix/sasl/sasl_passwd"
	if err := s.Etc.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}

	// Data for the sasl_passwd file
	data := fmt.Sprintf(
		"[%s]:%s %s:%s",
		s.SMTPHost, s.SMTPPort, s.SMTPUsername, s.SMTPPassword,
	)

	// Create sasl_passwd file
	if err := s.Etc.WriteFile(name, []byte(data), 0644); err != nil {
		return err
	}

	// Generate sasl_passwd.db file
	absName := filepath.Join(s.Etc.Name(), name)
	cmd = utils.Command("postmap", absName)
	if err := cmd.Run(); err != nil {
		return err
	}

	// Protect the sasl_passwd files
	for _, file := range []string{name, name + ".db"} {
		if err := s.Etc.Chmod(file, 0600); err != nil {
			return err
		}
	}

	// Configure SMTP relay settings
	postConfs := []string{
		"mydestination = $myhostname, localhost.localdomain, localhost",
		fmt.Sprintf("relayhost = [%s]:%s", s.SMTPHost, s.SMTPPort),
		"smtp_sasl_auth_enable = yes",
		"smtp_sasl_security_options = noanonymous",
		fmt.Sprintf("smtp_sasl_password_maps = hash:%s", absName),
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
