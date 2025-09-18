package config

import (
	"os"
)

type Config struct {
    Environment string
    DBPath      string
    ContentDir  string
    SiteBaseURL string
    SiteTitle   string
    Port        string
    ReloadToken string
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
    return &Config{
        Environment: getEnv("ENV", "prod"),
        DBPath:      getEnv("DB_PATH", "./notebook.db"),
        ContentDir:  getEnv("CONTENT_DIR", "./content"),
        SiteBaseURL: getEnv("SITE_BASEURL", "https://notebook.oceanheart.ai"),
        SiteTitle:   getEnv("SITE_TITLE", "Oceanheart Notebook"),
        Port:        getEnv("PORT", "8003"),
        ReloadToken: getEnv("RELOAD_TOKEN", ""),
    }
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsDev returns true if running in development mode
func (c *Config) IsDev() bool {
	return c.Environment == "dev"
}

// IsAdmin returns true if admin endpoints should be enabled
func (c *Config) IsAdmin() bool {
    return c.IsDev()
}

// AllowReload returns true if a reload action is permitted for the provided token.
// In dev, any request is allowed. In prod, RELOAD_TOKEN must be set and match.
func (c *Config) AllowReload(token string) bool {
    if c.IsDev() {
        return true
    }
    if c.ReloadToken == "" {
        return false
    }
    return token == c.ReloadToken
}
