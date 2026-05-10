package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port int
	DB   *DBConfig
}

type DBConfig struct {
	Dialect  string
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	Charset  string
}

func GetConfig() *Config {
	return &Config{
		Port: getEnvAsInt("APP_PORT", 3000),
		DB: &DBConfig{
			Dialect:  "postgres",
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnvAsInt("DB_PORT", 7432),
			Username: getEnv("DB_USER", "guest"),
			Password: getEnv("DB_PASSWORD", "Guest0000!"),
			Name:     getEnv("DB_NAME", "freellm"),
			Charset:  "utf8",
		},
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}
