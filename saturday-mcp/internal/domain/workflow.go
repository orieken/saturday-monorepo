package domain

import (
	"context"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
)

// Workflow is a multi-step orchestration exposed to MCP clients as a single
// tool invocation. Workflows differ from Tools in intent, not in transport:
// a Tool is one bounded operation (generate a page, parse a log line), a
// Workflow composes multiple use-cases into a coherent journey (load metrics,
// analyze, format).
//
// In the current codebase, run_tests and prioritize_tests are the two
// handlers that look like Tools but read like Workflows — each does a
// parse-input → invoke-adapter → analyze → format pipeline. Extracting them
// as Workflows (Phase D) surfaces that shape and gives the extraction seam
// for later Milestone 2 concerns (per-step OTel spans, resumability).
//
// Contract:
//   - Name and Description are the identifiers the client sees when the
//     workflow is adapted into a Tool.
//   - InputSchema is the JSON Schema for the workflow's aggregate input; the
//     schema is the workflow's public contract, not the composed steps'.
//   - Run executes the pipeline. It receives the raw request so argument
//     parsing lives beside the schema declaration.
//
// Adaptation: a Workflow becomes an MCP tool via a thin WorkflowTool adapter
// (Phase D op 10) that delegates Tool.Execute → Workflow.Run. This keeps the
// MCP protocol layer ignorant of whether it's driving a single-step Tool or
// a multi-step Workflow.
type Workflow interface {
	Name() string
	Description() string
	InputSchema() mcp.ToolInputSchema
	// OutputSchema mirrors Tool.OutputSchema so the WorkflowTool adapter
	// (internal/tools/workflow_tool.go) can delegate straight through and
	// keep the two interfaces symmetric. Added in mcp-add-plan Phase F.
	OutputSchema() *jsonschema.Schema
	Run(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

// WorkflowProvider yields the full set of workflows a server should register.
// The server layer iterates this slice; it never hardcodes the list.
type WorkflowProvider interface {
	Workflows() []Workflow
}
