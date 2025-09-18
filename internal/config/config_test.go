package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test default configuration
	cfg := LoadConfig()
	if cfg.Environment != "prod" {
		t.Errorf("Expected environment 'prod', got %s", cfg.Environment)
	}
	if cfg.DBPath != "./notebook.db" {
		t.Errorf("Expected DB path './notebook.db', got %s", cfg.DBPath)
	}
	if cfg.Port != "8003" {
		t.Errorf("Expected port '8003', got %s", cfg.Port)
	}

	// Test environment variable override
	os.Setenv("ENV", "dev")
	os.Setenv("PORT", "3000")
	defer os.Unsetenv("ENV")
	defer os.Unsetenv("PORT")

	cfg = LoadConfig()
	if cfg.Environment != "dev" {
		t.Errorf("Expected environment 'dev', got %s", cfg.Environment)
	}
	if cfg.Port != "3000" {
		t.Errorf("Expected port '3000', got %s", cfg.Port)
	}
}

func TestConfigMethods(t *testing.T) {
	// Test production mode
	cfg := &Config{Environment: "prod"}
	if cfg.IsDev() {
		t.Error("Expected IsDev() to return false for prod environment")
	}
	if cfg.IsAdmin() {
		t.Error("Expected IsAdmin() to return false for prod environment")
	}

	// Test development mode
	cfg = &Config{Environment: "dev"}
	if !cfg.IsDev() {
		t.Error("Expected IsDev() to return true for dev environment")
	}
	if !cfg.IsAdmin() {
		t.Error("Expected IsAdmin() to return true for dev environment")
	}
}
