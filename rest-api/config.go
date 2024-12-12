package main

import (
	"os"
)

type Configuration struct {
	Port string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Port: getEnvOrDefault("PORT", "3000"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
