package config

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
	EnvDev   = "dev"
)

type Config struct {
	Enviremant string `env:"ENVIRONMENT" env-default:"local"`
	App        App
	DB         DB
	HTTPServer HTTPServer
}

type App struct {
	AppLogin   string        `env:"APP_USERNAME" env-default:"root"`
	AppPassord string        `env:"APP_PASSWORD" env-default:"root"`
	TockenTTL  time.Duration `env:"APP_TOCKEN_TTL" env-default:"1h"`
	SigningKey string        `env:"-"`
}

type DB struct {
	Provider string `env:"DB_PROVIDER" env-default:"postgres"`
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     int    `env:"DB_PORT" env-default:"5432"`
	Username string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASS" env-required:"true"`
	DBName   string `env:"DB_NAME" env-required:"true"`
	SSLMode  string `env:"DB_SSLMode" env-default:"disable"`
	Pool     DBPool
	Retry    DBRetry
}

type DBPool struct {
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" env-default:"10"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" env-default:"5"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" env-default:"30m"`
	MaxRequestTime  time.Duration `env:"DB_MAX_REQUEST_TIME" env-default:"5s"`
}

type DBRetry struct {
	MaxAttempts int           `env:"DB_MAX_ATTEMPTS" env-default:"10"`
	Delay       time.Duration `env:"DB_DELAY" env-default:"2s"`
	MaxDelay    time.Duration `env:"DB_MAX_DELAY" env-default:"10s"`
}

type HTTPServer struct {
	Address         string        `env:"VHOST" env-required:"true"`
	Port            string        `env:"VHOST_PORT" env-required:"true"`
	ReadTimeout     time.Duration `env:"SERV_READ_TIMEOUT" env-default:"4s"`
	WriteTimeout    time.Duration `env:"SERV_WRITE_TIMEOUT" env-default:"4s"`
	ShutdownTimeout time.Duration `env:"SERV_SHUTDOWN_TIMEOUT" env-default:"4s"`
	IdleTimeout     time.Duration `env:"SERV_IDLE_TIMEOUT" env-default:"60s"`
}

func LoadConfig(path string) (*Config, error) {
	var cfg Config
	if path == "" {
		return &cfg, errors.New("Config file is not set")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &cfg, fmt.Errorf("%w. Config file %s does not exist", err, path)
	}
	err := godotenv.Load(path)
	if err != nil {
		return &cfg, fmt.Errorf("%w. Cannot read the config file %s", err, path)
	}
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return &cfg, fmt.Errorf("%w. Cannot read from environment variables", err)
	}
	cfg.App.SigningKey, err = generateRandomString(16)
	if err != nil {
		return &cfg, fmt.Errorf("%w. Cannot generate SigningKey", err)
	}
	return &cfg, nil
}

func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	return string(b), nil
}

func IsLocalEnvierement(cfg *Config) bool {
	return cfg.Enviremant == EnvLocal
}

func IsDevEnvierement(cfg *Config) bool {
	return cfg.Enviremant == EnvDev
}

func IsProdEnvierement(cfg *Config) bool {
	return cfg.Enviremant == EnvProd
}
