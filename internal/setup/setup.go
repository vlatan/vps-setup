package setup

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/vlatan/vps-setup/internal/utils"

	// Autoload env vars from .env file
	// Will not override existing env vars
	_ "github.com/joho/godotenv/autoload"
)

type Setup struct {

	// Env vars
	Swappiness      string `env:"SWAPPINESS"`
	UbuntuProToken  string `env:"UBUNTU_PRO_TOKEN"`
	Hostname        string `env:"HOSTNAME"`
	Timezone        string `env:"TIMEZONE"`
	Username        string `env:"USERNAME"`
	SSHPort         string `env:"SSH_PORT"`
	PostfixMailname string `env:"POSTFIX_MAILNAME"`
	SMTPHost        string `env:"SMTP_HOST"`
	SMTPPort        string `env:"SMTP_PORT"`
	SMTPUsername    string `env:"SMTP_USERNAME"`
	SMTPPassword    string `env:"SMTP_PASSWORD"`
	GitRepoName     string `env:"GIT_REPO_NAME"`

	// Dynamic fields
	Scanner *bufio.Scanner
	Etc     *os.Root
	Home    *os.Root
}

// New creates new config object
func New() *Setup {

	// Parse the config from the environment
	var cfg Setup
	if err := env.Parse(&cfg); err != nil {
		fmt.Println("Unable to parse environment variables")
	}

	cfg.Scanner = bufio.NewScanner(os.Stdin)

	// Open /etc as root
	etc, err := os.OpenRoot("/etc")
	if err != nil {
		utils.Exit(err)
	}

	cfg.Etc = etc

	return &cfg
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
