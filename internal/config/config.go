package config

import (
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	ServerPort, DBHost, DBPort, DBUser, DBPass, DBName string
	SMTPHost, SMTPPort, SMTPUser, SMTPPass, SMTPFrom   string
	RateRPS, RateBurst                                 int
}

func Load() *Config {
	return &Config{
		ServerPort: env("SERVER_PORT", "8080"), DBHost: env("DB_HOST", "localhost"),
		DBPort: env("DB_PORT", "5432"), DBUser: env("DB_USER", "postgres"),
		DBPass: env("DB_PASSWORD", "secret"), DBName: env("DB_NAME", "kibt_leads"),
		SMTPHost: env("SMTP_HOST", "smtp.gmail.com"), SMTPPort: env("SMTP_PORT", "587"),
		SMTPUser: env("SMTP_USER", ""), SMTPPass: env("SMTP_PASS", ""),
		SMTPFrom: env("SMTP_FROM", "noreply@kibt.ru"),
		RateRPS:  envInt("RATE_LIMIT_RPS", 10), RateBurst: envInt("RATE_LIMIT_BURST", 20),
	}
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func envInt(k string, def int) int {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		slog.Warn("Invalid env int, using default", "key", k, "value", v, "default", def, "error", err)
		return def
	}
	return i
}