package workflows

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/analyzers"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/observability"
)

// PrioritizeTestsWorkflow orchestrates a metrics-driven ranking of test
// coverage priorities: load PageMetric records from a JSON file, run them
// through UsageAnalyzer, and serialize the result for the MCP client.
// Extracted from server.Handler.handlePrioritizeTests per mcp-add-plan
// Phase D.
//
// Milestone 2 concern: Run still calls observability.NewFileMetricsProvider
// inline. That is the "direct instantiation of infrastructure" concern
// architecture-guardrails.md flags — Phase E will introduce a
// domain.MetricsReader interface, inject it via the constructor, and let
// the wiring layer choose the file-backed adapter. Keeping the direct call
// here preserves the exact behavior server.Handler had, so Milestone 1
// stays behavior-neutral.
type PrioritizeTestsWorkflow struct {
	logger        *logging.Logger
	usageAnalyzer *analyzers.UsageAnalyzer
}

// NewPrioritizeTestsWorkflow wires the workflow with its collaborators.
func NewPrioritizeTestsWorkflow(logger *logging.Logger, usageAnalyzer *analyzers.UsageAnalyzer) *PrioritizeTestsWorkflow {
	return &PrioritizeTestsWorkflow{logger: logger, usageAnalyzer: usageAnalyzer}
}

// Name is the MCP tool identifier the WorkflowTool adapter advertises.
func (w *PrioritizeTestsWorkflow) Name() string { return "prioritize_tests" }

// Description is the human-readable label sent to MCP clients.
func (w *PrioritizeTestsWorkflow) Description() string {
	return "Rank test coverage needs based on production usage metrics"
}

// InputSchema is the JSON Schema advertised via MCP tool discovery.
func (w *PrioritizeTestsWorkflow) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"metricsFile"},
		Properties: map[string]interface{}{
			"metricsFile": map[string]interface{}{
				"type":        "string",
				"description": "Absolute path to metrics.json file",
			},
			"projectPath": map[string]interface{}{
				"type":        "string",
				"description": "Optional path to project to match files",
			},
		},
	}
}

// OutputSchema advertises the []analyzers.PagePriority response shape.
func (w *PrioritizeTestsWorkflow) OutputSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&[]analyzers.PagePriority{})
}

// Run executes the workflow. Response shape preserved from the old
// server.Handler.handlePrioritizeTests.
func (w *PrioritizeTestsWorkflow) Run(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	w.logger.Info("Handling prioritize_tests request")

	metricsFile, ok := request.Params.Arguments["metricsFile"].(string)
	if !ok || metricsFile == "" {
		return mcp.NewToolResultError("metricsFile argument is required"), nil
	}

	// 1. Load Metrics
	// Phase E TODO: replace this direct construction with an injected
	// domain.MetricsReader; today it matches the pre-extraction handler
	// exactly so no behavior changes.
	provider := observability.NewFileMetricsProvider(metricsFile)
	metrics, err := provider.GetPageMetrics()
	if err != nil {
		w.logger.Error("Failed to read metrics", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to load metrics: %v", err)), nil
	}

	// 2. Analyze
	priorities := w.usageAnalyzer.Analyze(metrics)

	// 3. (Optional) Match with Codebase
	// projectPath, hasProject := request.Params.Arguments["projectPath"].(string)
	// if hasProject && projectPath != "" {
	//    files, _ := listFiles(projectPath)
	//    priorities = w.usageAnalyzer.MatchWithCodebase(priorities, files)
	// }

	resultJSON, err := json.Marshal(priorities)
	if err != nil {
		w.logger.Error("Failed to marshal priorities", "error", err)
		return mcp.NewToolResultError("Failed to format priorities"), nil
	}

	w.logger.Info("Prioritized tests", "count", len(priorities))
	return mcp.NewToolResultText(string(resultJSON)), nil
}
