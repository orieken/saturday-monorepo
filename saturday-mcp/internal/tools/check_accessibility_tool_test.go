package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/orieken/saturday-mcp/internal/analyzers"
)

// Tests for CheckAccessibilityTool follow the same pattern established
// by analyze_complexity_tool_test.go: metadata assertions, an Execute
// success path against a t.TempDir() fixture, a missing-argument
// error path. Coverage target ≥ 85% per mcp-expand M1 Op 3 spec.

func TestCheckAccessibilityTool_Metadata(t *testing.T) {
	tool := NewCheckAccessibilityTool(silentLogger(), analyzers.NewAccessibilityAnalyzer())

	if tool.Name() != "check_accessibility" {
		t.Errorf("Name(): got %q, want %q", tool.Name(), "check_accessibility")
	}
	if tool.Description() == "" {
		t.Error("Description() should not be empty")
	}
	if tool.OutputSchema() == nil {
		t.Error("OutputSchema() should not be nil")
	}
	schema := tool.InputSchema()
	if schema.Type != "object" {
		t.Errorf("InputSchema Type: got %q, want object", schema.Type)
	}
	if _, ok := schema.Properties["filePath"]; !ok {
		t.Error("InputSchema is missing filePath property")
	}
	if _, ok := schema.Properties["projectPath"]; !ok {
		t.Error("InputSchema is missing projectPath property")
	}
}

func TestCheckAccessibilityTool_Execute_SingleFileWithViolations(t *testing.T) {
	tmp := t.TempDir()
	body := `<html>
<body>
<div onclick="doIt()">click</div>
<img src="pic.png">
<a href="/somewhere"></a>
<input type="text" id="orphan">
</body>
</html>
`
	filePath := filepath.Join(tmp, "page.html")
	if err := os.WriteFile(filePath, []byte(body), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	tool := NewCheckAccessibilityTool(silentLogger(), analyzers.NewAccessibilityAnalyzer())
	req := buildRequest(map[string]any{"filePath": filePath})

	res, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if res.IsError {
		t.Fatalf("Execute returned error result: %s", extractText(t, res))
	}

	var out AccessibilityReportResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if !out.Success {
		t.Error("expected Success=true")
	}
	if out.TotalFiles != 1 {
		t.Errorf("TotalFiles: got %d, want 1", out.TotalFiles)
	}
	if out.ViolationsCount != 5 {
		t.Errorf("ViolationsCount: got %d, want 5 (html-lang, clickable-non-interactive, img-alt, anchor-name, input-label)", out.ViolationsCount)
	}
	assertRulesPresent(t, out.Violations, "html-lang", "clickable-non-interactive", "img-alt", "anchor-name", "input-label")
}

func TestCheckAccessibilityTool_Execute_CleanFile(t *testing.T) {
	tmp := t.TempDir()
	body := `<html lang="en">
<body>
<button onclick="doIt()">click</button>
<img src="pic.png" alt="a picture">
<a href="/somewhere">go somewhere</a>
<a href="/other" aria-label="other place"></a>
<label for="named">Name</label>
<input type="text" id="named">
<input type="text" aria-label="Search">
<input type="hidden" name="csrf">
</body>
</html>
`
	filePath := filepath.Join(tmp, "clean.html")
	if err := os.WriteFile(filePath, []byte(body), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	tool := NewCheckAccessibilityTool(silentLogger(), analyzers.NewAccessibilityAnalyzer())
	req := buildRequest(map[string]any{"filePath": filePath})

	res, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out AccessibilityReportResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.ViolationsCount != 0 {
		t.Errorf("ViolationsCount: got %d, want 0; violations=%+v", out.ViolationsCount, out.Violations)
	}
	if out.Summary != "No accessibility violations found" {
		t.Errorf("Summary: got %q, want clean summary", out.Summary)
	}
}

func TestCheckAccessibilityTool_Execute_ProjectPathWalksDirectory(t *testing.T) {
	tmp := t.TempDir()
	// Two template files with violations plus a non-template file that
	// should be ignored, and a node_modules directory that should be
	// skipped entirely.
	writeFile(t, filepath.Join(tmp, "index.html"), `<html><img src="a.png"></html>`)
	writeFile(t, filepath.Join(tmp, "components/Widget.vue"), `<template><div onclick="x"/></template>`)
	writeFile(t, filepath.Join(tmp, "README.md"), `# not scanned`)
	writeFile(t, filepath.Join(tmp, "node_modules/pkg/junk.html"), `<img>`)

	tool := NewCheckAccessibilityTool(silentLogger(), analyzers.NewAccessibilityAnalyzer())
	req := buildRequest(map[string]any{"projectPath": tmp})

	res, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out AccessibilityReportResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.TotalFiles != 2 {
		t.Errorf("TotalFiles: got %d, want 2 (README.md excluded, node_modules skipped)", out.TotalFiles)
	}
	if out.ViolationsCount == 0 {
		t.Error("expected violations across the walked directory")
	}
}

func TestCheckAccessibilityTool_Execute_FilePathWinsOverProjectPath(t *testing.T) {
	tmp := t.TempDir()
	// A clean single file, and a dirty projectPath sibling.
	clean := filepath.Join(tmp, "clean.html")
	writeFile(t, clean, `<html lang="en"><body></body></html>`)
	writeFile(t, filepath.Join(tmp, "dirty.html"), `<html><img></html>`)

	tool := NewCheckAccessibilityTool(silentLogger(), analyzers.NewAccessibilityAnalyzer())
	req := buildRequest(map[string]any{
		"filePath":    clean,
		"projectPath": tmp,
	})

	res, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out AccessibilityReportResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.TotalFiles != 1 {
		t.Errorf("TotalFiles: got %d, want 1 (filePath should win)", out.TotalFiles)
	}
	if out.ViolationsCount != 0 {
		t.Errorf("ViolationsCount: got %d, want 0 (clean file was scanned)", out.ViolationsCount)
	}
}

func TestCheckAccessibilityTool_Execute_MissingArguments(t *testing.T) {
	tool := NewCheckAccessibilityTool(silentLogger(), analyzers.NewAccessibilityAnalyzer())
	req := buildRequest(map[string]any{})

	res, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if !res.IsError {
		t.Error("expected IsError=true when neither filePath nor projectPath is supplied")
	}
}

func TestCheckAccessibilityTool_Execute_NonexistentPath(t *testing.T) {
	tool := NewCheckAccessibilityTool(silentLogger(), analyzers.NewAccessibilityAnalyzer())
	req := buildRequest(map[string]any{"projectPath": filepath.Join(t.TempDir(), "does-not-exist")})

	res, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if !res.IsError {
		t.Error("expected IsError=true for a nonexistent target path")
	}
}

// writeFile writes body to path, creating parent directories as needed.
// Local to this test file (testfixtures_test.go's writeTSFile is TS-
// specific in intent).
func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// assertRulesPresent fails the test if any expected rule ID is not
// represented in the violations slice.
func assertRulesPresent(t *testing.T, violations []AccessibilityViolation, want ...string) {
	t.Helper()
	seen := map[string]struct{}{}
	for _, v := range violations {
		seen[v.Rule] = struct{}{}
	}
	for _, r := range want {
		if _, ok := seen[r]; !ok {
			t.Errorf("expected violation rule %q not found; got: %+v", r, violations)
		}
	}
}
