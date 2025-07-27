package tools

import (
	"context"
	"fmt"

	"github.com/NoASLR/htb-mcp-server/pkg/htb"
	"github.com/NoASLR/htb-mcp-server/pkg/mcp"
)

// GetUserProfile tool for getting user profile information
type GetUserProfile struct {
	client *htb.Client
}

func NewGetUserProfile(client *htb.Client) *GetUserProfile {
	return &GetUserProfile{client: client}
}

func (t *GetUserProfile) Name() string {
	return "get_user_profile"
}

func (t *GetUserProfile) Description() string {
	return "Get the authenticated user's profile information including points, rank, and subscription status"
}

func (t *GetUserProfile) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type:       "object",
		Properties: map[string]mcp.Property{},
	}
}

func (t *GetUserProfile) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	// Make API request to get user info
	data, err := t.client.GetWithParsing(ctx, "/user/info", "info")
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
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

// GetUserProgress tool for getting user progress and statistics
type GetUserProgress struct {
	client *htb.Client
}

func NewGetUserProgress(client *htb.Client) *GetUserProgress {
	return &GetUserProgress{client: client}
}

func (t *GetUserProgress) Name() string {
	return "get_user_progress"
}

func (t *GetUserProgress) Description() string {
	return "Get user progress including completed challenges, machines, and achievements"
}

func (t *GetUserProgress) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"type": {
				Type:        "string",
				Description: "Type of progress to retrieve",
				Enum:        []string{"overview", "machines", "challenges"},
				Default:     "overview",
			},
			"limit": {
				Type:        "integer",
				Description: "Limit the number of results",
				Default:     50,
			},
		},
	}
}

func (t *GetUserProgress) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	progressType := "overview"
	if pt, ok := args["type"].(string); ok {
		progressType = pt
	}

	var endpoint string
	var dataField string

	switch progressType {
	case "machines":
		// Get user's machine progress
		endpoint = "/user/info"
		dataField = "info"
	case "challenges":
		// Get user's challenge progress
		endpoint = "/user/info"
		dataField = "info"
	default:
		// Get general user info
		endpoint = "/user/info"
		dataField = "info"
	}

	// Make API request
	data, err := t.client.GetWithParsing(ctx, endpoint, dataField)
	if err != nil {
		return nil, fmt.Errorf("failed to get user progress: %w", err)
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
