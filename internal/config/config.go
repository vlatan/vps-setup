package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	// Autoload env vars from .env file
	// Will not override existing env vars
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
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
}

// New creates new config object
func New() *Config {

	// Parse the config from the environment
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		fmt.Println("Was unable to parse environment variables")
	}

	return &cfg
}
