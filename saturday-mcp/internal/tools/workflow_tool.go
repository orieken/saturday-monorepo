package tools

import (
	"context"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/domain"
)

// WorkflowTool adapts any domain.Workflow into a domain.Tool so the MCP
// server can register workflows through the same code path as single-step
// tools. See mcp-add-plan Phase D op 10.
//
// The adapter is intentionally a thin, non-generic struct: MCP tool
// arguments are already the untyped map[string]interface{} that both Tool
// and Workflow accept, so type parameters would add ceremony without
// buying safety. Each Workflow owns its own request-parsing, mirroring
// how each Tool owns its own schema declaration.
type WorkflowTool struct {
	workflow domain.Workflow
}

// NewWorkflowTool wraps a Workflow so it satisfies domain.Tool.
func NewWorkflowTool(workflow domain.Workflow) *WorkflowTool {
	return &WorkflowTool{workflow: workflow}
}

// Name delegates to the underlying workflow so the MCP tool identifier
// stays a single source of truth on the Workflow itself.
func (t *WorkflowTool) Name() string { return t.workflow.Name() }

// Description delegates to the underlying workflow.
func (t *WorkflowTool) Description() string { return t.workflow.Description() }

// InputSchema delegates to the underlying workflow.
func (t *WorkflowTool) InputSchema() mcp.ToolInputSchema { return t.workflow.InputSchema() }

// OutputSchema delegates to the underlying workflow so the Tool and
// Workflow interfaces stay symmetric under the adapter.
func (t *WorkflowTool) OutputSchema() *jsonschema.Schema { return t.workflow.OutputSchema() }

// Execute forwards the MCP tool call to the workflow's Run method.
func (t *WorkflowTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return t.workflow.Run(ctx, request)
}
