package setup

import (
	"bufio"
	"errors"
	"os"

	"github.com/caarlos0/env/v11"

	// Autoload env vars from .env file
	// Will not override existing env vars
	_ "github.com/joho/godotenv/autoload"
)

type Setup struct {

	// Env vars
	Swappiness      string `env:"SWAPPINESS" envDefault:"20"`
	UbuntuProToken  string `env:"UBUNTU_PRO_TOKEN"`
	Hostname        string `env:"HOSTNAME"`
	Timezone        string `env:"TIMEZONE" envDefault:"UTC"`
	Username        string `env:"USERNAME"`
	Password        string `env:"PASSWORD"`
	SSHPort         string `env:"SSH_PORT" envDefault:"22"`
	SSHPubKey       string `env:"SSH_PUBKEY"`
	PostfixMailname string `env:"POSTFIX_MAILNAME"`
	SMTPHost        string `env:"SMTP_HOST"`
	SMTPPort        string `env:"SMTP_PORT"`
	SMTPUsername    string `env:"SMTP_USERNAME"`
	SMTPPassword    string `env:"SMTP_PASSWORD"`
	GitRepoName     string `env:"GIT_REPO_NAME"`

	// Dynamic fields
	Scanner   *bufio.Scanner
	Etc, Home *os.Root
	Uid, Gid  int
}

// New creates new config object
func New() (*Setup, error) {

	// Parse the config from the environment
	var cfg Setup
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	cfg.Scanner = bufio.NewScanner(os.Stdin)

	// Open /etc as root
	etc, err := os.OpenRoot("/etc")
	if err != nil {
		return nil, err
	}

	cfg.Etc = etc
	return &cfg, nil
}

// Close closes opened roots
func (s *Setup) Close() error {
	var errs []error

	if s.Etc != nil {
		errs = append(errs, s.Etc.Close())
	}

	if s.Home != nil {
		errs = append(errs, s.Home.Close())
	}

	return errors.Join(errs...)
}
