package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NoASLR/htb-mcp-server/internal/tools"
	"github.com/NoASLR/htb-mcp-server/pkg/config"
	"github.com/NoASLR/htb-mcp-server/pkg/htb"
	"github.com/NoASLR/htb-mcp-server/pkg/mcp"
)

// Server represents the MCP server
type Server struct {
	config       *config.Config
	htbClient    *htb.Client
	toolRegistry *tools.Registry
	startTime    time.Time
	input        io.Reader
	output       io.Writer
}

// New creates a new MCP server instance
func New(cfg *config.Config) *Server {
	htbClient := htb.NewClient(cfg)

	return &Server{
		config:       cfg,
		htbClient:    htbClient,
		toolRegistry: tools.NewRegistry(htbClient),
		startTime:    time.Now(),
		input:        os.Stdin,
		output:       os.Stdout,
	}
}

// Start begins the MCP server operation
func (s *Server) Start(ctx context.Context) error {
	// Verify HTB API connection
	if err := s.htbClient.HealthCheck(ctx); err != nil {
		return fmt.Errorf("HTB API health check failed: %w", err)
	}

	log.Printf("HTB MCP Server starting on stdio transport")
	log.Printf("HTB API connection verified")

	// Start processing messages
	go s.processMessages(ctx)

	return nil
}

// Wait waits for shutdown signals
func (s *Server) Wait() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down HTB MCP Server...")
}

// processMessages handles incoming MCP messages
func (s *Server) processMessages(ctx context.Context) {
	scanner := bufio.NewScanner(s.input)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if err := s.handleMessage(ctx, line); err != nil {
			log.Printf("Error handling message: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from stdin: %v", err)
	}
}

// handleMessage processes a single MCP message
func (s *Server) handleMessage(ctx context.Context, line string) error {
	var msg mcp.Message
	if err := json.Unmarshal([]byte(line), &msg); err != nil {
		s.sendErrorResponse(nil, mcp.ErrorCodeParseError, "Parse error", err.Error())
		return nil
	}

	switch msg.Method {
	case mcp.MethodInitialize:
		return s.handleInitialize(ctx, &msg)
	case mcp.MethodListTools:
		return s.handleListTools(ctx, &msg)
	case mcp.MethodCallTool:
		return s.handleCallTool(ctx, &msg)
	default:
		s.sendErrorResponse(msg.ID, mcp.ErrorCodeMethodNotFound, "Method not found", fmt.Sprintf("Unknown method: %s", msg.Method))
		return nil
	}
}

// handleInitialize handles the initialize request
func (s *Server) handleInitialize(ctx context.Context, msg *mcp.Message) error {
	var req mcp.InitializeRequest
	if err := s.parseParams(msg.Params, &req); err != nil {
		s.sendErrorResponse(msg.ID, mcp.ErrorCodeInvalidParams, "Invalid params", err.Error())
		return nil
	}

	// Verify protocol version compatibility
	if req.ProtocolVersion != mcp.MCPVersion {
		log.Printf("Warning: Client protocol version %s differs from server version %s", req.ProtocolVersion, mcp.MCPVersion)
	}

	response := mcp.InitializeResponse{
		ProtocolVersion: mcp.MCPVersion,
		Capabilities: mcp.ServerCapabilities{
			Tools: &mcp.ToolsCapability{
				ListChanged: false,
			},
		},
		ServerInfo: mcp.ServerInfo{
			Name:    "htb-mcp-server",
			Version: "1.0.0",
		},
	}

	return s.sendResponse(msg.ID, response)
}

// handleListTools handles the list tools request
func (s *Server) handleListTools(ctx context.Context, msg *mcp.Message) error {
	tools := s.toolRegistry.GetTools()
	response := map[string]interface{}{
		"tools": tools,
	}

	return s.sendResponse(msg.ID, response)
}

// handleCallTool handles tool call requests
func (s *Server) handleCallTool(ctx context.Context, msg *mcp.Message) error {
	var req mcp.CallToolRequest
	if err := s.parseParams(msg.Params, &req); err != nil {
		s.sendErrorResponse(msg.ID, mcp.ErrorCodeInvalidParams, "Invalid params", err.Error())
		return nil
	}

	// Execute the tool
	result, err := s.toolRegistry.ExecuteTool(ctx, req.Name, req.Arguments)
	if err != nil {
		response := mcp.CallToolResponse{
			Content: []mcp.Content{
				mcp.CreateTextContent(fmt.Sprintf("Error executing tool: %v", err)),
			},
			IsError: true,
		}
		return s.sendResponse(msg.ID, response)
	}

	return s.sendResponse(msg.ID, result)
}

// sendResponse sends a successful response
func (s *Server) sendResponse(id interface{}, result interface{}) error {
	response := mcp.NewResponse(id, result)
	return s.sendMessage(response)
}

// sendErrorResponse sends an error response
func (s *Server) sendErrorResponse(id interface{}, code int, message, data string) error {
	response := mcp.NewErrorResponse(id, code, message, data)
	return s.sendMessage(response)
}

// sendMessage sends a message to the output
func (s *Server) sendMessage(msg *mcp.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if _, err := fmt.Fprintf(s.output, "%s\n", data); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

// parseParams parses message parameters into a struct
func (s *Server) parseParams(params interface{}, target interface{}) error {
	if params == nil {
		return fmt.Errorf("missing parameters")
	}

	data, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %w", err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal params: %w", err)
	}

	return nil
}

// GetUptime returns the server uptime
func (s *Server) GetUptime() time.Duration {
	return time.Since(s.startTime)
}
