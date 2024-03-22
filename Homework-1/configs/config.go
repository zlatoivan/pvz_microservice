package configs

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpPort  string `env:"HTTP_PORT"`
	HttpsPort string `env:"HTTPS_PORT"`
	Pg
	User
}

type Pg struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBname   string `env:"POSTGRES_DB"`
}

type User struct {
	Login    string `env:"USER_LOGIN"`
	Password string `env:"USER_PASSWORD"`
}

func NewConfig() (Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig("configs/.env", &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("cleanenv.ReadEnv: %w", err)
	}

	return cfg, nil
}
