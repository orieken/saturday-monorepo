package main

import (
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/orieken/saturday-mcp/internal/logging"
	saturdayServer "github.com/orieken/saturday-mcp/internal/server"
)

const (
	serverName    = "saturday-mcp"
	serverVersion = "0.1.0"
)

func main() {
	// Initialize logger
	logger := logging.NewLogger(os.Stderr)
	logger.Info("Starting Saturday MCP Server", "version", serverVersion)

	// Create MCP server
	mcpServer := server.NewMCPServer(
		serverName,
		serverVersion,
		server.WithLogging(),
	)

	// Initialize Saturday server components
	saturdayHandler, err := saturdayServer.NewHandler(logger)
	if err != nil {
		logger.Error("Failed to initialize server handler", "error", err)
		os.Exit(1)
	}

	// Register tools
	if err := saturdayHandler.RegisterTools(mcpServer); err != nil {
		logger.Error("Failed to register tools", "error", err)
		os.Exit(1)
	}

	// Register resources
	if err := saturdayHandler.RegisterResources(mcpServer); err != nil {
		logger.Error("Failed to register resources", "error", err)
		os.Exit(1)
	}

	// Register prompts
	if err := saturdayHandler.RegisterPrompts(mcpServer); err != nil {
		logger.Error("Failed to register prompts", "error", err)
		os.Exit(1)
	}

	logger.Info("Saturday MCP Server initialized successfully")

	// Start stdio transport
	if err := server.ServeStdio(mcpServer); err != nil {
		logger.Error("Server error", "error", err)
		os.Exit(1)
	}
}
