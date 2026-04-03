package config

import (
	"errors"
	"os"
	"time"
)

type Config struct {
	RefreshTokenTTL time.Duration
}

func GetConfig() (*Config, error) {
	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		return nil, errors.New("POSTGRES_URL is not set")
	}

	strTtl := os.Getenv("REFRESH_TOKEN_TTL")
	if strTtl == "" {
		strTtl = "720h"
	}

	ttl, err := time.ParseDuration(strTtl)
	if err != nil {
		return nil, err
	}

	conf := Config{
		RefreshTokenTTL: ttl,
	}

	return &conf, nil
}
