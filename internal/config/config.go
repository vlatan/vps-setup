package config

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

type Config struct {

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
func New() *Config {

	// Parse the config from the environment
	var cfg Config
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

	// Open /home as root
	home, err := os.OpenRoot("/home")
	if err != nil {
		utils.Exit(err)
	}

	cfg.Home = home

	return &cfg
}

// Close closes opened roots
func (c *Config) Close() error {
	return errors.Join(c.Etc.Close(), c.Home.Close())
}
