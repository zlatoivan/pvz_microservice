package configs

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpPort  string `env:"HTTP_PORT"`
	HttpsPort string `env:"HTTPS_PORT"`
	Pg
}

type Pg struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBname   string `env:"POSTGRES_DB"`
}

func NewConfig() (Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig("configs/.env", &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("cleanenv.ReadEnv: %w", err)
	}

	return cfg, nil
}
