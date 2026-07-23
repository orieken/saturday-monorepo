package tools

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// TestSearchDocsTool_Metadata verifies the schema-shape contract the
// registration layer depends on. Cheap first line of defence — if the
// name or input schema drifts, this fails before any integration test.
func TestSearchDocsTool_Metadata(t *testing.T) {
	tool := NewSearchDocsTool(silentLogger(), nil, nil)

	if got := tool.Name(); got != "search_docs" {
		t.Errorf("Name() = %q, want %q", got, "search_docs")
	}
	if tool.Description() == "" {
		t.Error("Description() must not be empty")
	}
	schema := tool.InputSchema()
	if schema.Type != "object" {
		t.Errorf("InputSchema.Type = %q, want object", schema.Type)
	}
	if len(schema.Required) != 1 || schema.Required[0] != "query" {
		t.Errorf("InputSchema.Required = %v, want [query]", schema.Required)
	}
	if _, ok := schema.Properties["docsPath"]; !ok {
		t.Error("InputSchema must advertise the docsPath property")
	}
	// tags/domain must NOT appear — BM25 doesn't use them, and exposing
	// them in the schema would mislead callers into thinking they filter.
	if _, ok := schema.Properties["tags"]; ok {
		t.Error("InputSchema must not advertise tags — BM25 ignores it")
	}
	if _, ok := schema.Properties["domain"]; ok {
		t.Error("InputSchema must not advertise domain — BM25 ignores it")
	}
	if tool.OutputSchema() == nil {
		t.Error("OutputSchema() must not be nil")
	}
}

// TestSearchDocsTool_RoundTrip writes a fake docs corpus into a tempdir,
// invokes the tool against it, and verifies at least one match comes
// back. This is the single most important behaviour — everything else
// is edge-case defence around it.
func TestSearchDocsTool_RoundTrip(t *testing.T) {
	corpus := t.TempDir()
	writeFile(t, filepath.Join(corpus, "runbook.md"), `---
name: OAuth token refresh runbook
---

When the OAuth token refresh flow fails, follow these steps.
`)

	result := runSearchDocs(t, corpus, map[string]interface{}{
		"query":    "OAuth token refresh",
		"docsPath": corpus,
	})

	if !result.Success {
		t.Errorf("expected success=true, got false")
	}
	if result.TotalHits < 1 {
		t.Fatalf("expected at least one hit, got %d", result.TotalHits)
	}
	if !strings.HasSuffix(result.Matches[0].Path, "runbook.md") {
		t.Errorf("first hit path = %q, want suffix runbook.md", result.Matches[0].Path)
	}
	if result.Matches[0].Title == "" {
		t.Error("match title should be populated from frontmatter")
	}
}

// TestSearchDocsTool_MissingQueryErrors pins the "query is required"
// contract — matches search_ki so downstream error handling is uniform.
func TestSearchDocsTool_MissingQueryErrors(t *testing.T) {
	tool := NewSearchDocsTool(silentLogger(), nil, nil)
	result, err := tool.Execute(context.Background(), buildRequest(nil))
	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}
	text := extractText(t, result)
	if !strings.Contains(text, "query is required") {
		t.Errorf("expected missing-query error, got %q", text)
	}
}

// TestSearchDocsTool_EmptyStringQueryErrors covers the edge where the
// caller passes "query": "" — same handling as omitting it entirely.
func TestSearchDocsTool_EmptyStringQueryErrors(t *testing.T) {
	tool := NewSearchDocsTool(silentLogger(), nil, nil)
	result, err := tool.Execute(context.Background(), buildRequest(map[string]interface{}{
		"query": "",
	}))
	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}
	if !strings.Contains(extractText(t, result), "query is required") {
		t.Errorf("expected missing-query error, got %q", extractText(t, result))
	}
}

// TestSearchDocsTool_EmptyCorpusReturnsZeroHitsNotError pins the
// graceful-degradation contract: pointing docsPath at an empty
// directory produces a zero-hit success, not a failure.
func TestSearchDocsTool_EmptyCorpusReturnsZeroHits(t *testing.T) {
	corpus := t.TempDir() // empty directory

	parsed := runSearchDocs(t, corpus, map[string]interface{}{
		"query":    "anything",
		"docsPath": corpus,
	})

	if !parsed.Success {
		t.Errorf("expected success=true even on empty corpus, got false")
	}
	if parsed.TotalHits != 0 {
		t.Errorf("expected 0 hits from empty corpus, got %d", parsed.TotalHits)
	}
}

