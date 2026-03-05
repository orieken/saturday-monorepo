package executor

import (
	"os"
	"testing"

	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/models"
)

func TestTestExecutor_Run(t *testing.T) {
	logger := logging.NewLogger(os.Stderr)
	executor := NewTestExecutor(logger)

	// Test with a simple echo command to simulate success
	req := models.TestExecutionRequest{
		ProjectPath: ".",
		Command:     "echo 'Test Passed'",
	}

	result, err := executor.Run(req)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected success for echo command")
	}

	if result.Output == "" {
		t.Error("Expected output")
	}

	// Test with a failing command
	reqFail := models.TestExecutionRequest{
		ProjectPath: ".",
		Command:     "false", // returns exit code 1
	}

	resultFail, _ := executor.Run(reqFail)
	if resultFail.Success {
		t.Error("Expected failure for false command")
	}
}
