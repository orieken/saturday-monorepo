package server

import (
	"bytes"
	"testing"

	"github.com/orieken/saturday-mcp/internal/logging"
)

func TestNewHandler(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.NewLogger(&buf)

	handler, err := NewHandler(logger)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}

	if handler.logger == nil {
		t.Error("Expected handler to have logger, got nil")
	}
}

func TestHandlerRegistration(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.NewLogger(&buf)

	handler, err := NewHandler(logger)
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	// Note: We can't fully test RegisterTools/Resources/Prompts without a real MCPServer
	// These tests verify the handler is created correctly
	// Integration tests would verify the full registration flow

	if handler == nil {
		t.Error("Handler should not be nil")
	}
}
