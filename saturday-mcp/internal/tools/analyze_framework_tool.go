package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/analyzers"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/models"
)

// AnalyzeFrameworkTool analyzes the existing framework structure and
// patterns of a project. See mcp-add-plan Phase C op 15.
type AnalyzeFrameworkTool struct {
	logger   *logging.Logger
	analyzer *analyzers.FrameworkAnalyzer
}

// NewAnalyzeFrameworkTool wires the tool with its dependencies.
func NewAnalyzeFrameworkTool(logger *logging.Logger, analyzer *analyzers.FrameworkAnalyzer) *AnalyzeFrameworkTool {
	return &AnalyzeFrameworkTool{logger: logger, analyzer: analyzer}
}

func (t *AnalyzeFrameworkTool) Name() string { return "analyze_framework" }

func (t *AnalyzeFrameworkTool) Description() string {
	return "Analyze existing framework structure and patterns"
}

func (t *AnalyzeFrameworkTool) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"projectPath"},
		Properties: map[string]interface{}{
			"projectPath": projectPathProperty(),
			"patterns": map[string]interface{}{
				"type":        "array",
				"description": "Optional list of patterns to look for",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
	}
}

// OutputSchema advertises the models.AnalysisResult response shape.
func (t *AnalyzeFrameworkTool) OutputSchema() *jsonschema.Schema {
	return reflectSchema(&models.AnalysisResult{})
}

// Execute runs the framework analyzer. Response shape preserved verbatim
// from server.Handler.handleAnalyzeFramework.
func (t *AnalyzeFrameworkTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling analyze_framework request")

	args := request.GetArguments()
	projectPath, _ := args["projectPath"].(string)

	if projectPath == "" {
		return mcp.NewToolResultError("projectPath is required"), nil
	}

	result, err := t.analyzer.Analyze(projectPath)
	if err != nil {
		t.logger.Error("Analysis failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Analysis failed: %v", err)), nil
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		t.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	t.logger.Info("Analysis completed successfully", "path", projectPath)
	return mcp.NewToolResultText(string(resultJSON)), nil
}
