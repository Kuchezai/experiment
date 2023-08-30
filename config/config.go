package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP HTTP `yaml:"http"`
		DB   DB   `yaml:"db"`
	}
	HTTP struct {
		Address     string        `yaml:"address"`
		Timeout     time.Duration `yaml:"timeout"`
		IdleTimeout time.Duration `yaml:"idle_timeout"`
		JWTSecret   string        `yaml:"jwt_secret"`
	}
	DB struct {
		Host         string        `env:"DB_HOST"`
		Port         string        `env:"DB_PORT"`
		User         string        `env:"DB_USER"`
		Pass         string        `env:"DB_PASSWORD"`
		Name         string        `env:"DB_NAME"`
		PoolSize     int           `yaml:"pool-size"`
		ConnAttempts int           `yaml:"conn-attempts"`
		ConnTimeout  time.Duration `yaml:"conn-timeout"`
	}
)

const (
	configPath = "config/config.yaml"
)

func New() (*Config, error) {
	op := "config.New"

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &cfg, nil
}
