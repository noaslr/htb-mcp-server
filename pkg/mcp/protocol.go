package mcp

import (
	"encoding/json"
	"fmt"
)

// Protocol version
const MCPVersion = "2024-11-05"

// Message types
const (
	MessageTypeRequest      = "request"
	MessageTypeResponse     = "response"
	MessageTypeNotification = "notification"
)

// Request methods
const (
	MethodInitialize    = "initialize"
	MethodListTools     = "tools/list"
	MethodCallTool      = "tools/call"
	MethodListResources = "resources/list"
	MethodReadResource  = "resources/read"
	MethodListPrompts   = "prompts/list"
	MethodGetPrompt     = "prompts/get"
)

// Base message structure
type Message struct {
	JSONRPCVersion string      `json:"jsonrpc"`
	ID             interface{} `json:"id,omitempty"`
	Method         string      `json:"method,omitempty"`
	Params         interface{} `json:"params,omitempty"`
	Result         interface{} `json:"result,omitempty"`
	Error          *Error      `json:"error,omitempty"`
}

// Error represents a JSON-RPC error
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Initialize request and response
type InitializeRequest struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    ClientCapabilities     `json:"capabilities"`
	ClientInfo      ClientInfo             `json:"clientInfo"`
	Meta            map[string]interface{} `json:"meta,omitempty"`
}

type InitializeResponse struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
}

type ClientCapabilities struct {
	Experimental map[string]interface{} `json:"experimental,omitempty"`
	Sampling     map[string]interface{} `json:"sampling,omitempty"`
}

type ServerCapabilities struct {
	Experimental map[string]interface{} `json:"experimental,omitempty"`
	Logging      map[string]interface{} `json:"logging,omitempty"`
	Prompts      *PromptsCapability     `json:"prompts,omitempty"`
	Resources    *ResourcesCapability   `json:"resources,omitempty"`
	Tools        *ToolsCapability       `json:"tools,omitempty"`
}

type PromptsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type ResourcesCapability struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Tool definitions
type Tool struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	InputSchema ToolSchema `json:"inputSchema"`
}

type ToolSchema struct {
	Type                 string              `json:"type"`
	Properties           map[string]Property `json:"properties"`
	Required             []string            `json:"required,omitempty"`
	AdditionalProperties interface{}         `json:"additionalProperties,omitempty"`
}

type Property struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Enum        []string    `json:"enum,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Items       *Property   `json:"items,omitempty"`
}

// Tool call request and response
type CallToolRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

type CallToolResponse struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

// Content types
type Content struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	Data     string `json:"data,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

// Resource definitions
type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

type ReadResourceRequest struct {
	URI string `json:"uri"`
}

type ReadResourceResponse struct {
	Contents []ResourceContent `json:"contents"`
}

type ResourceContent struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType"`
	Text     string `json:"text,omitempty"`
	Blob     string `json:"blob,omitempty"`
}

// Helper functions
func NewRequest(id interface{}, method string, params interface{}) *Message {
	return &Message{
		JSONRPCVersion: "2.0",
		ID:             id,
		Method:         method,
		Params:         params,
	}
}

func NewResponse(id interface{}, result interface{}) *Message {
	return &Message{
		JSONRPCVersion: "2.0",
		ID:             id,
		Result:         result,
	}
}

func NewErrorResponse(id interface{}, code int, message string, data interface{}) *Message {
	return &Message{
		JSONRPCVersion: "2.0",
		ID:             id,
		Error: &Error{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}

func NewNotification(method string, params interface{}) *Message {
	return &Message{
		JSONRPCVersion: "2.0",
		Method:         method,
		Params:         params,
	}
}

// Error codes
const (
	ErrorCodeParseError     = -32700
	ErrorCodeInvalidRequest = -32600
	ErrorCodeMethodNotFound = -32601
	ErrorCodeInvalidParams  = -32602
	ErrorCodeInternalError  = -32603
)

// CreateTextContent creates a text content object
func CreateTextContent(text string) Content {
	return Content{
		Type: "text",
		Text: text,
	}
}

// CreateJSONContent creates a JSON content object
func CreateJSONContent(data interface{}) (Content, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return Content{}, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return Content{
		Type:     "text",
		Text:     string(jsonData),
		MimeType: "application/json",
	}, nil
}
