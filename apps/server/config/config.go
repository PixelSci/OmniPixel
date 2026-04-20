package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr string
	Env  string
	DB   DBConfig
	Auth AuthConfig
}

type AuthConfig struct {
	JWTSecret string
	JWTTTL    time.Duration
	JWTIssuer string
}

type DBConfig struct {
	DSN             string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	AutoMigrate     bool
}

func Load() *Config {
	return &Config{
		Addr: getEnv("SERVER_ADDR", ":8080"),
		Env:  getEnv("APP_ENV", "dev"),
		DB: DBConfig{
			DSN:             getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/omni_pixel?sslmode=disable"),
			MaxConns:        int32(getEnvInt("DB_MAX_CONNS", 20)),
			MinConns:        int32(getEnvInt("DB_MIN_CONNS", 2)),
			MaxConnLifetime: getEnvDuration("DB_MAX_CONN_LIFETIME", time.Hour),
			MaxConnIdleTime: getEnvDuration("DB_MAX_CONN_IDLE_TIME", 30*time.Minute),
			AutoMigrate:     getEnvBool("DB_AUTO_MIGRATE", true),
		},
		Auth: AuthConfig{
			JWTSecret: getEnv("JWT_SECRET", "dev-insecure-secret-change-me"),
			JWTTTL:    getEnvDuration("JWT_TTL", 24*time.Hour),
			JWTIssuer: getEnv("JWT_ISSUER", "omni-pixel"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
