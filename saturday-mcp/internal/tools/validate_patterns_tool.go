package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/analyzers"
	"github.com/orieken/saturday-mcp/internal/logging"
)

// ValidatePatternsTool validates a project's code against Saturday framework
// patterns. See mcp-add-plan Phase C op 16.
type ValidatePatternsTool struct {
	logger    *logging.Logger
	validator *analyzers.PatternValidator
}

// NewValidatePatternsTool wires the tool with its dependencies.
func NewValidatePatternsTool(logger *logging.Logger, validator *analyzers.PatternValidator) *ValidatePatternsTool {
	return &ValidatePatternsTool{logger: logger, validator: validator}
}

func (t *ValidatePatternsTool) Name() string { return "validate_patterns" }

func (t *ValidatePatternsTool) Description() string {
	return "Validate code against Saturday framework patterns"
}

func (t *ValidatePatternsTool) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"projectPath"},
		Properties: map[string]interface{}{
			"projectPath": projectPathProperty(),
			"checkTypes": map[string]interface{}{
				"type":        "array",
				"description": "Optional kinds of checks to run (naming, inheritance, structure)",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
	}
}

// Execute runs the pattern validator. Response shape preserved verbatim
// from server.Handler.handleValidatePatterns.
func (t *ValidatePatternsTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling validate_patterns request")

	args := request.Params.Arguments
	projectPath, _ := args["projectPath"].(string)

	if projectPath == "" {
		return mcp.NewToolResultError("projectPath is required"), nil
	}

	result, err := t.validator.Validate(projectPath)
	if err != nil {
		t.logger.Error("Validation failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Validation failed: %v", err)), nil
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		t.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	t.logger.Info("Validation completed successfully", "path", projectPath, "valid", result.Valid)
	return mcp.NewToolResultText(string(resultJSON)), nil
}
