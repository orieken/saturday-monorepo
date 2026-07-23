package tools

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// TestSearchKITool_Metadata verifies the schema-shape contract the
// registration layer depends on. Cheap first line of defence — if the
// tool name or input schema drifts, this fails before any downstream
// integration test does.
func TestSearchKITool_Metadata(t *testing.T) {
	tool := NewSearchKITool(silentLogger(), nil)
	if got := tool.Name(); got != "search_ki" {
		t.Errorf("Name() = %q, want %q", got, "search_ki")
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
	if tool.OutputSchema() == nil {
		t.Error("OutputSchema() must not be nil")
	}
}

// TestSearchKITool_TagMatchOutranksKeyword pins the ADR-002 scoring
// contract: exact tag match (+2.0) must beat body keyword match (+0.3).
// This is the property that makes the calling LLM's semantic filtering
// downstream easier — obvious tag hits float to the top.
func TestSearchKITool_TagMatchOutranksKeyword(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "keyword-only.md"), `---
name: keyword-only-ki
tags: [unrelated]
domain: other
created: 2026-07-23
---

This body mentions retrieval prominently.
`)
	writeFile(t, filepath.Join(root, "tag-match.md"), `---
name: tag-match-ki
tags: [retrieval]
domain: other
created: 2026-07-23
---

Body without the keyword.
`)

	result := runSearchKI(t, root, map[string]interface{}{
		"query": "retrieval",
		"tags":  []interface{}{"retrieval"},
	})

	if result.TotalHits < 2 {
		t.Fatalf("expected both KIs to match, got %d hits", result.TotalHits)
	}
	if result.Matches[0].Path == "" || !strings.HasSuffix(result.Matches[0].Path, "tag-match.md") {
		t.Errorf("expected tag-match KI first, got %q", result.Matches[0].Path)
	}
}

// TestSearchKITool_DomainMatchScores confirms the domain-match boost
// (+1.5 per ADR-002) reaches results even when the query token itself
// doesn't hit the title or body — a common real query pattern.
func TestSearchKITool_DomainMatchScores(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "domain-hit.md"), `---
name: domain-hit-ki
tags: [foo]
domain: retrieval
created: 2026-07-23
---

Body about something else entirely.
`)

	result := runSearchKI(t, root, map[string]interface{}{
		"query":  "retrieval",
		"domain": "retrieval",
	})

	if result.TotalHits == 0 {
		t.Fatal("expected domain match to surface at least one hit")
	}
}

// TestSearchKITool_MissingQueryErrors pins the "query is required"
// contract — matches the analytical tools' missing-required-arg
// behaviour so downstream error handling is uniform.
func TestSearchKITool_MissingQueryErrors(t *testing.T) {
	tool := NewSearchKITool(silentLogger(), NewKICorpusRetriever(nil))
	result, err := tool.Execute(context.Background(), buildRequest(nil))
	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}
	text := extractText(t, result)
	if !strings.Contains(text, "query is required") {
		t.Errorf("expected missing-query error, got %q", text)
	}
}

// TestSearchKITool_EmptyCorpusReturnsZeroHits confirms that an
// unconfigured corpus surfaces cleanly as zero hits rather than an
// error — the caller distinguishes "no matches" from "no corpus" via
// the diagnostic note in the Query field.
func TestSearchKITool_EmptyCorpusReturnsZeroHits(t *testing.T) {
	tool := NewSearchKITool(silentLogger(), NewKICorpusRetriever(nil))
	result, err := tool.Execute(context.Background(), buildRequest(map[string]interface{}{
		"query": "anything",
	}))
	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}

	var parsed KISearchResult
	if err := json.Unmarshal([]byte(extractText(t, result)), &parsed); err != nil {
		t.Fatalf("failed to parse result: %v", err)
	}
	if parsed.TotalHits != 0 {
		t.Errorf("expected 0 hits from empty corpus, got %d", parsed.TotalHits)
	}
	if !strings.Contains(parsed.Query, "no retriever configured") && !strings.Contains(parsed.Query, "anything") {
		t.Errorf("expected diagnostic query text, got %q", parsed.Query)
	}
}

// TestSearchKITool_MalformedFrontmatterSkipped confirms that a broken
// KI file does not fail the whole retrieval — critical for a large
// corpus where one bad file mustn't take out searchability.
func TestSearchKITool_MalformedFrontmatterSkipped(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "broken.md"), "not a valid frontmatter file\n")
	writeFile(t, filepath.Join(root, "good.md"), `---
name: good-ki
tags: [retrieval]
domain: other
created: 2026-07-23
---

A valid KI mentioning retrieval.
`)

	result := runSearchKI(t, root, map[string]interface{}{
		"query": "retrieval",
	})

	if result.TotalHits < 1 {
		t.Errorf("expected the good KI to surface despite the broken sibling; got %d hits", result.TotalHits)
	}
}

// runSearchKI is the moist-test helper for the tag/domain/malformed
// cases — DRYs the corpus-path wiring so each test body reads as
// "arrange the corpus, ask this question, assert on the answer" and
// keeps the assertion path visible where the framework's testing
// conventions expect it.
func runSearchKI(t *testing.T, corpusRoot string, args map[string]interface{}) KISearchResult {
	t.Helper()
	retriever := NewKICorpusRetriever([]string{corpusRoot})
	tool := NewSearchKITool(silentLogger(), retriever)
	result, err := tool.Execute(context.Background(), buildRequest(args))
	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}
	if result == nil {
		t.Fatal("Execute() returned nil result")
	}
	var parsed KISearchResult
	if err := json.Unmarshal([]byte(extractText(t, result)), &parsed); err != nil {
		t.Fatalf("failed to parse result JSON: %v", err)
	}
	return parsed
}

// unused suppresses the mcp import in case a future test doesn't need
// it directly — keeps the import stable across refactors.
var _ = mcp.NewToolResultText
