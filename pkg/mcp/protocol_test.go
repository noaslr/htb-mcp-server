package mcp

import (
	"encoding/json"
	"testing"
)

func TestNewRequest(t *testing.T) {
	req := NewRequest(1, MethodInitialize, map[string]string{"test": "value"})

	if req.JSONRPCVersion != "2.0" {
		t.Errorf("Expected JSON-RPC version 2.0, got %s", req.JSONRPCVersion)
	}

	if req.ID != 1 {
		t.Errorf("Expected ID 1, got %v", req.ID)
	}

	if req.Method != MethodInitialize {
		t.Errorf("Expected method %s, got %s", MethodInitialize, req.Method)
	}

	if req.Params == nil {
		t.Errorf("Expected params to be set")
	}
}

func TestNewResponse(t *testing.T) {
	result := map[string]string{"status": "ok"}
	resp := NewResponse(1, result)

	if resp.JSONRPCVersion != "2.0" {
		t.Errorf("Expected JSON-RPC version 2.0, got %s", resp.JSONRPCVersion)
	}

	if resp.ID != 1 {
		t.Errorf("Expected ID 1, got %v", resp.ID)
	}

	if resp.Result == nil {
		t.Errorf("Expected result to be set")
	}

	if resp.Error != nil {
		t.Errorf("Expected error to be nil, got %v", resp.Error)
	}
}

func TestNewErrorResponse(t *testing.T) {
	resp := NewErrorResponse(1, ErrorCodeInvalidRequest, "Invalid request", "additional data")

	if resp.JSONRPCVersion != "2.0" {
		t.Errorf("Expected JSON-RPC version 2.0, got %s", resp.JSONRPCVersion)
	}

	if resp.ID != 1 {
		t.Errorf("Expected ID 1, got %v", resp.ID)
	}

	if resp.Error == nil {
		t.Errorf("Expected error to be set")
	}

	if resp.Error.Code != ErrorCodeInvalidRequest {
		t.Errorf("Expected error code %d, got %d", ErrorCodeInvalidRequest, resp.Error.Code)
	}

	if resp.Error.Message != "Invalid request" {
		t.Errorf("Expected error message 'Invalid request', got %s", resp.Error.Message)
	}

	if resp.Result != nil {
		t.Errorf("Expected result to be nil, got %v", resp.Result)
	}
}

func TestNewNotification(t *testing.T) {
	notif := NewNotification("test/notification", map[string]string{"data": "value"})

	if notif.JSONRPCVersion != "2.0" {
		t.Errorf("Expected JSON-RPC version 2.0, got %s", notif.JSONRPCVersion)
	}

	if notif.ID != nil {
		t.Errorf("Expected ID to be nil for notification, got %v", notif.ID)
	}

	if notif.Method != "test/notification" {
		t.Errorf("Expected method 'test/notification', got %s", notif.Method)
	}

	if notif.Params == nil {
		t.Errorf("Expected params to be set")
	}
}

func TestCreateTextContent(t *testing.T) {
	text := "Hello, world!"
	content := CreateTextContent(text)

	if content.Type != "text" {
		t.Errorf("Expected type 'text', got %s", content.Type)
	}

	if content.Text != text {
		t.Errorf("Expected text '%s', got %s", text, content.Text)
	}
}

func TestCreateJSONContent(t *testing.T) {
	data := map[string]interface{}{
		"name":   "test",
		"value":  123,
		"active": true,
	}

	content, err := CreateJSONContent(data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if content.Type != "text" {
		t.Errorf("Expected type 'text', got %s", content.Type)
	}

	if content.MimeType != "application/json" {
		t.Errorf("Expected mime type 'application/json', got %s", content.MimeType)
	}

	// Verify the JSON is valid
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(content.Text), &parsed); err != nil {
		t.Errorf("Generated JSON is invalid: %v", err)
	}

	// Verify the data matches
	if parsed["name"] != "test" {
		t.Errorf("Expected name 'test', got %v", parsed["name"])
	}

	if parsed["value"].(float64) != 123 {
		t.Errorf("Expected value 123, got %v", parsed["value"])
	}

	if parsed["active"] != true {
		t.Errorf("Expected active true, got %v", parsed["active"])
	}
}

