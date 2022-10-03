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
}

func New() (*Config, error) {
	telegramToken := os.Getenv(envTelegramToken)
	if telegramToken == "" {
		return nil, fmt.Errorf("required environment variable %s not present!", envTelegramToken)
	}
	return &Config{
		TelegramToken: telegramToken,
	}, nil
}
