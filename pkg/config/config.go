package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	TelegramBotToken string `yaml:"telegram_token"`
}

func NewConfig(cfgPath string) (Config, error) {
	fd, err := os.Open(cfgPath)
	if err != nil {
		return Config{}, fmt.Errorf("could not open config path: %w", err)
	}

	defer fd.Close()

	var cfg Config
	if err := yaml.NewDecoder(fd).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("could not decode config data: %w", err)
	}

	return cfg, nil
}
