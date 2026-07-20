package server

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/orieken/saturday-mcp/internal/domain"
)

// This file owns MCP registration. Handler exposes each provider slice
// (Tools, Resources, Prompts) and Register* iterates it once. Adding a
// tool becomes: drop a file in internal/tools/, wire it in NewHandler,
// append it to Tools() — no code change here. See mcp-add-plan Phase I
// op 21.

// Tools returns every domain.Tool the handler exposes to MCP clients.
// The concrete ordering is decided once inside buildTools; this method
// only surfaces the slice for the registration loop and any consumer
// (tests, diagnostics) that wants to walk the full inventory.
func (h *Handler) Tools() []domain.Tool {
	return h.tools
}

// RegisterTools registers every tool returned by h.Tools() with the MCP
// server. Each tool's declarative OutputSchema (produced by invopop/jsonschema
// against the tool's typed response struct) is marshaled onto the MCP
// Tool.RawOutputSchema field so clients can discover the response shape.
func (h *Handler) RegisterTools(s *server.MCPServer) error {
	h.logger.Info("Registering tools")

	for _, t := range h.Tools() {
		tool := t
		mcpTool := mcp.Tool{
			Name:        tool.Name(),
			Description: tool.Description(),
			InputSchema: tool.InputSchema(),
		}
		if outSchema := tool.OutputSchema(); outSchema != nil {
			if raw, err := json.Marshal(outSchema); err == nil {
				mcpTool.RawOutputSchema = raw
			} else {
				h.logger.Warn("Failed to marshal output schema", "tool", tool.Name(), "error", err)
			}
		}
		s.AddTool(mcpTool, tool.Execute)
	}

	h.logger.Info("Tools registered successfully")
	return nil
}

// RegisterResources registers every resource yielded by the resource
// provider. The provider owns its own list; this function only wires it.
func (h *Handler) RegisterResources(s *server.MCPServer) error {
	h.logger.Info("Registering resources")

	for _, r := range h.resourceProvider.List() {
		s.AddResource(r, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			return h.resourceProvider.Read(request.Params.URI)
		})
	}

	h.logger.Info("Resources registered successfully")
	return nil
}

// RegisterPrompts registers every prompt yielded by the prompt provider.
// The provider owns its own list; this function only wires it.
func (h *Handler) RegisterPrompts(s *server.MCPServer) error {
	h.logger.Info("Registering prompts")

	for _, p := range h.promptProvider.List() {
		pName := p.Name
		s.AddPrompt(p, func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			messages, err := h.promptProvider.Get(pName, req.Params.Arguments)
			if err != nil {
				return nil, err
			}
			return &mcp.GetPromptResult{
				Messages: messages,
			}, nil
		})
	}

	h.logger.Info("Prompts registered successfully")
	return nil
}
