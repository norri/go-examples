package main

import (
	"log/slog"
	"os"
)

type Configuration struct {
	Port        string
	DatabaseURL string
}

func NewConfiguration() *Configuration {
	dbURL := getEnvOrDefault("DATABASE_URL", "")
	if dbURL == "" {
		slog.Warn("DATABASE_URL is not set")
	}
	return &Configuration{
		Port:        getEnvOrDefault("PORT", "3000"),
		DatabaseURL: dbURL,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
