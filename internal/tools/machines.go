package tools

import (
	"context"
	"fmt"

	"github.com/NoASLR/htb-mcp-server/pkg/htb"
	"github.com/NoASLR/htb-mcp-server/pkg/mcp"
)

// ListMachines tool for listing HTB machines
type ListMachines struct {
	client *htb.Client
}

func NewListMachines(client *htb.Client) *ListMachines {
	return &ListMachines{client: client}
}

func (t *ListMachines) Name() string {
	return "list_machines"
}

func (t *ListMachines) Description() string {
	return "Get a list of HackTheBox machines with optional filtering by status, difficulty, and OS"
}

func (t *ListMachines) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"status": {
				Type:        "string",
				Description: "Filter by machine status",
				Enum:        []string{"active", "retired"},
				Default:     "active",
			},
			"difficulty": {
				Type:        "string",
				Description: "Filter by difficulty level",
				Enum:        []string{"Easy", "Medium", "Hard", "Insane"},
			},
			"os": {
				Type:        "string",
				Description: "Filter by operating system",
				Enum:        []string{"Linux", "Windows"},
			},
			"page": {
				Type:        "integer",
				Description: "Page number for pagination",
				Default:     1,
			},
			"per_page": {
				Type:        "integer",
				Description: "Number of machines per page",
				Default:     20,
			},
		},
	}
}

func (t *ListMachines) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	// Extract parameters
	status := "active"
	if s, ok := args["status"].(string); ok {
		status = s
	}

	perPage := 20
	if pp, ok := args["per_page"].(float64); ok {
		perPage = int(pp)
	}

	// Build endpoint URL based on status
	var endpoint string
	if status == "retired" {
		endpoint = fmt.Sprintf("/machine/list/retired/paginated/?per_page=%d&sort_by=release-date", perPage)
	} else {
		endpoint = fmt.Sprintf("/machine/paginated/?per_page=%d", perPage)
	}

	// Make API request
	data, err := t.client.GetWithParsing(ctx, endpoint, "data")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch machines: %w", err)
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

// StartMachine tool for starting a HTB machine
type StartMachine struct {
	client *htb.Client
}

func NewStartMachine(client *htb.Client) *StartMachine {
	return &StartMachine{client: client}
}

func (t *StartMachine) Name() string {
	return "start_machine"
}

func (t *StartMachine) Description() string {
	return "Start a HackTheBox machine by ID and get connection details"
}

func (t *StartMachine) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"machine_id": {
				Type:        "integer",
				Description: "The ID of the machine to start",
			},
		},
		Required: []string{"machine_id"},
	}
}

func (t *StartMachine) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	machineID, ok := args["machine_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("machine_id is required")
	}

	// Build request payload
	payload := htb.MachineActionRequest{
		MachineID: int(machineID),
	}

	// Determine the correct endpoint based on machine type
	// For now, we'll use the standard machine endpoint
	endpoint := fmt.Sprintf("/machine/play/%d", int(machineID))

	// Make API request
	data, err := t.client.PostWithParsing(ctx, endpoint, payload, "")
	if err != nil {
		return nil, fmt.Errorf("failed to start machine: %w", err)
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

// GetMachineIP tool for getting machine IP address
type GetMachineIP struct {
	client *htb.Client
}

func NewGetMachineIP(client *htb.Client) *GetMachineIP {
	return &GetMachineIP{client: client}
}

func (t *GetMachineIP) Name() string {
	return "get_machine_ip"
}

func (t *GetMachineIP) Description() string {
	return "Get the IP address of the currently active machine"
}

func (t *GetMachineIP) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"machine_id": {
				Type:        "integer",
				Description: "Optional machine ID. If not provided, gets the active machine IP",
			},
		},
	}
}

func (t *GetMachineIP) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	// Get active machine information
	data, err := t.client.GetWithParsing(ctx, "/machine/active", "info")
	if err != nil {
		return nil, fmt.Errorf("failed to get active machine: %w", err)
	}

	if data == nil {
		content := mcp.CreateTextContent("No machine is currently active")
		return &mcp.CallToolResponse{
			Content: []mcp.Content{content},
		}, nil
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

// SubmitUserFlag tool for submitting user flags
type SubmitUserFlag struct {
	client *htb.Client
}

func NewSubmitUserFlag(client *htb.Client) *SubmitUserFlag {
	return &SubmitUserFlag{client: client}
}

func (t *SubmitUserFlag) Name() string {
	return "submit_user_flag"
}

func (t *SubmitUserFlag) Description() string {
	return "Submit a user flag for a HackTheBox machine"
}

func (t *SubmitUserFlag) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"machine_id": {
				Type:        "integer",
				Description: "The ID of the machine",
			},
			"flag": {
				Type:        "string",
				Description: "The user flag to submit",
			},
		},
		Required: []string{"machine_id", "flag"},
	}
}

func (t *SubmitUserFlag) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	machineID, ok := args["machine_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("machine_id is required")
	}

	flag, ok := args["flag"].(string)
	if !ok {
		return nil, fmt.Errorf("flag is required")
	}

	// Build request payload
	payload := htb.FlagSubmissionRequest{
		ID:   int(machineID),
		Flag: flag,
	}

	// Make API request
	data, err := t.client.PostWithParsing(ctx, "/machine/own", payload, "message")
	if err != nil {
		return nil, fmt.Errorf("failed to submit user flag: %w", err)
	}

	// Create text content with result
	message := fmt.Sprintf("User flag submission result: %v", data)
	content := mcp.CreateTextContent(message)

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}

// SubmitRootFlag tool for submitting root flags
type SubmitRootFlag struct {
	client *htb.Client
}

func NewSubmitRootFlag(client *htb.Client) *SubmitRootFlag {
	return &SubmitRootFlag{client: client}
}

func (t *SubmitRootFlag) Name() string {
	return "submit_root_flag"
}

func (t *SubmitRootFlag) Description() string {
	return "Submit a root flag for a HackTheBox machine"
}

func (t *SubmitRootFlag) Schema() mcp.ToolSchema {
	return mcp.ToolSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"machine_id": {
				Type:        "integer",
				Description: "The ID of the machine",
			},
			"flag": {
				Type:        "string",
				Description: "The root flag to submit",
			},
		},
		Required: []string{"machine_id", "flag"},
	}
}

func (t *SubmitRootFlag) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	machineID, ok := args["machine_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("machine_id is required")
	}

	flag, ok := args["flag"].(string)
	if !ok {
		return nil, fmt.Errorf("flag is required")
	}

	// Build request payload
	payload := htb.FlagSubmissionRequest{
		ID:   int(machineID),
		Flag: flag,
	}

	// Make API request to the same endpoint (HTB API handles flag type detection)
	data, err := t.client.PostWithParsing(ctx, "/machine/own", payload, "message")
	if err != nil {
		return nil, fmt.Errorf("failed to submit root flag: %w", err)
	}

	// Create text content with result
	message := fmt.Sprintf("Root flag submission result: %v", data)
	content := mcp.CreateTextContent(message)

	return &mcp.CallToolResponse{
		Content: []mcp.Content{content},
	}, nil
}
