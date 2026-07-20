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

// AnalyzeImpactTool analyzes the ripple effect of modifying a specific file
// against the project's dependency graph. See mcp-add-plan Phase C op 18.
type AnalyzeImpactTool struct {
	logger        *logging.Logger
	graphAnalyzer *analyzers.GraphAnalyzer
}

// NewAnalyzeImpactTool wires the tool with its dependencies.
func NewAnalyzeImpactTool(logger *logging.Logger, graphAnalyzer *analyzers.GraphAnalyzer) *AnalyzeImpactTool {
	return &AnalyzeImpactTool{logger: logger, graphAnalyzer: graphAnalyzer}
}

func (t *AnalyzeImpactTool) Name() string { return "analyze_impact" }

func (t *AnalyzeImpactTool) Description() string {
	return "Analyze the impact of modifying a specific file (dependency graph)"
}

func (t *AnalyzeImpactTool) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"projectPath", "targetFile"},
		Properties: map[string]interface{}{
			"projectPath": projectPathProperty(),
			"targetFile": map[string]interface{}{
				"type":        "string",
				"description": "Relative path of the file to modify",
			},
		},
	}
}

// OutputSchema advertises the ImpactResult response shape.
func (t *AnalyzeImpactTool) OutputSchema() *jsonschema.Schema {
	return reflectSchema(&ImpactResult{})
}

// Execute builds a dependency graph and analyzes the impact of modifying
// targetFile. Response shape preserved verbatim from
// server.Handler.handleAnalyzeImpact.
func (t *AnalyzeImpactTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling analyze_impact request")

	args := request.GetArguments()
	projectPath, _ := args["projectPath"].(string)
	targetFile, _ := args["targetFile"].(string)

	if projectPath == "" || targetFile == "" {
		return mcp.NewToolResultError("projectPath and targetFile are required"), nil
	}

	graph, err := t.graphAnalyzer.Build(projectPath)
	if err != nil {
		t.logger.Error("Graph build failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build dependency graph: %v", err)), nil
	}

	impactedFiles, err := t.graphAnalyzer.AnalyzeImpact(graph, targetFile)
	if err != nil {
		t.logger.Error("Impact analysis failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Analysis failed: %v", err)), nil
	}

	result := map[string]interface{}{
		"target":   targetFile,
		"impacted": impactedFiles,
		"count":    len(impactedFiles),
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		t.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	t.logger.Info("Impact analysis completed successfully", "target", targetFile, "impacted_count", len(impactedFiles))
	return mcp.NewToolResultText(string(resultJSON)), nil
}
