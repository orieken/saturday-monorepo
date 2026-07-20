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

// newSuggestImprovementsTool wires the tool with a real
// ImprovementAnalyzer (scans for waitForTimeout, console.log, and
// `: any` typing). Same moist-tests rationale as the other analyzer
// tests — we let the analyzer actually walk the fixture project.
func newSuggestImprovementsTool(t *testing.T) *SuggestImprovementsTool {
	t.Helper()
	analyzer := analyzers.NewImprovementAnalyzer(silentLogger())
	return NewSuggestImprovementsTool(silentLogger(), analyzer)
}

func TestSuggestImprovementsTool_Metadata(t *testing.T) {
	tool := newSuggestImprovementsTool(t)

	if tool.Name() != "suggest_improvements" {
		t.Errorf("Name: got %q, want %q", tool.Name(), "suggest_improvements")
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

func TestSuggestImprovementsTool_Execute_Success(t *testing.T) {
	tool := newSuggestImprovementsTool(t)

	// One TS file with all three anti-patterns the analyzer looks for:
	// a hard wait (no-hard-waits), a console.log (no-console-log), and
	// an `: any` type annotation (no-any). Expect one suggestion per
	// rule, hence three total.
	projectDir := t.TempDir()
	writeTSFile(t, projectDir, "src/bad.ts", `
export function bad(x: any) {
  console.log('hello');
  await page.waitForTimeout(1000);
  return x;
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
	if len(got.Suggestions) != 3 {
		t.Fatalf("Suggestions len: got %d, want 3; got=%+v", len(got.Suggestions), got.Suggestions)
	}

	seenRules := map[string]bool{}
	for _, s := range got.Suggestions {
		seenRules[s.Rule] = true
	}
	for _, want := range []string{"no-hard-waits", "no-console-log", "no-any"} {
		if !seenRules[want] {
			t.Errorf("expected rule %q in suggestions, got rules=%v", want, seenRules)
		}
	}

	// Summary is a free-form map[string]interface{}; the analyzer
	// records totalSuggestions there. json unmarshals the numeric
	// value as float64, so compare accordingly.
	if got.Summary == nil {
		t.Fatal("expected non-nil Summary")
	}
	if total, _ := got.Summary["totalSuggestions"].(float64); int(total) != 3 {
		t.Errorf("Summary.totalSuggestions: got %v, want 3", got.Summary["totalSuggestions"])
	}
}

func TestSuggestImprovementsTool_Execute_MissingProjectPathErrors(t *testing.T) {
	tool := newSuggestImprovementsTool(t)

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

func TestSuggestImprovementsTool_Execute_InvalidProjectPathTypeErrors(t *testing.T) {
	tool := newSuggestImprovementsTool(t)

	// projectPath as a slice — type assertion falls through, same
	// "projectPath is required" guard fires.
	req := buildRequest(map[string]any{
		"projectPath": []any{"nope"},
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

func TestSuggestImprovementsTool_Execute_AnalysisFailureSurfacesAsToolError(t *testing.T) {
	tool := newSuggestImprovementsTool(t)

	// Non-existent projectPath — filepath.Walk fails, analyzer returns
	// the error, tool wraps with "Analysis failed:" prefix (same
	// prefix as analyze_framework, distinct from validate_patterns'
	// "Validation failed:").
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
