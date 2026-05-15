package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
	RedisAddr    string
	RedisPass    string
	WorkerCount  int
	MaxRetries   int
}

func Load() (*Config, error) {
	// Best effort: load .env for local development.
	// Environment variables already set in the OS still take precedence.
	_ = godotenv.Load(".env")

	requiredVars := []string{
		"SMTP_HOST",
		"SMTP_PORT",
		"FROM_EMAIL",
		"FROM_NAME",
		"REDIS_ADDR",
	}

	missing := make([]string, 0)
	for _, key := range requiredVars {
		if strings.TrimSpace(os.Getenv(key)) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}

	smtpPort, err := parseIntEnv("SMTP_PORT", 587)
	if err != nil {
		return nil, err
	}
	workerCount, err := parseIntEnv("WORKER_COUNT", 5)
	if err != nil {
		return nil, err
	}
	maxRetries, err := parseIntEnv("MAX_RETRIES", 3)
	if err != nil {
		return nil, err
	}

	return &Config{
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     smtpPort,
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASS", ""),
		FromEmail:    getEnv("FROM_EMAIL", ""),
		FromName:     getEnv("FROM_NAME", ""),
		RedisAddr:    getEnv("REDIS_ADDR", ""),
		RedisPass:    getEnv("REDIS_PASS", ""),
		WorkerCount:  workerCount,
		MaxRetries:   maxRetries,
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseIntEnv(key string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %q must be a number", key, value)
	}
	if n <= 0 {
		return 0, fmt.Errorf("invalid %s: %d must be > 0", key, n)
	}
	return n, nil
}
