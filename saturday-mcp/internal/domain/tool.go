// Package domain defines the trinity of MCP domain concepts — Tool, Persona,
// Workflow — as pure interfaces with zero SDK-provider coupling beyond the
// mcp-go value types used on their signatures.
//
// These interfaces are the seam between the MCP protocol layer (registration,
// dispatch) and the concrete implementations under internal/tools/,
// internal/workflows/, and internal/prompts/. They exist so that adding a
// new tool means dropping a new file into internal/tools/ and appending it to
// a provider slice — never touching the server package.
package domain

import (
	"context"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
)

// Tool is a single MCP tool the client can call. One Go type per registered
// tool — Extract Class was the primary Fowler operation used to arrive at
// this shape (see mcp-add-plan Phase C).
//
// Contract:
//   - Name and Description are stable identifiers; they participate in the
//     MCP tool discovery response and must not change without a version bump.
//   - InputSchema is the JSON Schema the SDK advertises to clients; it is the
//     public contract for what arguments Execute accepts.
//   - Execute is the tool's business behavior. It receives the raw
//     mcp.CallToolRequest so that argument parsing stays close to the tool
//     that owns the schema, rather than being centralized in the server layer.
//
// OutputSchema is generated from the typed response struct each tool
// serializes, using github.com/invopop/jsonschema (see internal/tools/
// responses.go). RegisterTools marshals the returned schema onto
// mcp.Tool.RawOutputSchema so it participates in MCP tool discovery.
// Phase F op 14 added the interface method; the mcp-go v0.56.0 upgrade
// wired it into the registration path.
type Tool interface {
	Name() string
	Description() string
	InputSchema() mcp.ToolInputSchema
	OutputSchema() *jsonschema.Schema
	Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

// ToolProvider yields the full set of tools a server should register. The
// server layer iterates this slice in RegisterTools; it never hardcodes the
// list. Handler composition is the only place tools are constructed.
type ToolProvider interface {
	Tools() []Tool
}
