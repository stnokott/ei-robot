package config

import (
	"fmt"
	"os"
)

const (
	envTelegramToken string = "TELEGRAM_TOKEN"
)

type Config struct {
	TelegramToken string
	DbDir         string
}

const dataDir = "/var/lib/data"

func New() (*Config, error) {
	telegramToken := os.Getenv(envTelegramToken)
	if telegramToken == "" {
		return nil, fmt.Errorf("required environment variable %s not present", envTelegramToken)
	}

	if err := os.MkdirAll(dataDir, 0666); err != nil {
		return nil, fmt.Errorf("could not create data dir: %w", err)
	}

	return &Config{
		TelegramToken: telegramToken,
		DbDir:         dataDir,
	}, nil
}
