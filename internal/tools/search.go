package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/NoASLR/htb-mcp-server/pkg/htb"
	"github.com/NoASLR/htb-mcp-server/pkg/mcp"
)

// SearchContent tool for searching across HTB platform
type SearchContent struct {
	client *htb.Client
}

func NewSearchContent(client *htb.Client) *SearchContent {
	return &SearchContent{client: client}
}

func (t *SearchContent) Name() string {
	return "search_content"
}

func (t *SearchContent) Description() string {
	return "Search across HackTheBox challenges, machines, and users by name or keyword"
}

func (t *SearchContent) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"query": {
				Type:        "string",
				Description: "Search query string",
			},
			"type": {
				Type:        "string",
				Description: "Type of content to search",
				Enum:        []string{"all", "machines", "challenges", "users"},
				Default:     "all",
			},
		},
		Required: []string{"query"},
	}
}

func (t *SearchContent) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	query, ok := args["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query is required")
	}

	searchType := "all"
	if st, ok := args["type"].(string); ok {
		searchType = st
	}

	// Build search endpoint URL
	endpoint := fmt.Sprintf("/search/fetch?query=%s", query)

	// Make API request
	data, err := t.client.GetWithParsing(ctx, endpoint, "")
	if err != nil {
		return nil, fmt.Errorf("failed to search content: %w", err)
	}

	// Filter results based on search type if specified
	if searchType != "all" && data != nil {
		if dataMap, ok := data.(map[string]interface{}); ok {
			filteredData := make(map[string]interface{})
			if searchType == "machines" && dataMap["machines"] != nil {
				filteredData["machines"] = dataMap["machines"]
			} else if searchType == "challenges" && dataMap["challenges"] != nil {
				filteredData["challenges"] = dataMap["challenges"]
			} else if searchType == "users" && dataMap["users"] != nil {
				filteredData["users"] = dataMap["users"]
			}
			data = filteredData
		}
	}

	// Create JSON content
	content, err := mcp.CreateJSONContent(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON content: %w", err)
	}

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}

// GetServerStatus tool for server health and status information
type GetServerStatus struct {
	client    *htb.Client
	startTime time.Time
}

func NewGetServerStatus(client *htb.Client) *GetServerStatus {
	return &GetServerStatus{
		client:    client,
		startTime: time.Now(),
	}
}

func (t *GetServerStatus) Name() string {
	return "get_server_status"
}

func (t *GetServerStatus) Description() string {
	return "Get MCP server health status and HTB API connectivity information"
}

func (t *GetServerStatus) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type:       "object",
		Properties: map[string]mcp.Property{},
	}
}

func (t *GetServerStatus) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	// Check HTB API health
	htbStatus := "healthy"
	if err := t.client.HealthCheck(ctx); err != nil {
		htbStatus = fmt.Sprintf("unhealthy: %v", err)
	}

	// Calculate uptime
	uptime := time.Since(t.startTime)

	// Build status response
	status := htb.ServerStatus{
		Status:       "running",
		Version:      "1.0.0",
		HTBAPIStatus: htbStatus,
		Uptime:       uptime.String(),
		Timestamp:    time.Now(),
	}

	// Create JSON content
	content, err := mcp.CreateJSONContent(status)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON content: %w", err)
	}

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}
