package main

import (
	"testing"
)

func TestNewConfiguration(t *testing.T) {
	t.Setenv("PORT", "8080")

	conf := NewConfiguration()

	if conf.Port != "8080" {
		t.Fatalf("Expected 8080, but got %v", conf.Port)
	}
}

func TestNewConfiguration_Defaults(t *testing.T) {
	conf := NewConfiguration()

	if conf.Port != "3000" {
		t.Fatalf("Expected 3000, but got %v", conf.Port)
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	t.Setenv("TEST_ENV", "value")

	value := getEnvOrDefault("TEST_ENV", "default")
	if value != "value" {
		t.Fatalf("Expected value, but got %v", value)
	}

	value = getEnvOrDefault("NON_EXISTENT_ENV", "default")
	if value != "default" {
		t.Fatalf("Expected default, but got %v", value)
	}
}
