package config

import "os"

type Config struct {
	Addr string
	Env  string
}

func Load() *Config {
	return &Config{
		Addr: getEnv("SERVER_ADDR", ":8080"),
		Env:  getEnv("APP_ENV", "dev"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
