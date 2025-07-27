package main

import (
	"context"
	"log"

	"github.com/NoASLR/htb-mcp-server/internal/server"
	"github.com/NoASLR/htb-mcp-server/pkg/config"
)

func main() {
	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create and start the MCP server
	srv := server.New(cfg)

	ctx := context.Background()
	if err := srv.Start(ctx); err != nil {
		log.Fatalf("Failed to start MCP server: %v", err)
	}

	// Wait for shutdown signal
	srv.Wait()
}
