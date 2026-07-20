package workflows

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/analyzers"
	"github.com/orieken/saturday-mcp/internal/domain/metrics"
	"github.com/orieken/saturday-mcp/internal/logging"
)

// PrioritizeTestsWorkflow orchestrates a metrics-driven ranking of test
// coverage priorities: load PageMetric records via a metrics.Reader,
// run them through UsageAnalyzer, and serialize the result for the MCP
// client. Extracted from server.Handler.handlePrioritizeTests per
// mcp-add-plan Phase D and adapter-isolated in Phase E op 13.
//
// The reader collaborator is the domain metrics.Reader interface, not a
// concrete file adapter — the workflow no longer imports adapters,
// closing the "direct instantiation of infrastructure" gap the plan
// flagged (Pattern Conformance row #2). The file path per invocation
// still travels through the tool input; the Reader interface takes it
// as a per-call arg so one Reader instance serves many requests.
type PrioritizeTestsWorkflow struct {
	logger        *logging.Logger
	usageAnalyzer *analyzers.UsageAnalyzer
	reader        metrics.Reader
}

// NewPrioritizeTestsWorkflow wires the workflow with its collaborators.
func NewPrioritizeTestsWorkflow(logger *logging.Logger, usageAnalyzer *analyzers.UsageAnalyzer, reader metrics.Reader) *PrioritizeTestsWorkflow {
	return &PrioritizeTestsWorkflow{logger: logger, usageAnalyzer: usageAnalyzer, reader: reader}
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

	metricsFile, ok := request.GetArguments()["metricsFile"].(string)
	if !ok || metricsFile == "" {
		return mcp.NewToolResultError("metricsFile argument is required"), nil
	}

	// 1. Load Metrics via the injected domain.Reader — no infrastructure
	// import inside this workflow. Phase E op 13 closed the direct-adapter
	// construction gap the prior code had.
	pageMetrics, err := w.reader.ReadPageMetrics(metricsFile)
	if err != nil {
		w.logger.Error("Failed to read metrics", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to load metrics: %v", err)), nil
	}

	// 2. Analyze
	priorities := w.usageAnalyzer.Analyze(pageMetrics)

	// 3. (Optional) Match with Codebase
	// projectPath, hasProject := request.GetArguments()["projectPath"].(string)
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
