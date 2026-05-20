package main

import (
	"os"
	"testing"
)

func TestGetEnvReturnsExistingValue(t *testing.T) {
	t.Setenv("TEST_ENV_VALUE", "custom")

	value := getEnv("TEST_ENV_VALUE", "fallback")

	if value != "custom" {
		t.Errorf("expected custom, got %s", value)
	}
}

func TestGetEnvReturnsFallback(t *testing.T) {
	_ = os.Unsetenv("TEST_ENV_MISSING")

	value := getEnv("TEST_ENV_MISSING", "fallback")

	if value != "fallback" {
		t.Errorf("expected fallback, got %s", value)
	}
}