package executor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/models"
)

// TestExecutor manages test execution
type TestExecutor struct {
	logger *logging.Logger
}

// NewTestExecutor creates a new test executor
func NewTestExecutor(logger *logging.Logger) *TestExecutor {
	return &TestExecutor{
		logger: logger,
	}
}

// Run executes the tests based on the request
func (e *TestExecutor) Run(req models.TestExecutionRequest) (*models.TestExecutionResult, error) {
	e.logger.Info("Starting test execution", "path", req.ProjectPath, "cmd", req.Command)
	
	start := time.Now()

	// Default command if not provided
	cmdStr := req.Command
	if cmdStr == "" {
		cmdStr = "npx playwright test"
	}

	// Add filter if provided
	if req.Filter != "" {
		cmdStr += fmt.Sprintf(" -g '%s'", req.Filter)
	}

	// Split command string into executable and args
	// This is simple splitting, might handle quotes poorly, but fine for MVP
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	head := parts[0]
	args := parts[1:]

	cmd := exec.Command(head, args...)
	cmd.Dir = req.ProjectPath
	
	// Environment variables
	cmd.Env = os.Environ()
	for k, v := range req.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run command
	err := cmd.Run()
	duration := time.Since(start).String()

	output := stdout.String() + "\n" + stderr.String()
	success := err == nil

	// Basic summary extraction (optional, can be improved with regex/parsing)
	summary := "Tests executed"
	if success {
		summary = "Tests Passed"
	} else {
		summary = "Tests Failed"
	}

	e.logger.Info("Test execution complete", "success", success, "duration", duration)

	return &models.TestExecutionResult{
		Success:  success,
		Output:   output,
		Duration: duration,
		Summary:  summary,
	}, nil
}
