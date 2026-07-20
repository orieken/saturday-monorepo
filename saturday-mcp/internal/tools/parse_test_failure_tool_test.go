package tools

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/orieken/saturday-mcp/internal/analyzers"
)

// newParseTestFailureTool wires the tool with a real TestLogAnalyzer.
// This tool is purely a string-in/JSON-out transform — no filesystem,
// no walking a project — so the "analyzer" here is genuinely just a
// regex-driven parser and letting it run for real keeps the parsed
// TestFailureInfo shape visible on the assertion path.
func newParseTestFailureTool(t *testing.T) *ParseTestFailureTool {
	t.Helper()
	analyzer := analyzers.NewTestLogAnalyzer(silentLogger())
	return NewParseTestFailureTool(silentLogger(), analyzer)
}

func TestParseTestFailureTool_Metadata(t *testing.T) {
	tool := newParseTestFailureTool(t)

	if tool.Name() != "parse_test_failure" {
		t.Errorf("Name: got %q, want %q", tool.Name(), "parse_test_failure")
	}
	if tool.Description() == "" {
		t.Error("Description should not be empty")
	}

	schema := tool.InputSchema()
	if schema.Type != "object" {
		t.Errorf("InputSchema.Type: got %q", schema.Type)
	}
	// Unlike the analyzer tools, parse_test_failure takes raw log output
	// rather than a filesystem path — the sole required arg is "output".
	wantRequired := map[string]bool{"output": true}
	if len(schema.Required) != len(wantRequired) {
		t.Errorf("Required len: got %d, want %d", len(schema.Required), len(wantRequired))
	}
	for _, r := range schema.Required {
		if !wantRequired[r] {
			t.Errorf("unexpected required field: %s", r)
		}
	}
	if _, ok := schema.Properties["output"]; !ok {
		t.Error("expected 'output' property in schema")
	}
	if tool.OutputSchema() == nil {
		t.Error("expected non-nil OutputSchema")
	}
}

func TestParseTestFailureTool_Execute_Success(t *testing.T) {
	tool := newParseTestFailureTool(t)

	// Hand-crafted Playwright-shaped log excerpt. The analyzer's header
	// regex is `\d+\)\s+\[.*?\]\s+›\s+(.*?):(\d+):\d+\s+›\s+(.*)`, so
	// each failure needs the numbered `N) [browser] › file:line:col ›
	// title` shape. Two failures let us confirm the tool serializes a
	// list, not just a singleton.
	logOutput := `
  1) [chromium] › tests/login.spec.ts:15:5 › Login Flow
    Timeout of 10000ms exceeded.
      15 |     await page.click('button#submit');

  2) [chromium] › tests/cart.spec.ts:42:10 › Cart Flow
    Error: expect(received).toBe(expected)
      42 |     expect(cartItems).toBe(5);
`

	req := buildRequest(map[string]any{
		"output": logOutput,
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if result.IsError {
		t.Fatalf("expected success, got IsError=true; content=%q", extractText(t, result))
	}

	var got []analyzers.TestFailureInfo
	if err := json.Unmarshal([]byte(extractText(t, result)), &got); err != nil {
		t.Fatalf("result not valid JSON: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("failures: got %d, want 2 (%+v)", len(got), got)
	}
	if got[0].File != "tests/login.spec.ts" {
		t.Errorf("failures[0].File: got %q, want tests/login.spec.ts", got[0].File)
	}
	if got[0].Line != 15 {
		t.Errorf("failures[0].Line: got %d, want 15", got[0].Line)
	}
	if got[1].File != "tests/cart.spec.ts" {
		t.Errorf("failures[1].File: got %q, want tests/cart.spec.ts", got[1].File)
	}
	if got[1].Line != 42 {
		t.Errorf("failures[1].Line: got %d, want 42", got[1].Line)
	}
}

func TestParseTestFailureTool_Execute_NoFailuresInLog(t *testing.T) {
	tool := newParseTestFailureTool(t)

	// Log with prose that matches no header shape — analyzer returns an
	// empty slice with no error, tool serializes to "null" (encoding/json
	// marshals nil slices as null). Either "null" or "[]" is a valid
	// success shape here; the point is IsError must be false.
	req := buildRequest(map[string]any{
		"output": "All tests passed successfully.\n",
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if result.IsError {
		t.Fatalf("expected success on clean log, got IsError=true; content=%q", extractText(t, result))
	}

	text := extractText(t, result)
	if text != "null" && text != "[]" {
		// If a future refactor pre-allocates the failures slice, the
		// serialized form flips from "null" to "[]" — both mean "no
		// failures", so we accept either rather than pinning to one.
		var got []analyzers.TestFailureInfo
		if err := json.Unmarshal([]byte(text), &got); err != nil {
			t.Fatalf("result not valid JSON: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected empty failures list, got %d entries", len(got))
		}
	}
}

func TestParseTestFailureTool_Execute_MissingOutputErrors(t *testing.T) {
	tool := newParseTestFailureTool(t)

	// No "output" key at all — the type assertion falls through and the
	// required-arg guard fires (line 57 of parse_test_failure_tool.go).
	req := buildRequest(map[string]any{})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true when output missing")
	}
	if !strings.Contains(extractText(t, result), "output argument is required") {
		t.Errorf("expected required-arg error, got %q", extractText(t, result))
	}
}

func TestParseTestFailureTool_Execute_EmptyOutputErrors(t *testing.T) {
	tool := newParseTestFailureTool(t)

	// Empty-string output — assertion succeeds but the value == ""
	// branch of the guard fires. Kept distinct from MissingOutput so a
	// future refactor that splits the two conditions is easy to notice.
	req := buildRequest(map[string]any{
		"output": "",
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true when output is empty")
	}
	if !strings.Contains(extractText(t, result), "output argument is required") {
		t.Errorf("expected required-arg error, got %q", extractText(t, result))
	}
}
