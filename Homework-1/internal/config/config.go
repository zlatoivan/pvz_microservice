package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server `yaml:"server"`
	Pg     `yaml:"postgres"`
}

type Pg struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBname   string `yaml:"dbname"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Server struct {
	HttpPort  string `yaml:"http_port"`
	HttpsPort string `yaml:"https_port"`
	Login     string `yaml:"pvz_auth_login"`
	Password  string `yaml:"pvz_auth_password"`
}

func New() (Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return Config{}, fmt.Errorf("CONFIG_PATH is not set")
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

	fmt.Println(cfg)

	return cfg, nil
}
