package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DatabaseURL        string
	Port               string
	LogLevel           string
	BaselineWindowDays int
	AnomalyStrictMode  bool
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL:        strings.TrimSpace(os.Getenv("DATABASE_URL")),
		Port:               strings.TrimSpace(os.Getenv("PORT")),
		LogLevel:           strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL"))),
		BaselineWindowDays: 30,
		AnomalyStrictMode:  false,
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if cfg.Port == "" {
		cfg.Port = "8000"
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	if v := strings.TrimSpace(os.Getenv("BASELINE_WINDOW_DAYS")); v != "" {
		if days, err := strconv.Atoi(v); err == nil && days > 0 {
			cfg.BaselineWindowDays = days
		}
	}

	if v := strings.TrimSpace(os.Getenv("ANOMALY_STRICT_MODE")); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.AnomalyStrictMode = b
		}
	}

	return cfg, nil
}
