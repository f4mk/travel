package config

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type API struct {
	HostName        string        `env:"API_HOST_NAME,required"`
	Port            string        `env:"API_PORT,required"`
	ReadTimeout     time.Duration `env:"API_READ_TIMEOUT" envDefault:"5s"`
	WriteTimeout    time.Duration `env:"API_WRITE_TIMEOUT" envDefault:"5s"`
	ShutdownTimeout time.Duration `env:"API_SHUTDOWN_TIMEOUT" envDefault:"5s"`
	RequestTimeout  time.Duration `env:"API_REQUEST_TIMEOUT" envDefault:"5s"`
	KeyFile         string        `env:"API_KEY_FILE,required"`
	RateLimit       int           `env:"API_RATE_LIMIT,required"`
}

type Service struct {
	ServiceName string `env:"SERVICE_NAME,required"`
	DomainName  string `env:"DOMAIN_NAME,required"`
}

type Debug struct {
	HostName string `env:"DEBUG_HOST_NAME,required"`
	Port     string `env:"DEBUG_PORT,required"`
}

type Auth struct {
	KeyPath         string        `env:"AUTH_KEY_PATH,required"`
	AuthDuration    time.Duration `env:"AUTH_AUTH_DURATION,required"`
	RefreshDuration time.Duration `env:"AUTH_REFRESH_DURATION,required"`
}

type Cache struct {
	HostName     string `env:"RD_HOST_NAME,required"`
	Port         string `env:"RD_PORT,required"`
	DB           int    `env:"RD_DB,required"`
	PoolSize     int    `env:"RD_POOL_SIZE,required"`
	MinIdleConns int    `env:"RD_MIN_IDLE_CONNS,required"`
}

type MessageBroker struct {
	HostName string `env:"RABBIT_HOST_NAME,required"`
	Port     string `env:"RABBIT_PORT,required"`
	User     string `env:"RABBIT_NAME,required"`
	Password string `env:"RABBIT_PASSWORD,required"`
}

type MailService struct {
	PrivateKey string `env:"MAIL_PRIVATE_KEY_PATH,required"`
	PublicKey  string `env:"MAIL_PUBLIC_KEY_PATH,required"`
}

type Log struct {
	LogLevel int `env:"LOG_LEVEL,required"`
}

type Telemetry struct {
	Port     string  `env:"OTEL_PORT,required"`
	HostName string  `env:"OTEL_HOST_NAME,required"`
	Route    string  `env:"OTEL_ROUTE,required"`
	Prob     float64 `env:"OTEL_PROB,required"`
}

type DB struct {
	User        string `env:"PG_USER,required"`
	Password    string `env:"PG_PASSWORD,required,unset"`
	HostName    string `env:"PG_HOST_NAME,required"`
	Port        string `env:"PG_PORT,required"`
	DBName      string `env:"PG_DB_NAME,required"`
	MaxIdleConn int    `env:"PG_MAX_IDLE_CONN" envDefault:"2"`
	MaxOpenConn int    `env:"PG_MAX_OPEN_CONN" envDefault:"0"`
	DisableTLS  bool   `env:"PG_DISABLE_TLS,required"`
}

type Config struct {
	Environment   string `env:"ENVIRONMENT,required"`
	Service       Service
	Log           Log
	API           API
	Debug         Debug
	Auth          Auth
	DB            DB
	Cache         Cache
	MessageBroker MessageBroker
	MailService   MailService
	Telemetry     Telemetry
}

func New(configPath string) (*Config, error) {

	if err := godotenv.Load(configPath); err != nil {
		return nil, err
	}

	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	// TODO: refactor this
	privateKeyContent, err := os.ReadFile(cfg.MailService.PrivateKey)
	if err != nil {
		fmt.Printf("Error reading private key: %v\n", err)
		os.Exit(1)
	}

	publicKeyContent, err := os.ReadFile(cfg.MailService.PublicKey)
	if err != nil {
		fmt.Printf("Error reading public key: %v\n", err)
		os.Exit(1)
	}

	cfg.MailService.PrivateKey = string(privateKeyContent)
	cfg.MailService.PublicKey = string(publicKeyContent)

	return &cfg, nil
}
