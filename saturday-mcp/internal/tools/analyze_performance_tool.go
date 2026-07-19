package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/analyzers"
	"github.com/orieken/saturday-mcp/internal/logging"
)

// AnalyzePerformanceTool scans a project for performance bottlenecks using
// the concurrent analyzer. See mcp-add-plan Phase C op 13.
type AnalyzePerformanceTool struct {
	logger   *logging.Logger
	analyzer *analyzers.PerformanceAnalyzer
}

// NewAnalyzePerformanceTool wires the tool with its dependencies.
func NewAnalyzePerformanceTool(logger *logging.Logger, analyzer *analyzers.PerformanceAnalyzer) *AnalyzePerformanceTool {
	return &AnalyzePerformanceTool{logger: logger, analyzer: analyzer}
}

func (t *AnalyzePerformanceTool) Name() string { return "analyze_performance" }

func (t *AnalyzePerformanceTool) Description() string {
	return "Analyze code for performance bottlenecks (concurrent)"
}

func (t *AnalyzePerformanceTool) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"projectPath"},
		Properties: map[string]interface{}{
			"projectPath": map[string]interface{}{
				"type":        "string",
				"description": "Absolute path to the project root",
			},
		},
	}
}

// Execute runs the performance analyzer. Response shape preserved verbatim
// from server.Handler.handleAnalyzePerformance.
func (t *AnalyzePerformanceTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling analyze_performance request")

	args := request.Params.Arguments
	projectPath, _ := args["projectPath"].(string)

	if projectPath == "" {
		return mcp.NewToolResultError("projectPath is required"), nil
	}

	result, err := t.analyzer.Analyze(projectPath)
	if err != nil {
		t.logger.Error("Performance analysis failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Analysis failed: %v", err)), nil
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		t.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	t.logger.Info("Performance analysis generated successfully", "path", projectPath, "issues", len(result.Suggestions))
	return mcp.NewToolResultText(string(resultJSON)), nil
}
