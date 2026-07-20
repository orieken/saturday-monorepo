// Package workflows contains multi-step orchestrations that MCP clients
// invoke as if they were single tools. See mcp-add-plan Phase D — each
// workflow does a parse-input -> invoke-adapter -> analyze -> format
// pipeline that Extract Class surfaced from the old server.Handler.
package workflows

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/domain"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/models"
)

// RunTestsWorkflow orchestrates a full test-execution pass: parse the tool
// request into a TestExecutionRequest, delegate to a domain.TestRunner,
// then serialize the result for the MCP client. Behavior is verbatim from
// server.Handler.handleRunTests so the e2e_test.go regression surface
// continues to pass.
//
// The runner collaborator is the domain.TestRunner interface, not the
// concrete os/exec adapter. Phase E op 11 introduced this seam so the
// workflow no longer imports internal/executor (or its Phase E successor
// internal/adapters/testrunner) — and so unit tests can inject a fake
// runner. Timeout policy is applied inside the runner's context.
type RunTestsWorkflow struct {
	logger *logging.Logger
	runner domain.TestRunner
}

// NewRunTestsWorkflow wires the workflow with its collaborators.
func NewRunTestsWorkflow(logger *logging.Logger, runner domain.TestRunner) *RunTestsWorkflow {
	return &RunTestsWorkflow{logger: logger, runner: runner}
}

// Name is the MCP tool identifier the WorkflowTool adapter advertises.
func (w *RunTestsWorkflow) Name() string { return "run_tests" }

// Description is the human-readable label sent to MCP clients.
func (w *RunTestsWorkflow) Description() string {
	return "Execute tests and capture output"
}

// InputSchema is the JSON Schema advertised via MCP tool discovery.
func (w *RunTestsWorkflow) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"projectPath"},
		Properties: map[string]interface{}{
			"projectPath": map[string]interface{}{
				"type":        "string",
				"description": "Absolute path to the project root",
			},
			"command": map[string]interface{}{
				"type":        "string",
				"description": "Test command (default: npx playwright test)",
			},
			"filter": map[string]interface{}{
				"type":        "string",
				"description": "Grep filter for tests",
			},
		},
	}
}

// OutputSchema advertises the models.TestExecutionResult response shape.
func (w *RunTestsWorkflow) OutputSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&models.TestExecutionResult{})
}

// Run executes the workflow. Response shape preserved from the old
// server.Handler.handleRunTests.
func (w *RunTestsWorkflow) Run(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	w.logger.Info("Handling run_tests request")

	args := request.GetArguments()
	var req models.TestExecutionRequest
	argsJSON, err := json.Marshal(args)
	if err != nil {
		w.logger.Error("Failed to marshal arguments", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	if err := json.Unmarshal(argsJSON, &req); err != nil {
		w.logger.Error("Failed to unmarshal request", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid request format: %v", err)), nil
	}

	if req.ProjectPath == "" {
		return mcp.NewToolResultError("projectPath is required"), nil
	}

	result, err := w.runner.Run(ctx, req)
	if err != nil {
		w.logger.Error("Test execution failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Test execution failed: %v", err)), nil
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		w.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	w.logger.Info("Tests executed", "success", result.Success, "summary", result.Summary)
	return mcp.NewToolResultText(string(resultJSON)), nil
}
