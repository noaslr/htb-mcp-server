package tools

import (
	"context"
	"fmt"

	"github.com/NoASLR/htb-mcp-server/pkg/htb"
	"github.com/NoASLR/htb-mcp-server/pkg/mcp"
)

// Registry manages all available MCP tools
type Registry struct {
	tools     map[string]Tool
	htbClient *htb.Client
}

// Tool interface that all HTB tools must implement
type Tool interface {
	Name() string
	Description() string
	Schema() mcp.ToolSchema
	Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error)
}

// NewRegistry creates a new tool registry
func NewRegistry(htbClient *htb.Client) *Registry {
	registry := &Registry{
		tools:     make(map[string]Tool),
		htbClient: htbClient,
	}

	// Register all available tools
	registry.registerTools()

	return registry
}

// registerTools registers all available HTB tools
func (r *Registry) registerTools() {
	// Challenge management tools
	r.RegisterTool(NewListChallenges(r.htbClient))
	r.RegisterTool(NewStartChallenge(r.htbClient))
	r.RegisterTool(NewSubmitChallengeFlag(r.htbClient))

	// Machine management tools
	r.RegisterTool(NewListMachines(r.htbClient))
	r.RegisterTool(NewStartMachine(r.htbClient))
	r.RegisterTool(NewGetMachineIP(r.htbClient))
	r.RegisterTool(NewSubmitUserFlag(r.htbClient))
	r.RegisterTool(NewSubmitRootFlag(r.htbClient))

	// User management tools
	r.RegisterTool(NewGetUserProfile(r.htbClient))
	r.RegisterTool(NewGetUserProgress(r.htbClient))

	// Search and utility tools
	r.RegisterTool(NewSearchContent(r.htbClient))
	r.RegisterTool(NewGetServerStatus(r.htbClient))
}

// RegisterTool registers a new tool
func (r *Registry) RegisterTool(tool Tool) {
	r.tools[tool.Name()] = tool
}

// GetTool returns a tool by name
func (r *Registry) GetTool(name string) (Tool, bool) {
	tool, exists := r.tools[name]
	return tool, exists
}

// GetTools returns all registered tools in MCP format
func (r *Registry) GetTools() []mcp.Tool {
	var tools []mcp.Tool

	for _, tool := range r.tools {
		tools = append(tools, mcp.Tool{
			Name:        tool.Name(),
			Description: tool.Description(),
			InputSchema: tool.Schema(),
		})
	}

	return tools
}

// ExecuteTool executes a tool by name with the given arguments
func (r *Registry) ExecuteTool(ctx context.Context, name string, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	tool, exists := r.GetTool(name)
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", name)
	}

	return tool.Execute(ctx, args)
}

// ListToolNames returns a list of all registered tool names
func (r *Registry) ListToolNames() []string {
	var names []string
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}
