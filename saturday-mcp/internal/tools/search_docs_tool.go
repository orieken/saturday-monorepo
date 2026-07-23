// Package tools - search_docs_tool.go
//
// SearchDocsTool exposes the BM25Retriever (see bm25_retriever.go) as
// the MCP tool an agent calls to search the installed project's docs
// corpus. Per framework ADR-002 (Corpus-Aware Retrieval Strategy), this
// is retrieval tier 2 — BM25 via sqlite-fts5 over prose docs. Tier 1
// (KI + ADR corpus, LLM-as-retriever) is search_ki; tier 3 (feature
// archive, vector similarity) lands in M2 as search_features.
//
// **/reindex-docs skill spec.** No standalone skill exists yet — the
// tool refreshes its own index on every call by invoking EnsureIndex
// against the caller-supplied docsPath. This is fast for the small
// docs corpora M1 targets (fractions of a second for a few hundred
// files) and buys the "just works" property: the caller never has to
// remember to rebuild the index after editing markdown. If a future
// corpus grows large enough for the per-call refresh to bite, extract
// EnsureIndex out to a `/reindex-docs` skill and drop the per-call
// call here.
package tools

import (
	"context"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/logging"
)

// DocIndexer refreshes a docs index for the given corpus roots. Split
// out from Retriever so tests can inject a fake indexer without needing
// a real sqlite backend, and so a future no-op indexer (once
// /reindex-docs becomes a standalone skill) can be swapped in without
// touching the tool. BM25Retriever satisfies both interfaces, so
// production wiring stays a single object.
type DocIndexer interface {
	EnsureIndex(corpusPaths []string) error
}

// defaultDocsPath is the corpus root used when a caller omits docsPath.
// Matches framework ADR-002's tier-2 corpus definition (docs/features/,
// docs/adrs/, docs/patterns/, docs/runbooks/ — all live under docs/).
const defaultDocsPath = "docs/"

// SearchDocsTool implements domain.Tool for BM25 docs-corpus search.
// It owns neither the retriever nor the indexer — tool_provider.go
// constructs a BM25Retriever and passes it in for both roles.
type SearchDocsTool struct {
	logger    *logging.Logger
	retriever Retriever
	indexer   DocIndexer
}

// NewSearchDocsTool wires the tool with its retriever and indexer.
// Passing nil for either is legal — the tool then returns zero hits
// gracefully with a diagnostic summary in the Query field, mirroring
// search_ki's "no retriever configured" behaviour. This lets
// tool_provider.go register the tool even when the sqlite DB can't be
// opened, so the server never fails startup on retrieval-layer setup.
func NewSearchDocsTool(logger *logging.Logger, retriever Retriever, indexer DocIndexer) *SearchDocsTool {
	return &SearchDocsTool{logger: logger, retriever: retriever, indexer: indexer}
}

func (t *SearchDocsTool) Name() string { return "search_docs" }

func (t *SearchDocsTool) Description() string {
	return "Search the installed project's markdown docs (docs/features/, docs/adrs/, docs/patterns/, docs/runbooks/) via BM25. Returns lexically-ranked references (paths + snippets, never content copies) for the calling LLM to filter semantically."
}

func (t *SearchDocsTool) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"query"},
		Properties: map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "Free-text query. Whitespace-split into tokens; each token is matched via fts5 against doc titles (10x weight) and bodies (1x weight).",
			},
			"docsPath": map[string]interface{}{
				"type":        "string",
				"description": "Corpus root to index and search. Defaults to \"docs/\" (relative to the server's working directory).",
			},
		},
	}
}

func (t *SearchDocsTool) OutputSchema() *jsonschema.Schema {
	return reflectSchema(&DocSearchResult{})
}

func (t *SearchDocsTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling search_docs request")

	args := request.GetArguments()
	query, _ := args["query"].(string)
	if query == "" {
		return mcp.NewToolResultError("query is required"), nil
	}

	docsPath, _ := args["docsPath"].(string)
	if docsPath == "" {
		docsPath = defaultDocsPath
	}

	if t.retriever == nil || t.indexer == nil {
		return t.emptyResult(query, "no docs retriever configured")
	}

	if err := t.indexer.EnsureIndex([]string{docsPath}); err != nil {
		t.logger.Error("Docs index refresh failed", "error", err, "docsPath", docsPath)
		return mcp.NewToolResultError(fmt.Sprintf("Index refresh failed: %v", err)), nil
	}

	// tags and domain are nil — BM25 doesn't consume structured metadata
	// (see bm25_retriever.go's Retrieve docstring). We omit them from
	// the input schema entirely so callers aren't misled.
	refs, err := t.retriever.Retrieve(query, nil, "")
	if err != nil {
		t.logger.Error("Docs retrieval failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Retrieval failed: %v", err)), nil
	}

	result := DocSearchResult{
		Success:   true,
		Query:     query,
		TotalHits: len(refs),
		Matches:   convertReferencesToDocMatches(refs),
	}
	return marshalToolResult(t.logger, result, "search_docs")
}

// emptyResult returns a zero-hit success response with a diagnostic
// note in the query field, matching search_ki's pattern so the caller
// can distinguish "no matches found" from "no corpus configured"
// without an error.
func (t *SearchDocsTool) emptyResult(query, note string) (*mcp.CallToolResult, error) {
	result := DocSearchResult{
		Success:   true,
		Query:     fmt.Sprintf("%s (%s)", query, note),
		TotalHits: 0,
	}
	return marshalToolResult(t.logger, result, "search_docs")
}

// convertReferencesToDocMatches drops the Tags/Domain fields from the
// Retriever's Reference type (which are always zero-value from BM25) so
// the JSON output shape stays honest — no misleadingly empty fields.
func convertReferencesToDocMatches(refs []Reference) []DocMatch {
	matches := make([]DocMatch, 0, len(refs))
	for _, r := range refs {
		matches = append(matches, DocMatch{
			Title:     r.Title,
			Path:      r.Path,
			Summary:   r.Summary,
			Relevance: r.Relevance,
		})
	}
	return matches
}
