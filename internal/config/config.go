package config

import (
	"fmt"
	"os"
	"path"
)

const (
	envTelegramToken string = "TELEGRAM_TOKEN"
)

type Config struct {
	TelegramToken string
	DbDir         string
}

func New() (*Config, error) {
	telegramToken := os.Getenv(envTelegramToken)
	if telegramToken == "" {
		return nil, fmt.Errorf("required environment variable %s not present", envTelegramToken)
	}

	dbDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	dbDir = path.Join(dbDir, "eirobot-data")

	return &Config{
		TelegramToken: telegramToken,
		DbDir:         dbDir,
	}, nil
}
