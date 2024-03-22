package configs

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpPort  string `env:"HTTP_PORT"`
	HttpsPort string `env:"HTTPS_PORT"`
	Pg
	Auth
}

type Pg struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBname   string `env:"POSTGRES_DB"`
}

type Auth struct {
	Login    string `env:"AUTH_LOGIN"`
	Password string `env:"AUTH_PASSWORD"`
}

func NewConfig() (Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig("configs/.env", &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("cleanenv.ReadEnv: %w", err)
	}

	return cfg, nil
}
