package server

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/orieken/saturday-mcp/internal/domain"
)

// withTracing wraps a tool's Execute in a span. The span name is
// namespaced under "mcp.tool." + toolName so trace backends can group
// every framework tool call under a common prefix. Tags applied at start:
//   - tool.name
// Tags applied at end (via EndSpan variadic):
//   - tool.success  — false when the tool returned an error result OR
//                     the closure returned err != nil
//   - tool.error_class — populated only on failure; sourced from the
//                        error type or "tool_result_error" for
//                        CallToolResult.IsError=true responses.
// duration_ms is added by the tracer adapter itself — no need to
// re-derive it here.
//
// This is the ONLY place tool spans get created. Zero span calls exist
// inside individual tool Execute methods; the middleware is what
// enforces architecture-guardrails #8 (Observability Boundaries)
// without polluting the domain-layer tools with tracer references.
func withTracing(tracer domain.Tracer, toolName string, next mcpserver.ToolHandlerFunc) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ctx, end := tracer.StartSpan(ctx, "mcp.tool."+toolName,
			domain.String("tool.name", toolName),
		)

		result, err := next(ctx, request)

		success, errClass := classifyOutcome(result, err)
		end(err,
			domain.Bool("tool.success", success),
			domain.String("tool.error_class", errClass),
		)

		return result, err
	}
}

// classifyOutcome derives the two closure-time attributes the middleware
// records: success, and — on failure — an error_class string that is
// stable enough for backend filtering. The classes are intentionally
// coarse (transport_error vs tool_result_error) because fine-grained
// error typing lives inside individual tools, not the tracing layer.
func classifyOutcome(result *mcp.CallToolResult, err error) (success bool, errClass string) {
	if err != nil {
		return false, fmt.Sprintf("%T", err)
	}
	if result != nil && result.IsError {
		return false, "tool_result_error"
	}
	return true, ""
}
