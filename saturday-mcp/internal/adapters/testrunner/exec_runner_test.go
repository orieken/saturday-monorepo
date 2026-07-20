package testrunner

import (
	"context"
	"os"
	"testing"

	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/models"
)

func TestExecRunner_Run(t *testing.T) {
	logger := logging.NewLogger(os.Stderr)
	runner := NewExecRunner(logger)

	req := models.TestExecutionRequest{
		ProjectPath: ".",
		Command:     "echo 'Test Passed'",
	}

	result, err := runner.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected success for echo command")
	}

	if result.Output == "" {
		t.Error("Expected output")
	}

	reqFail := models.TestExecutionRequest{
		ProjectPath: ".",
		Command:     "false",
	}

	resultFail, _ := runner.Run(context.Background(), reqFail)
	if resultFail.Success {
		t.Error("Expected failure for false command")
	}
}
