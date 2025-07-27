package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Save original environment
	originalToken := os.Getenv("HTB_TOKEN")
	defer func() {
		if originalToken != "" {
			os.Setenv("HTB_TOKEN", originalToken)
		} else {
			os.Unsetenv("HTB_TOKEN")
		}
	}()

	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
		validate    func(*Config) error
	}{
		{
			name: "valid configuration with defaults",
			envVars: map[string]string{
				"HTB_TOKEN": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodGIiLCJpYXQiOjE2NzA0MjE2MDAsImV4cCI6MTY3MDQyNTIwMCwidXNlcl9pZCI6MTIzNDV9.abc123def456",
			},
			expectError: false,
			validate: func(cfg *Config) error {
				if cfg.HTBBaseURL != "https://labs.hackthebox.com/api/v4" {
					t.Errorf("Expected default HTB base URL, got %s", cfg.HTBBaseURL)
				}
				if cfg.ServerPort != 3000 {
					t.Errorf("Expected default port 3000, got %d", cfg.ServerPort)
				}
				if cfg.LogLevel != "INFO" {
					t.Errorf("Expected default log level INFO, got %s", cfg.LogLevel)
				}
				if cfg.RateLimitPerMinute != 100 {
					t.Errorf("Expected default rate limit 100, got %d", cfg.RateLimitPerMinute)
				}
				if cfg.CacheTTL != 5*time.Minute {
					t.Errorf("Expected default cache TTL 5m, got %v", cfg.CacheTTL)
				}
				return nil
			},
		},
		{
			name: "custom configuration",
			envVars: map[string]string{
				"HTB_TOKEN":               "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodGIiLCJpYXQiOjE2NzA0MjE2MDAsImV4cCI6MTY3MDQyNTIwMCwidXNlcl9pZCI6MTIzNDV9.abc123def456",
				"SERVER_PORT":             "8080",
				"LOG_LEVEL":               "DEBUG",
				"RATE_LIMIT_PER_MINUTE":   "200",
				"CACHE_TTL_SECONDS":       "600",
				"REQUEST_TIMEOUT_SECONDS": "60",
			},
			expectError: false,
			validate: func(cfg *Config) error {
				if cfg.ServerPort != 8080 {
					t.Errorf("Expected port 8080, got %d", cfg.ServerPort)
				}
				if cfg.LogLevel != "DEBUG" {
					t.Errorf("Expected log level DEBUG, got %s", cfg.LogLevel)
				}
				if cfg.RateLimitPerMinute != 200 {
					t.Errorf("Expected rate limit 200, got %d", cfg.RateLimitPerMinute)
				}
				if cfg.CacheTTL != 10*time.Minute {
					t.Errorf("Expected cache TTL 10m, got %v", cfg.CacheTTL)
				}
				if cfg.RequestTimeout != 60*time.Second {
					t.Errorf("Expected request timeout 60s, got %v", cfg.RequestTimeout)
				}
				return nil
			},
		},
		{
			name:        "missing HTB token",
			envVars:     map[string]string{},
			expectError: true,
		},
		{
			name: "invalid HTB token format",
			envVars: map[string]string{
				"HTB_TOKEN": "invalid-token",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Unsetenv("HTB_TOKEN")
			os.Unsetenv("SERVER_PORT")
			os.Unsetenv("LOG_LEVEL")
			os.Unsetenv("RATE_LIMIT_PER_MINUTE")
			os.Unsetenv("CACHE_TTL_SECONDS")
			os.Unsetenv("REQUEST_TIMEOUT_SECONDS")

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			cfg, err := Load()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.validate != nil {
				if err := tt.validate(cfg); err != nil {
					t.Errorf("Validation failed: %v", err)
				}
			}
		})
	}
}

func TestValidateHTBToken(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name:        "valid JWT token",
			token:       "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodGIiLCJpYXQiOjE2NzA0MjE2MDAsImV4cCI6MTY3MDQyNTIwMCwidXNlcl9pZCI6MTIzNDV9.abc123def456",
			expectError: false,
		},
		{
			name:        "empty token",
			token:       "",
			expectError: true,
		},
		{
			name:        "token with no dots",
			token:       "invalidtoken",
			expectError: true,
		},
		{
			name:        "token with one dot",
			token:       "header.payload",
			expectError: true,
		},
		{
			name:        "token with too many dots",
			token:       "header.payload.signature.extra",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateHTBToken(tt.token)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for token %s, but got none", tt.token)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for token %s: %v", tt.token, err)
				}
			}
		})
	}
}

func TestGetHTBAPIURL(t *testing.T) {
	cfg := &Config{
		HTBBaseURL: "https://labs.hackthebox.com/api/v4",
	}

	tests := []struct {
		name     string
		endpoint string
		expected string
	}{
		{
			name:     "simple endpoint",
			endpoint: "/user/info",
			expected: "https://labs.hackthebox.com/api/v4/user/info",
		},
		{
			name:     "endpoint with parameters",
			endpoint: "/machine/paginated/?per_page=20",
			expected: "https://labs.hackthebox.com/api/v4/machine/paginated/?per_page=20",
		},
		{
			name:     "empty endpoint",
			endpoint: "",
			expected: "https://labs.hackthebox.com/api/v4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cfg.GetHTBAPIURL(tt.endpoint)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
