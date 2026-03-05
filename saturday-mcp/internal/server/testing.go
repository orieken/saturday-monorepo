package server

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

// HandleListTools is an exported wrapper for testing
func (h *Handler) HandleListTools(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleListTools(ctx, request)
}

// HandleGenerateSite is an exported wrapper for testing
func (h *Handler) HandleGenerateSite(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleGenerateSite(ctx, request)
}

// HandleGeneratePage is an exported wrapper for testing
func (h *Handler) HandleGeneratePage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleGeneratePage(ctx, request)
}

// HandleGenerateFlow is an exported wrapper for testing
func (h *Handler) HandleGenerateFlow(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleGenerateFlow(ctx, request)
}

// HandleGenerateSteps is an exported wrapper for testing
func (h *Handler) HandleGenerateSteps(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleGenerateSteps(ctx, request)
}

// HandleAnalyzeFramework is an exported wrapper for testing
func (h *Handler) HandleAnalyzeFramework(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleAnalyzeFramework(ctx, request)
}

// HandleValidatePatterns is an exported wrapper for testing
func (h *Handler) HandleValidatePatterns(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleValidatePatterns(ctx, request)
}

// HandleGenerateElement is an exported wrapper for testing
func (h *Handler) HandleGenerateElement(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleGenerateElement(ctx, request)
}

// HandleGenerateService is an exported wrapper for testing
func (h *Handler) HandleGenerateService(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleGenerateService(ctx, request)
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
