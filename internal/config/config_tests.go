package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

func NewTest() (Config, error) {
	configPath := os.Getenv("CONFIG_PATH_TEST")
	if configPath == "" {
		return Config{}, fmt.Errorf("CONFIG_PATH_TEST is not set")
	}

	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		return Config{}, fmt.Errorf("config file %s not found", configPath)
	}

	var cfg Config

	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("cleanenv.ReadConfig: %w", err)
	}

	return cfg, nil
}
