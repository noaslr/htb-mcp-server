package htb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/NoASLR/htb-mcp-server/pkg/config"
)

// Client represents an HTB API client
type Client struct {
	httpClient *http.Client
	config     *config.Config
	baseURL    string
}

// NewClient creates a new HTB API client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: cfg.RequestTimeout,
		},
		config:  cfg,
		baseURL: cfg.HTBBaseURL,
	}
}

// Request makes an authenticated HTTP request to the HTB API
func (c *Client) Request(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.config.GetHTBAPIURL(endpoint)
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	req.Header.Set("User-Agent", "htb-mcp-server/1.0")
	req.Header.Set("Authorization", "Bearer "+c.config.HTBToken)

	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json, text/plain, */*")
	} else if method == http.MethodGet {
		req.Header.Set("Host", "labs.hackthebox.com")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	// Check for authentication errors
	if resp.StatusCode == 302 && resp.Header.Get("Location") != "" {
		resp.Body.Close()
		return nil, fmt.Errorf("HTB token appears invalid or expired")
	}

	if resp.StatusCode == 401 {
		resp.Body.Close()
		return nil, fmt.Errorf("unauthorized: HTB token is invalid")
	}

	return resp, nil
}

// Get makes a GET request to the HTB API
func (c *Client) Get(ctx context.Context, endpoint string) (*http.Response, error) {
	return c.Request(ctx, http.MethodGet, endpoint, nil)
}

// Post makes a POST request to the HTB API
func (c *Client) Post(ctx context.Context, endpoint string, body interface{}) (*http.Response, error) {
	return c.Request(ctx, http.MethodPost, endpoint, body)
}

// ParseResponse parses a JSON response and extracts a specific field
func (c *Client) ParseResponse(resp *http.Response, field string) (interface{}, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if field == "" {
		return result, nil
	}

	return result[field], nil
}

// GetWithParsing performs a GET request and parses the response
func (c *Client) GetWithParsing(ctx context.Context, endpoint, field string) (interface{}, error) {
	resp, err := c.Get(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	return c.ParseResponse(resp, field)
}

// PostWithParsing performs a POST request and parses the response
func (c *Client) PostWithParsing(ctx context.Context, endpoint string, body interface{}, field string) (interface{}, error) {
	resp, err := c.Post(ctx, endpoint, body)
	if err != nil {
		return nil, err
	}

	return c.ParseResponse(resp, field)
}

// HealthCheck verifies the HTB API connection and token validity
func (c *Client) HealthCheck(ctx context.Context) error {
	resp, err := c.Get(ctx, "/user/info")
	if err != nil {
		return fmt.Errorf("HTB API health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTB API health check failed with status: %d", resp.StatusCode)
	}

	return nil
}
