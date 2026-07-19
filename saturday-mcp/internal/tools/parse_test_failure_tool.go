package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/analyzers"
	"github.com/orieken/saturday-mcp/internal/logging"
)

// ParseTestFailureTool parses raw Playwright stdout to identify failing
// files and lines. See mcp-add-plan Phase C op 19.
type ParseTestFailureTool struct {
	logger      *logging.Logger
	logAnalyzer *analyzers.TestLogAnalyzer
}

// NewParseTestFailureTool wires the tool with its dependencies.
func NewParseTestFailureTool(logger *logging.Logger, logAnalyzer *analyzers.TestLogAnalyzer) *ParseTestFailureTool {
	return &ParseTestFailureTool{logger: logger, logAnalyzer: logAnalyzer}
}

func (t *ParseTestFailureTool) Name() string { return "parse_test_failure" }

func (t *ParseTestFailureTool) Description() string {
	return "Parse standard output from Playwright tests to identify failing files and lines"
}

func (t *ParseTestFailureTool) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"output"},
		Properties: map[string]interface{}{
			"output": map[string]interface{}{
				"type":        "string",
				"description": "Raw output string from run_tests",
			},
		},
	}
}

// OutputSchema advertises the []analyzers.TestFailureInfo response shape.
func (t *ParseTestFailureTool) OutputSchema() *jsonschema.Schema {
	return reflectSchema(&[]analyzers.TestFailureInfo{})
}

// Execute runs the log analyzer. Response shape preserved verbatim from
// server.Handler.handleParseTestFailure.
func (t *ParseTestFailureTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling parse_test_failure request")

	output, ok := request.Params.Arguments["output"].(string)
	if !ok || output == "" {
		return mcp.NewToolResultError("output argument is required"), nil
	}

	failures, err := t.logAnalyzer.Parse(output)
	if err != nil {
		t.logger.Error("Failed to parse output", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Parse failed: %v", err)), nil
	}

	resultJSON, err := json.Marshal(failures)
	if err != nil {
		t.logger.Error("Failed to marshal failures", "error", err)
		return mcp.NewToolResultError("Failed to format failures"), nil
	}

	t.logger.Info("Parsed test failures", "count", len(failures))
	return mcp.NewToolResultText(string(resultJSON)), nil
}
