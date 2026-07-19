package server

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

// HandleGenerateSite is an exported wrapper for testing; delegates to the
// extracted tools.GenerateSiteTool (Phase C op 6).
func (h *Handler) HandleGenerateSite(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.generateSiteTool.Execute(ctx, request)
}

// HandleGeneratePage is an exported wrapper for testing; delegates to
// tools.GeneratePageTool (Phase C op 7).
func (h *Handler) HandleGeneratePage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.generatePageTool.Execute(ctx, request)
}

// HandleGenerateFlow is an exported wrapper for testing; delegates to
// tools.GenerateFlowTool (Phase C op 8).
func (h *Handler) HandleGenerateFlow(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.generateFlowTool.Execute(ctx, request)
}

// HandleGenerateSteps is an exported wrapper for testing; delegates to
// tools.GenerateStepsTool (Phase C op 9).
func (h *Handler) HandleGenerateSteps(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.generateStepsTool.Execute(ctx, request)
}

// HandleAnalyzeFramework is an exported wrapper for testing
func (h *Handler) HandleAnalyzeFramework(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleAnalyzeFramework(ctx, request)
}

// HandleValidatePatterns is an exported wrapper for testing
func (h *Handler) HandleValidatePatterns(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleValidatePatterns(ctx, request)
}

// HandleGenerateElement is an exported wrapper for testing; delegates to
// tools.GenerateElementTool (Phase C op 10).
func (h *Handler) HandleGenerateElement(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.generateElementTool.Execute(ctx, request)
}

// HandleGenerateService is an exported wrapper for testing; delegates to
// tools.GenerateServiceTool (Phase C op 11).
func (h *Handler) HandleGenerateService(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.generateServiceTool.Execute(ctx, request)
}

// HandleSuggestImprovements is an exported wrapper for testing
func (h *Handler) HandleSuggestImprovements(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleSuggestImprovements(ctx, request)
}

// HandleMigrateCode is an exported wrapper for testing
func (h *Handler) HandleMigrateCode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleMigrateCode(ctx, request)
}

// HandleAnalyzePerformance is an exported wrapper for testing
func (h *Handler) HandleAnalyzePerformance(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleAnalyzePerformance(ctx, request)
}

// HandleGenerateDocumentation is an exported wrapper for testing
func (h *Handler) HandleGenerateDocumentation(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleGenerateDocumentation(ctx, request)
}

// HandleAnalyzeImpact is an exported wrapper for testing
func (h *Handler) HandleAnalyzeImpact(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleAnalyzeImpact(ctx, request)
}

// HandleRunTests is an exported wrapper for testing
func (h *Handler) HandleRunTests(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleRunTests(ctx, request)
}

// HandlePrioritizeTests is an exported wrapper for testing
func (h *Handler) HandlePrioritizeTests(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handlePrioritizeTests(ctx, request)
}

// HandleParseTestFailure is an exported wrapper for testing
func (h *Handler) HandleParseTestFailure(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleParseTestFailure(ctx, request)
}
