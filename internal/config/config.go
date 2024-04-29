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

type Server struct {
	HttpPort           string   `yaml:"http_port" env-default:"9000"`
	HttpsPort          string   `yaml:"https_port" env-default:"9001"`
	GatewayPort        string   `yaml:"gateway_port" env-default:"9002"`
	MetricsPort        string   `yaml:"metrics_port" env-default:"9003"`
	GrpcPort           string   `yaml:"grpc_port" env-default:"50051"`
	TracerExporterPort string   `yaml:"tracer_exporter_port" env-default:"4318"`
	TracerJaegerUIPort string   `yaml:"tracer_jaeger_ui_port" env-default:"16686"`
	PVZLogin           string   `yaml:"pvz_auth_login" env-default:""`
	PVZPassword        string   `yaml:"pvz_auth_password" env-default:""`
	OrderLogin         string   `yaml:"order_auth_login" env-default:""`
	OrderPassword      string   `yaml:"order_auth_password" env-default:""`
	Brokers            []string `yaml:"brokers" env-default:""`
	Topic              string   `yaml:"topic" env-default:""`
	Redis              string   `yaml:"redis" env-default:""`
}

type Pg struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	DBname   string `yaml:"dbname" env-default:"postgres"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
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

	return cfg, nil
}
