package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the HTB MCP Server
type Config struct {
	// HTB API Configuration
	HTBToken   string
	HTBBaseURL string

	// Server Configuration
	ServerPort int
	LogLevel   string

	// Rate Limiting
	RateLimitPerMinute int

	// Caching
	CacheTTL time.Duration

	// Timeouts
	RequestTimeout time.Duration
}

// Load creates a new configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		// Default values
		HTBBaseURL:         "https://labs.hackthebox.com/api/v4",
		ServerPort:         3000,
		LogLevel:           "INFO",
		RateLimitPerMinute: 100,
		CacheTTL:           5 * time.Minute,
		RequestTimeout:     30 * time.Second,
	}

	// Required environment variables
	cfg.HTBToken = os.Getenv("HTB_TOKEN")
	if cfg.HTBToken == "" {
		return nil, fmt.Errorf("HTB_TOKEN environment variable is required")
	}

	// Validate HTB token format (should be JWT with 3 parts)
	if err := validateHTBToken(cfg.HTBToken); err != nil {
		return nil, fmt.Errorf("invalid HTB_TOKEN format: %v", err)
	}

	// Optional environment variables
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.ServerPort = p
		}
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}

	if rateLimit := os.Getenv("RATE_LIMIT_PER_MINUTE"); rateLimit != "" {
		if rl, err := strconv.Atoi(rateLimit); err == nil {
			cfg.RateLimitPerMinute = rl
		}
	}

	if cacheTTL := os.Getenv("CACHE_TTL_SECONDS"); cacheTTL != "" {
		if ttl, err := strconv.Atoi(cacheTTL); err == nil {
			cfg.CacheTTL = time.Duration(ttl) * time.Second
		}
	}

	if timeout := os.Getenv("REQUEST_TIMEOUT_SECONDS"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			cfg.RequestTimeout = time.Duration(t) * time.Second
		}
	}

	return cfg, nil
}

// validateHTBToken checks if the token has the correct JWT format
func validateHTBToken(token string) error {
	// Basic JWT validation - should have 3 parts separated by dots
	parts := 0
	for _, char := range token {
		if char == '.' {
			parts++
		}
	}

	if parts != 2 {
		return fmt.Errorf("HTB token must be a valid JWT with 3 parts separated by dots")
	}

	return nil
}

// GetHTBAPIURL returns the full URL for an HTB API endpoint
func (c *Config) GetHTBAPIURL(endpoint string) string {
	return c.HTBBaseURL + endpoint
}
