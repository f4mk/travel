package config

import (
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type API struct {
	HostName       string        `env:"API_HOST_NAME,required"`
	Port           string        `env:"API_PORT,required"`
	RequestTimeout time.Duration `env:"API_REQUEST_TIMEOUT" envDefault:"5s"`
}
type Target struct {
	HostName string `env:"TARGET_HOST_NAME,required"`
	Port     string `env:"TARGET_PORT,required"`
}
type Service struct {
	ServiceName string `env:"SERVICE_NAME,required"`
	DomainName  string `env:"DOMAIN_NAME,required"`
}

type Log struct {
	LogLevel int `env:"LOG_LEVEL,required"`
}

type Config struct {
	Environment string `env:"ENVIRONMENT,required"`
	API         API
	Target      Target
	Service     Service
	Log         Log
}

func New(configPath string) (*Config, error) {

	if err := godotenv.Load(configPath); err != nil {
		return nil, err
	}

	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
