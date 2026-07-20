package tools

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/orieken/saturday-mcp/internal/analyzers"
	"github.com/orieken/saturday-mcp/internal/models"
)

// newAnalyzePerformanceTool wires the tool with a real concurrent
// PerformanceAnalyzer. Same moist-tests rationale — we let the
// analyzer actually walk the fixture project so the code the tool
// serializes back is the real analyzer output.
func newAnalyzePerformanceTool(t *testing.T) *AnalyzePerformanceTool {
	t.Helper()
	analyzer := analyzers.NewPerformanceAnalyzer(silentLogger())
	return NewAnalyzePerformanceTool(silentLogger(), analyzer)
}

func TestAnalyzePerformanceTool_Metadata(t *testing.T) {
	tool := newAnalyzePerformanceTool(t)

	if tool.Name() != "analyze_performance" {
		t.Errorf("Name: got %q, want %q", tool.Name(), "analyze_performance")
	}
	if tool.Description() == "" {
		t.Error("Description should not be empty")
	}

	schema := tool.InputSchema()
	if schema.Type != "object" {
		t.Errorf("InputSchema.Type: got %q", schema.Type)
	}
	// projectPathOnlySchema: single required field, single property.
	wantRequired := map[string]bool{"projectPath": true}
	if len(schema.Required) != len(wantRequired) {
		t.Errorf("Required len: got %d, want %d", len(schema.Required), len(wantRequired))
	}
	for _, r := range schema.Required {
		if !wantRequired[r] {
			t.Errorf("unexpected required field: %s", r)
		}
	}
	if tool.OutputSchema() == nil {
		t.Error("expected non-nil OutputSchema")
	}
}

func TestAnalyzePerformanceTool_Execute_Success(t *testing.T) {
	tool := newAnalyzePerformanceTool(t)

	// One TS file with two of the analyzer's line-level checks:
	// a >100-char selector (simplify-selector) and a 5-digit timeout
	// value (avoid-high-timeouts). The large-file/large-file-size
	// rules would require 500+ lines / 50 KB of content, which is
	// wasteful for a unit test — the two line-level rules are enough
	// to prove the tool wires results through.
	projectDir := t.TempDir()
	longSelector := strings.Repeat("a", 120)
	writeTSFile(t, projectDir, "src/slow.ts", `
export function slow() {
  const sel = "`+longSelector+`";
  return { timeout: 30000 };
}
`)

	req := buildRequest(map[string]any{
		"projectPath": projectDir,
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if result.IsError {
		t.Fatalf("expected success, got IsError=true; content=%q", extractText(t, result))
	}

	var got models.ImprovementResponse
	if err := json.Unmarshal([]byte(extractText(t, result)), &got); err != nil {
		t.Fatalf("result not valid JSON: %v", err)
	}
	if len(got.Suggestions) < 2 {
		t.Fatalf("Suggestions len: got %d, want >=2; got=%+v", len(got.Suggestions), got.Suggestions)
	}

	seenRules := map[string]bool{}
	for _, s := range got.Suggestions {
		seenRules[s.Rule] = true
	}
	for _, want := range []string{"simplify-selector", "avoid-high-timeouts"} {
		if !seenRules[want] {
			t.Errorf("expected rule %q in suggestions, got rules=%v", want, seenRules)
		}
	}

	// Summary carries scanDurationMs and issuesFound. json unmarshals
	// numeric fields as float64.
	if got.Summary == nil {
		t.Fatal("expected non-nil Summary")
	}
	if _, ok := got.Summary["scanDurationMs"]; !ok {
		t.Error("expected scanDurationMs in Summary")
	}
	if issues, _ := got.Summary["issuesFound"].(float64); int(issues) < 2 {
		t.Errorf("Summary.issuesFound: got %v, want >=2", got.Summary["issuesFound"])
	}
}

func TestAnalyzePerformanceTool_Execute_MissingProjectPathErrors(t *testing.T) {
	tool := newAnalyzePerformanceTool(t)

	req := buildRequest(map[string]any{
		// projectPath deliberately omitted
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true when projectPath missing")
	}
	if !strings.Contains(extractText(t, result), "projectPath is required") {
		t.Errorf("expected required-field error, got %q", extractText(t, result))
	}
}

func TestAnalyzePerformanceTool_Execute_InvalidProjectPathTypeErrors(t *testing.T) {
	tool := newAnalyzePerformanceTool(t)

	// projectPath as a map — type assertion falls through, same
	// "projectPath is required" guard fires.
	req := buildRequest(map[string]any{
		"projectPath": map[string]any{"nope": true},
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true when projectPath is not a string")
	}
	if !strings.Contains(extractText(t, result), "projectPath is required") {
		t.Errorf("expected required-field error, got %q", extractText(t, result))
	}
}

func TestAnalyzePerformanceTool_Execute_AnalysisFailureSurfacesAsToolError(t *testing.T) {
	tool := newAnalyzePerformanceTool(t)

	// Non-existent projectPath — filepath.Walk fails, analyzer returns
	// the error, tool wraps with "Analysis failed:" prefix.
	req := buildRequest(map[string]any{
		"projectPath": filepath.Join(t.TempDir(), "does-not-exist"),
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute should not surface transport err, got: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true when analysis fails")
	}
	msg := extractText(t, result)
	if !strings.Contains(msg, "Analysis failed") {
		t.Errorf("expected 'Analysis failed' prefix, got %q", msg)
	}
}
