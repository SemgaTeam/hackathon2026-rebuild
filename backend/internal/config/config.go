package config

import (
	"errors"
	"os"
)

type Config struct {
}

func GetConfig() (*Config, error) {
	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		return nil, errors.New("POSTGRES_URL is not set")
	}

	conf := Config{
	}

	return &conf, nil
}
