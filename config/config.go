package config

import "os"

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

func Load() *Config {
	return &Config{
		SMTPHost:     getEnv("SMTP_HOST", "smtp.mailtrap.io"),
		SMTPPort:     587,
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASS", ""),
		FromEmail:    getEnv("FROM_EMAIL", "no-reply@yourdomain.com"),
		FromName:     getEnv("FROM_NAME", "Email Service"),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:    getEnv("REDIS_PASS", ""),
		WorkerCount:  5,
		MaxRetries:   3,
	}
}

func getEnv(key, falback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return falback
}