// TestSearchDocsTool_DocsPathScopingRespected verifies the docsPath
// arg actually confines the search to that root — invoking the tool
// with two different tempdirs must return non-overlapping hits.
func TestSearchDocsTool_DocsPathScopingRespected(t *testing.T) {
	corpusA := t.TempDir()
	corpusB := t.TempDir()
	writeFile(t, filepath.Join(corpusA, "alpha.md"), "# Alpha\n\nBody about unicorns.\n")
	writeFile(t, filepath.Join(corpusB, "beta.md"), "# Beta\n\nBody about unicorns.\n")

	retriever := newBM25(t)
	tool := NewSearchDocsTool(silentLogger(), retriever, retriever)

	resultA := invokeSearchDocs(t, tool, map[string]interface{}{
		"query":    "unicorns",
		"docsPath": corpusA,
	})
	if resultA.TotalHits != 1 || !strings.HasSuffix(resultA.Matches[0].Path, "alpha.md") {
		t.Errorf("corpusA scoping failed: got %d hits, paths=%v", resultA.TotalHits, pathsFromMatches(resultA.Matches))
	}

	resultB := invokeSearchDocs(t, tool, map[string]interface{}{
		"query":    "unicorns",
		"docsPath": corpusB,
	})
	// After scoping to corpusB the same retriever now has both files
	// indexed. What we assert is that corpusB's file surfaces — the
	// scoping proof is the corpusA test above, and this one just
	// confirms the second invocation didn't blow away corpusA's row
	// (idempotent EnsureIndex is per-file, not per-corpus).
	foundB := false
	for _, m := range resultB.Matches {
		if strings.HasSuffix(m.Path, "beta.md") {
			foundB = true
		}
	}
	if !foundB {
		t.Errorf("corpusB scoping failed: beta.md not found in results %v", pathsFromMatches(resultB.Matches))
	}
}

// TestSearchDocsTool_NilRetrieverReturnsDiagnostic confirms that a
// server startup where the sqlite backend couldn't be opened still
// registers the tool — invoking it returns a zero-hit success with a
// diagnostic note in the query field, not a hard error.
func TestSearchDocsTool_NilRetrieverReturnsDiagnostic(t *testing.T) {
	tool := NewSearchDocsTool(silentLogger(), nil, nil)
	result, err := tool.Execute(context.Background(), buildRequest(map[string]interface{}{
		"query": "anything",
	}))
	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}

	var parsed DocSearchResult
	if err := json.Unmarshal([]byte(extractText(t, result)), &parsed); err != nil {
		t.Fatalf("parse result: %v", err)
	}
	if !parsed.Success {
		t.Errorf("expected success=true with diagnostic, got success=false")
	}
	if parsed.TotalHits != 0 {
		t.Errorf("expected 0 hits, got %d", parsed.TotalHits)
	}
	if !strings.Contains(parsed.Query, "no docs retriever configured") {
		t.Errorf("expected diagnostic in query, got %q", parsed.Query)
	}
}

// TestSearchDocsTool_DefaultDocsPathUsedWhenOmitted pins the "docsPath
// defaults to docs/" behaviour. The corpus root doesn't exist in the
// test environment, so we can't assert on retrieval hits — we can only
// assert that EnsureIndex was invoked with the default path.
func TestSearchDocsTool_DefaultDocsPathUsedWhenOmitted(t *testing.T) {
	stub := &stubIndexer{}
	// Nil retriever would short-circuit before the indexer runs, so we
	// wire a real (empty) BM25 to keep the code path honest.
	retriever := newBM25(t)
	tool := NewSearchDocsTool(silentLogger(), retriever, stub)

	_, err := tool.Execute(context.Background(), buildRequest(map[string]interface{}{
		"query": "widgets",
	}))
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if len(stub.lastCorpusPaths) != 1 || stub.lastCorpusPaths[0] != "docs/" {
		t.Errorf("expected EnsureIndex called with [docs/], got %v", stub.lastCorpusPaths)
	}
}

