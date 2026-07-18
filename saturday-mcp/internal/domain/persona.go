package domain

import "github.com/mark3labs/mcp-go/mcp"

// Persona is a named, argument-shaped prompt template the client can invoke.
// In the mcp-go SDK vocabulary a persona is registered as an mcp.Prompt; in
// this codebase's ubiquitous language a Persona is one of the three
// first-class domain types (Tool, Persona, Workflow) — see
// shared/blueprints/mcp-server.md.
//
// Contract:
//   - Name and Description are stable identifiers used by MCP prompt discovery.
//   - Arguments declares the required and optional inputs the persona expects,
//     following the mcp-go PromptArgument shape.
//   - Render produces the concrete PromptMessage sequence the client sends to
//     the model, given the invocation arguments. Rendering is pure — no I/O,
//     no side effects; that's what keeps personas testable without an MCP
//     transport.
type Persona interface {
	Name() string
	Description() string
	Arguments() []mcp.PromptArgument
	Render(arguments map[string]string) ([]mcp.PromptMessage, error)
}

// PersonaProvider yields the full set of personas a server should register.
// The server layer iterates this slice in RegisterPrompts; it never hardcodes
// the list. Symmetric with ToolProvider and WorkflowProvider.
type PersonaProvider interface {
	Personas() []Persona
}
