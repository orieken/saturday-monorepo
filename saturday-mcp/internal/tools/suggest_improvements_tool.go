package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/analyzers"
	"github.com/orieken/saturday-mcp/internal/logging"
)

// SuggestImprovementsTool suggests code improvements based on Saturday
// framework best practices. See mcp-add-plan Phase C op 17.
type SuggestImprovementsTool struct {
	logger   *logging.Logger
	analyzer *analyzers.ImprovementAnalyzer
}

// NewSuggestImprovementsTool wires the tool with its dependencies.
func NewSuggestImprovementsTool(logger *logging.Logger, analyzer *analyzers.ImprovementAnalyzer) *SuggestImprovementsTool {
	return &SuggestImprovementsTool{logger: logger, analyzer: analyzer}
}

func (t *SuggestImprovementsTool) Name() string { return "suggest_improvements" }

func (t *SuggestImprovementsTool) Description() string {
	return "Suggest code improvements based on Saturday framework best practices"
}

func (t *SuggestImprovementsTool) InputSchema() mcp.ToolInputSchema {
	return projectPathOnlySchema()
}

// Execute runs the improvement analyzer. Response shape preserved verbatim
// from server.Handler.handleSuggestImprovements.
func (t *SuggestImprovementsTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling suggest_improvements request")

	args := request.Params.Arguments
	projectPath, _ := args["projectPath"].(string)

	if projectPath == "" {
		return mcp.NewToolResultError("projectPath is required"), nil
	}

	result, err := t.analyzer.Analyze(projectPath)
	if err != nil {
		t.logger.Error("Improvement analysis failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Analysis failed: %v", err)), nil
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		t.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	t.logger.Info("Improvement suggestions generated successfully", "path", projectPath, "count", len(result.Suggestions))
	return mcp.NewToolResultText(string(resultJSON)), nil
}
