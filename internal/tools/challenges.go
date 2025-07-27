package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/NoASLR/htb-mcp-server/pkg/htb"
	"github.com/NoASLR/htb-mcp-server/pkg/mcp"
)

// ListChallenges tool for listing HTB challenges
type ListChallenges struct {
	client *htb.Client
}

func NewListChallenges(client *htb.Client) *ListChallenges {
	return &ListChallenges{client: client}
}

func (t *ListChallenges) Name() string {
	return "list_challenges"
}

func (t *ListChallenges) Description() string {
	return "Get a paginated list of HackTheBox challenges with optional filtering by category, difficulty, and status"
}

func (t *ListChallenges) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"category": {
				Type:        "string",
				Description: "Filter by challenge category (Web, Pwn, Crypto, Forensics, etc.)",
			},
			"difficulty": {
				Type:        "string",
				Description: "Filter by difficulty level",
				Enum:        []string{"Easy", "Medium", "Hard", "Insane"},
			},
			"status": {
				Type:        "string",
				Description: "Filter by challenge status",
				Enum:        []string{"active", "retired"},
				Default:     "active",
			},
			"page": {
				Type:        "integer",
				Description: "Page number for pagination",
				Default:     1,
			},
			"per_page": {
				Type:        "integer",
				Description: "Number of challenges per page",
				Default:     20,
			},
		},
	}
}

func (t *ListChallenges) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	// Extract parameters
	status := "active"
	if s, ok := args["status"].(string); ok {
		status = s
	}

	// Build endpoint URL based on status
	endpoint := "/challenge/list"
	if status == "retired" {
		endpoint = "/challenge/list/retired"
	}

	// Make API request
	data, err := t.client.GetWithParsing(ctx, endpoint, "challenges")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch challenges: %w", err)
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

// StartChallenge tool for starting a HTB challenge
type StartChallenge struct {
	client *htb.Client
}

func NewStartChallenge(client *htb.Client) *StartChallenge {
	return &StartChallenge{client: client}
}

func (t *StartChallenge) Name() string {
	return "start_challenge"
}

func (t *StartChallenge) Description() string {
	return "Start a HackTheBox challenge by ID to initialize the challenge environment"
}

func (t *StartChallenge) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"challenge_id": {
				Type:        "string",
				Description: "The ID of the challenge to start",
			},
		},
		Required: []string{"challenge_id"},
	}
}

func (t *StartChallenge) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	challengeID, ok := args["challenge_id"].(string)
	if !ok {
		return nil, fmt.Errorf("challenge_id is required")
	}

	// Build endpoint URL
	endpoint := fmt.Sprintf("/challenge/%s/start", challengeID)

	// Make API request
	data, err := t.client.PostWithParsing(ctx, endpoint, nil, "")
	if err != nil {
		return nil, fmt.Errorf("failed to start challenge: %w", err)
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

// SubmitChallengeFlag tool for submitting challenge flags
type SubmitChallengeFlag struct {
	client *htb.Client
}

func NewSubmitChallengeFlag(client *htb.Client) *SubmitChallengeFlag {
	return &SubmitChallengeFlag{client: client}
}

func (t *SubmitChallengeFlag) Name() string {
	return "submit_challenge_flag"
}

func (t *SubmitChallengeFlag) Description() string {
	return "Submit a flag for a HackTheBox challenge"
}

func (t *SubmitChallengeFlag) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"challenge_id": {
				Type:        "string",
				Description: "The ID of the challenge",
			},
			"flag": {
				Type:        "string",
				Description: "The flag to submit",
			},
			"difficulty": {
				Type:        "integer",
				Description: "Difficulty rating (1-10)",
			},
		},
		Required: []string{"challenge_id", "flag", "difficulty"},
	}
}

func (t *SubmitChallengeFlag) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	challengeID, ok := args["challenge_id"].(string)
	if !ok {
		return nil, fmt.Errorf("challenge_id is required")
	}

	flag, ok := args["flag"].(string)
	if !ok {
		return nil, fmt.Errorf("flag is required")
	}

	difficulty, ok := args["difficulty"].(float64)
	if !ok {
		return nil, fmt.Errorf("difficulty is required")
	}

	// Convert difficulty to string (HTB API expects difficulty * 10)
	difficultyStr := strconv.Itoa(int(difficulty) * 10)

	// Build request payload
	payload := htb.FlagSubmissionRequest{
		ChallengeID: challengeID,
		Flag:        flag,
		Difficulty:  difficultyStr,
	}

	// Make API request
	data, err := t.client.PostWithParsing(ctx, "/challenge/own", payload, "message")
	if err != nil {
		return nil, fmt.Errorf("failed to submit flag: %w", err)
	}

	// Create text content with result
	message := fmt.Sprintf("Flag submission result: %v", data)
	content := mcp.CreateTextContent(message)

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}