func TestCreateJSONContentError(t *testing.T) {
	// Test with data that cannot be marshaled to JSON
	data := make(chan int) // channels cannot be marshaled to JSON

	_, err := CreateJSONContent(data)
	if err == nil {
		t.Errorf("Expected error when marshaling invalid data, got nil")
	}
}

func TestToolSchema(t *testing.T) {
	schema := ToolSchema{
		Type: "object",
		Properties: map[string]Property{
			"name": {
				Type:        "string",
				Description: "The name field",
			},
			"count": {
				Type:        "integer",
				Description: "The count field",
				Default:     0,
			},
			"tags": {
				Type:        "array",
				Description: "List of tags",
				Items: &Property{
					Type: "string",
				},
			},
		},
		Required: []string{"name"},
	}

	// Verify schema can be marshaled to JSON
	data, err := json.Marshal(schema)
	if err != nil {
		t.Errorf("Failed to marshal schema: %v", err)
	}

	// Verify schema can be unmarshaled from JSON
	var unmarshaled ToolSchema
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("Failed to unmarshal schema: %v", err)
	}

	// Verify properties are preserved
	if unmarshaled.Type != "object" {
		t.Errorf("Expected type 'object', got %s", unmarshaled.Type)
	}

	if len(unmarshaled.Properties) != 3 {
		t.Errorf("Expected 3 properties, got %d", len(unmarshaled.Properties))
	}

	if len(unmarshaled.Required) != 1 {
		t.Errorf("Expected 1 required field, got %d", len(unmarshaled.Required))
	}

	if unmarshaled.Required[0] != "name" {
		t.Errorf("Expected required field 'name', got %s", unmarshaled.Required[0])
	}
}

func TestMessageSerialization(t *testing.T) {
	// Test request serialization
	req := NewRequest(123, MethodListTools, nil)
	data, err := json.Marshal(req)
	if err != nil {
		t.Errorf("Failed to marshal request: %v", err)
	}

	var unmarshaled Message
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("Failed to unmarshal request: %v", err)
	}

	if unmarshaled.JSONRPCVersion != "2.0" {
		t.Errorf("Expected JSON-RPC version 2.0, got %s", unmarshaled.JSONRPCVersion)
	}

	if unmarshaled.Method != MethodListTools {
		t.Errorf("Expected method %s, got %s", MethodListTools, unmarshaled.Method)
	}

	// Test response serialization
	resp := NewResponse(123, map[string]string{"status": "ok"})
	data, err = json.Marshal(resp)
	if err != nil {
		t.Errorf("Failed to marshal response: %v", err)
	}

	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if unmarshaled.Result == nil {
		t.Errorf("Expected result to be set")
	}
}

func TestConstants(t *testing.T) {
	// Verify MCP version is set
	if MCPVersion == "" {
		t.Errorf("MCPVersion should not be empty")
	}

	// Verify error codes are negative (JSON-RPC standard)
	errorCodes := []int{
		ErrorCodeParseError,
		ErrorCodeInvalidRequest,
		ErrorCodeMethodNotFound,
		ErrorCodeInvalidParams,
		ErrorCodeInternalError,
	}

	for _, code := range errorCodes {
		if code >= 0 {
			t.Errorf("Error code %d should be negative", code)
		}
	}

	// Verify methods are properly defined
	methods := []string{
		MethodInitialize,
		MethodListTools,
		MethodCallTool,
		MethodListResources,
		MethodReadResource,
		MethodListPrompts,
		MethodGetPrompt,
	}

	for _, method := range methods {
		if method == "" {
			t.Errorf("Method should not be empty")
		}
	}
}