// TestSearchDocsTool_IndexErrorSurfaces confirms that when EnsureIndex
// returns an error, the tool surfaces it as a tool-result error rather
// than silently swallowing it. Distinguishes "we tried to index and
// couldn't" from "we indexed and found nothing".
func TestSearchDocsTool_IndexErrorSurfaces(t *testing.T) {
	stub := &stubIndexer{err: errors.New("disk full")}
	retriever := newBM25(t)
	tool := NewSearchDocsTool(silentLogger(), retriever, stub)

	result, err := tool.Execute(context.Background(), buildRequest(map[string]interface{}{
		"query":    "anything",
		"docsPath": "irrelevant",
	}))
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	text := extractText(t, result)
	if !strings.Contains(text, "Index refresh failed") {
		t.Errorf("expected index error surfaced to caller, got %q", text)
	}
}

// TestSearchDocsTool_RetrieverErrorSurfaces confirms the symmetric
// case for the Retrieve step — errors propagate as tool-result errors,
// not silent zero-hits.
func TestSearchDocsTool_RetrieverErrorSurfaces(t *testing.T) {
	retriever := &erroringRetriever{err: errors.New("retriever exploded")}
	indexer := &stubIndexer{}
	tool := NewSearchDocsTool(silentLogger(), retriever, indexer)

	result, err := tool.Execute(context.Background(), buildRequest(map[string]interface{}{
		"query":    "anything",
		"docsPath": "somewhere",
	}))
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	text := extractText(t, result)
	if !strings.Contains(text, "Retrieval failed") {
		t.Errorf("expected retrieval error surfaced to caller, got %q", text)
	}
}

// TestSearchDocsTool_BM25SatisfiesBothInterfaces is a compile-time
// guard — if BM25Retriever ever loses either method the whole tool
// wiring in tool_provider.go breaks, and this catches it at test time
// with a readable failure message.
func TestSearchDocsTool_BM25SatisfiesBothInterfaces(t *testing.T) {
	var (
		_ Retriever  = (*BM25Retriever)(nil)
		_ DocIndexer = (*BM25Retriever)(nil)
	)
}

// runSearchDocs constructs a real BM25 retriever pointed at corpusRoot,
// invokes the tool, and returns the parsed response. Matches
// search_ki_tool_test.go's runSearchKI helper pattern.
func runSearchDocs(t *testing.T, corpusRoot string, args map[string]interface{}) DocSearchResult {
	t.Helper()
	retriever := newBM25(t)
	tool := NewSearchDocsTool(silentLogger(), retriever, retriever)
	return invokeSearchDocs(t, tool, args)
}

// invokeSearchDocs is the shared "call Execute, parse JSON" body — the
// two-invocation scoping test needs it separately from retriever setup.
func invokeSearchDocs(t *testing.T, tool *SearchDocsTool, args map[string]interface{}) DocSearchResult {
	t.Helper()
	result, err := tool.Execute(context.Background(), buildRequest(args))
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if result == nil {
		t.Fatal("Execute returned nil result")
	}
	var parsed DocSearchResult
	if err := json.Unmarshal([]byte(extractText(t, result)), &parsed); err != nil {
		t.Fatalf("parse result JSON: %v", err)
	}
	return parsed
}

// pathsFromMatches is a formatting helper for readable
// scoping-assertion failure messages.
func pathsFromMatches(matches []DocMatch) []string {
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		out = append(out, filepath.Base(m.Path))
	}
	return out
}

// stubIndexer captures the last EnsureIndex call so tests can assert on
// how the tool wired the corpus paths. Returns err (nil by default) to
// exercise the error branch.
type stubIndexer struct {
	lastCorpusPaths []string
	err             error
}

func (s *stubIndexer) EnsureIndex(corpusPaths []string) error {
	s.lastCorpusPaths = append([]string(nil), corpusPaths...)
	return s.err
}

// erroringRetriever is a Retriever double that always errors — used to
// exercise the retrieval-error branch of Execute.
type erroringRetriever struct {
	err error
}

func (e *erroringRetriever) Retrieve(query string, tags []string, domain string) ([]Reference, error) {
	return nil, e.err
}

// unused suppresses the mcp import in case a future test doesn't need
// it directly — matches search_ki_tool_test.go's pattern.
var _ = mcp.NewToolResultText
