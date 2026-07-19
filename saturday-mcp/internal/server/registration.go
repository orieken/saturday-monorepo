package server

import (
	"context"

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
// Extracted workflows appear here as tools.WorkflowTool wrappers so
// registration is uniform. Order defines the tool-list order MCP
// advertises.
func (h *Handler) Tools() []domain.Tool {
	return []domain.Tool{
		h.generateSiteTool,
		h.generatePageTool,
		h.generateFlowTool,
		h.generateStepsTool,
		h.generateElementTool,
		h.generateServiceTool,
		h.migrateCodeTool,
		h.analyzePerformanceTool,
		h.generateDocumentationTool,
		h.analyzeFrameworkTool,
		h.validatePatternsTool,
		h.suggestImprovementsTool,
		h.analyzeImpactTool,
		h.runTestsTool,
		h.parseTestFailureTool,
		h.prioritizeTestsTool,
	}
}

// RegisterTools registers every tool returned by h.Tools() with the MCP
// server. mcp-go v0.8.0's mcp.Tool struct does not accept an output
// schema field, so tool.OutputSchema() is not passed here yet — see
// domain.Tool doc for the forward-looking rationale.
func (h *Handler) RegisterTools(s *server.MCPServer) error {
	h.logger.Info("Registering tools")

	for _, t := range h.Tools() {
		tool := t
		s.AddTool(mcp.Tool{
			Name:        tool.Name(),
			Description: tool.Description(),
			InputSchema: tool.InputSchema(),
		}, tool.Execute)
	}

	h.logger.Info("Tools registered successfully")
	return nil
}

// RegisterResources registers every resource yielded by the resource
// provider. The provider owns its own list; this function only wires it.
func (h *Handler) RegisterResources(s *server.MCPServer) error {
	h.logger.Info("Registering resources")

	for _, r := range h.resourceProvider.List() {
		s.AddResource(r, func(ctx context.Context, request mcp.ReadResourceRequest) ([]interface{}, error) {
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
