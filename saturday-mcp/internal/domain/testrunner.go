package domain

import (
	"context"

	"github.com/orieken/saturday-mcp/internal/models"
)

// TestRunner is the domain seam for executing a project's test suite.
// The concrete implementation (internal/adapters/testrunner/) wraps
// os/exec and enforces a context-based timeout; use-cases and workflows
// depend on this interface, never on the concrete type. See mcp-add-plan
// Phase E op 11 for the extraction rationale — the previous concrete
// dependency (executor.TestExecutor) coupled workflows.RunTestsWorkflow
// directly to the OS process boundary, violating architecture-guardrails
// #1 (Clean Architecture dependency direction).
//
// Contract:
//   - The context governs cancellation and timeouts. Adapters MUST honor
//     ctx.Done() and abort the underlying process when the deadline fires.
//   - Request/response types are the existing models.TestExecutionRequest
//     and models.TestExecutionResult — kept in models/ to avoid churning
//     downstream serialization and preserve the MCP tool response shape
//     that internal/integration/e2e_test.go asserts on.
type TestRunner interface {
	Run(ctx context.Context, req models.TestExecutionRequest) (*models.TestExecutionResult, error)
}
