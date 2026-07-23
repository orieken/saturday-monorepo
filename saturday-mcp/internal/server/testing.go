package server

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

// The HandleFoo wrappers below exist purely to keep internal/integration/
// e2e_test.go's response-shape regression suite green through the Phase C/D
// refactor. Each one delegates to the tool the MCP name resolves to via
// h.byName. Once Phase J backfills unit tests directly against the
// internal/tools/ and internal/workflows/ types, these wrappers can be
// deleted and e2e_test.go can be rewritten to consume Handler.Tools()
// (or the tool types directly) instead.

// dispatch looks up a tool by its MCP name and delegates Execute. If the
// name is unknown, it returns an MCP tool-result error so the tests see
// a real failure instead of a nil-pointer panic.
func (h *Handler) dispatch(ctx context.Context, name string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, ok := h.byName[name]
	if !ok {
		return mcp.NewToolResultError("unknown tool: " + name), nil
	}
	return tool.Execute(ctx, request)
}

// HandleGenerateSite delegates to the generate_site tool (Phase C op 6).
func (h *Handler) HandleGenerateSite(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "generate_site", request)
}

// HandleGeneratePage delegates to the generate_page tool (Phase C op 7).
func (h *Handler) HandleGeneratePage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "generate_page", request)
}

// HandleGenerateFlow delegates to the generate_flow tool (Phase C op 8).
func (h *Handler) HandleGenerateFlow(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "generate_flow", request)
}

// HandleGenerateSteps delegates to the generate_steps tool (Phase C op 9).
func (h *Handler) HandleGenerateSteps(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "generate_steps", request)
}

// HandleGenerateElement delegates to the generate_element tool (Phase C op 10).
func (h *Handler) HandleGenerateElement(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "generate_element", request)
}

// HandleGenerateService delegates to the generate_service tool (Phase C op 11).
func (h *Handler) HandleGenerateService(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "generate_service", request)
}

// HandleMigrateCode delegates to the migrate_code tool (Phase C op 12).
func (h *Handler) HandleMigrateCode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "migrate_code", request)
}

// HandleAnalyzePerformance delegates to the analyze_performance tool (Phase C op 13).
func (h *Handler) HandleAnalyzePerformance(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "analyze_performance", request)
}

// HandleGenerateDocumentation delegates to the generate_documentation tool (Phase C op 14).
func (h *Handler) HandleGenerateDocumentation(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "generate_documentation", request)
}

// HandleAnalyzeFramework delegates to the analyze_framework tool (Phase C op 15).
func (h *Handler) HandleAnalyzeFramework(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "analyze_framework", request)
}

// HandleValidatePatterns delegates to the validate_patterns tool (Phase C op 16).
func (h *Handler) HandleValidatePatterns(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "validate_patterns", request)
}

// HandleSuggestImprovements delegates to the suggest_improvements tool (Phase C op 17).
func (h *Handler) HandleSuggestImprovements(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "suggest_improvements", request)
}

// HandleAnalyzeImpact delegates to the analyze_impact tool (Phase C op 18).
func (h *Handler) HandleAnalyzeImpact(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "analyze_impact", request)
}

// HandleParseTestFailure delegates to the parse_test_failure tool (Phase C op 19).
func (h *Handler) HandleParseTestFailure(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "parse_test_failure", request)
}

// HandleRunTests delegates to the run_tests workflow-tool (Phase D op 3).
func (h *Handler) HandleRunTests(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "run_tests", request)
}

// HandlePrioritizeTests delegates to the prioritize_tests workflow-tool (Phase D op 3).
func (h *Handler) HandlePrioritizeTests(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "prioritize_tests", request)
}

// HandleAnalyzeComplexity delegates to the analyze_complexity tool (mcp-expand M1 Op 2).
func (h *Handler) HandleAnalyzeComplexity(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "analyze_complexity", request)
}

// HandleCheckAccessibility delegates to the check_accessibility tool (mcp-expand M1 Op 3).
func (h *Handler) HandleCheckAccessibility(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "check_accessibility", request)
}

// HandleCheckUbiquitousLanguage delegates to the check_ubiquitous_language tool (mcp-expand M1 Op 4).
func (h *Handler) HandleCheckUbiquitousLanguage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "check_ubiquitous_language", request)
}

// HandleVerifyDependencies delegates to the verify_dependencies tool (mcp-expand M1 Op 5).
func (h *Handler) HandleVerifyDependencies(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "verify_dependencies", request)
}

// HandleSearchKI delegates to the search_ki tool (mcp-expand M1 Op 6 — LLM-as-retriever).
func (h *Handler) HandleSearchKI(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "search_ki", request)
}

// HandleSearchDocs delegates to the search_docs tool (mcp-expand M1 Op 7 — BM25 backend).
func (h *Handler) HandleSearchDocs(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.dispatch(ctx, "search_docs", request)
}
