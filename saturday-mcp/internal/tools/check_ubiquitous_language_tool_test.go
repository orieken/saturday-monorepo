package tools

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/orieken/saturday-mcp/internal/analyzers"
)

// Tests for CheckUbiquitousLanguageTool follow the pattern established
// by check_accessibility_tool_test.go: metadata assertions, an Execute
// success path against a t.TempDir() fixture, and error paths for
// missing arguments and a missing dictionary. Coverage target ≥ 85% per
// mcp-expand M1 Op 4 spec. writeFile / silentLogger / buildRequest /
// extractText come from testfixtures_test.go.

func TestCheckUbiquitousLanguageTool_Metadata(t *testing.T) {
	tool := NewCheckUbiquitousLanguageTool(silentLogger(), analyzers.NewUbiquitousLanguageAnalyzer())

	if tool.Name() != "check_ubiquitous_language" {
		t.Errorf("Name(): got %q, want %q", tool.Name(), "check_ubiquitous_language")
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
	if len(schema.Required) != 2 {
		t.Errorf("Required: got %v, want two entries", schema.Required)
	}
	if _, ok := schema.Properties["projectPath"]; !ok {
		t.Error("InputSchema is missing projectPath property")
	}
	if _, ok := schema.Properties["dictionaryPath"]; !ok {
		t.Error("InputSchema is missing dictionaryPath property")
	}
}

func TestCheckUbiquitousLanguageTool_Execute_TableFormatViolations(t *testing.T) {
	tmp := t.TempDir()
	dict := filepath.Join(tmp, "DOMAIN_DICTIONARY.md")
	writeFile(t, dict, "# Domain Dictionary\n\n"+
		"## 2. Entities\n\n"+
		"| Term | Domain | Description | Synonyms to AVOID |\n"+
		"|---|---|---|---|\n"+
		"| **User** | Auth | End user of the system | `client`, `customer` (too generic) |\n")
	writeFile(t, filepath.Join(tmp, "code.go"),
		"package x\n// Look up a client record\n// Also a customer note\n")

	tool := NewCheckUbiquitousLanguageTool(silentLogger(), analyzers.NewUbiquitousLanguageAnalyzer())
	req := buildRequest(map[string]any{
		"projectPath":    tmp,
		"dictionaryPath": dict,
	})

	res, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if res.IsError {
		t.Fatalf("Execute returned error result: %s", extractText(t, res))
	}

	var out UbiquitousLanguageResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if !out.Success {
		t.Error("expected Success=true")
	}
	if out.ViolationsCount != 2 {
		t.Fatalf("ViolationsCount: got %d, want 2; violations=%+v", out.ViolationsCount, out.Violations)
	}
	for _, v := range out.Violations {
		if v.Suggested != "User" {
			t.Errorf("Suggested: got %q, want User", v.Suggested)
		}
	}
}

func TestCheckUbiquitousLanguageTool_Execute_SectionFormatViolations(t *testing.T) {
	tmp := t.TempDir()
	dict := filepath.Join(tmp, "DOMAIN_DICTIONARY.md")
	writeFile(t, dict, "# Dictionary\n\n"+
		"## User\n"+
		"End user of the system.\n"+
		"**Synonyms to avoid**: `client`, `customer`\n")
	writeFile(t, filepath.Join(tmp, "svc.py"),
		"# fetch client here\n    pass\n")

	tool := NewCheckUbiquitousLanguageTool(silentLogger(), analyzers.NewUbiquitousLanguageAnalyzer())
	req := buildRequest(map[string]any{
		"projectPath":    tmp,
		"dictionaryPath": dict,
	})

	res, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out UbiquitousLanguageResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.ViolationsCount != 1 {
		t.Fatalf("ViolationsCount: got %d, want 1; violations=%+v", out.ViolationsCount, out.Violations)
	}
	v := out.Violations[0]
	if v.InvalidTerm != "client" {
		t.Errorf("InvalidTerm: got %q, want client", v.InvalidTerm)
	}
	if v.Suggested != "User" {
		t.Errorf("Suggested: got %q, want User", v.Suggested)
	}
	if v.LineNumber != 1 {
		t.Errorf("LineNumber: got %d, want 1", v.LineNumber)
	}
}

func TestCheckUbiquitousLanguageTool_Execute_CommaSeparatedNoBackticks(t *testing.T) {
	tmp := t.TempDir()
	dict := filepath.Join(tmp, "DOMAIN_DICTIONARY.md")
	writeFile(t, dict, "## User\n**Not**: client, customer\n")
	writeFile(t, filepath.Join(tmp, "code.ts"), "const client = null;\n")

	tool := NewCheckUbiquitousLanguageTool(silentLogger(), analyzers.NewUbiquitousLanguageAnalyzer())
	req := buildRequest(map[string]any{
		"projectPath":    tmp,
		"dictionaryPath": dict,
	})

	res, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out UbiquitousLanguageResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.ViolationsCount != 1 {
		t.Fatalf("ViolationsCount: got %d, want 1; violations=%+v", out.ViolationsCount, out.Violations)
	}
	if out.Violations[0].Suggested != "User" {
		t.Errorf("Suggested: got %q, want User", out.Violations[0].Suggested)
	}
}

func TestCheckUbiquitousLanguageTool_Execute_CleanCode(t *testing.T) {
	tmp := t.TempDir()
	dict := filepath.Join(tmp, "DOMAIN_DICTIONARY.md")
	writeFile(t, dict, "## User\n**Synonyms to avoid**: `client`\n")
	writeFile(t, filepath.Join(tmp, "code.go"), "package x\nfunc handleUser() {}\n")

	tool := NewCheckUbiquitousLanguageTool(silentLogger(), analyzers.NewUbiquitousLanguageAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{
		"projectPath":    tmp,
		"dictionaryPath": dict,
	}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out UbiquitousLanguageResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.ViolationsCount != 0 {
		t.Errorf("ViolationsCount: got %d, want 0; violations=%+v", out.ViolationsCount, out.Violations)
	}
	if out.Summary != "No ubiquitous language violations found" {
		t.Errorf("Summary: got %q, want clean summary", out.Summary)
	}
}

func TestCheckUbiquitousLanguageTool_Execute_SkipsExcludedDirsAndNonSource(t *testing.T) {
	tmp := t.TempDir()
	dict := filepath.Join(tmp, "DOMAIN_DICTIONARY.md")
	writeFile(t, dict, "## User\n**Synonyms to avoid**: `client`\n")
	// One real source file that should be scanned.
	writeFile(t, filepath.Join(tmp, "src/main.go"), "package x\n// client here\n")
	// Excluded directories that should NOT be scanned.
	writeFile(t, filepath.Join(tmp, "node_modules/lib/x.js"), "// client\n")
	writeFile(t, filepath.Join(tmp, "vendor/dep/y.go"), "// client\n")
	writeFile(t, filepath.Join(tmp, "dist/out.ts"), "// client\n")
	writeFile(t, filepath.Join(tmp, "build/bundle.js"), "// client\n")
	writeFile(t, filepath.Join(tmp, ".venv/lib/z.py"), "# client\n")
	writeFile(t, filepath.Join(tmp, ".git/config.go"), "// client\n")
	// Non-source files that should be ignored regardless of directory.
	writeFile(t, filepath.Join(tmp, "src/README.md"), "client\n")
	writeFile(t, filepath.Join(tmp, "src/data.json"), "{\"client\":1}\n")

	tool := NewCheckUbiquitousLanguageTool(silentLogger(), analyzers.NewUbiquitousLanguageAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{
		"projectPath":    tmp,
		"dictionaryPath": dict,
	}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out UbiquitousLanguageResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.ViolationsCount != 1 {
		t.Errorf("ViolationsCount: got %d, want 1 (only src/main.go should be scanned); violations=%+v", out.ViolationsCount, out.Violations)
	}
}

func TestCheckUbiquitousLanguageTool_Execute_EmptyDictionary(t *testing.T) {
	tmp := t.TempDir()
	dict := filepath.Join(tmp, "DOMAIN_DICTIONARY.md")
	writeFile(t, dict, "# Dictionary\n\nNo parseable terms here.\n")
	writeFile(t, filepath.Join(tmp, "code.go"), "package x\n// client customer profile role\n")

	tool := NewCheckUbiquitousLanguageTool(silentLogger(), analyzers.NewUbiquitousLanguageAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{
		"projectPath":    tmp,
		"dictionaryPath": dict,
	}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if res.IsError {
		t.Fatalf("expected non-error result for empty dictionary: %s", extractText(t, res))
	}

	var out UbiquitousLanguageResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.ViolationsCount != 0 {
		t.Errorf("ViolationsCount: got %d, want 0 (empty dictionary should yield zero violations)", out.ViolationsCount)
	}
}

func TestCheckUbiquitousLanguageTool_Execute_NumberedHeadingsIgnored(t *testing.T) {
	// Regression guard for updateTermFromHeader: a section like
	// "## 1. Core Domains" is organizational, not a term name — synonym
	// lines directly under it must not be attributed to "1. Core Domains".
	tmp := t.TempDir()
	dict := filepath.Join(tmp, "DOMAIN_DICTIONARY.md")
	writeFile(t, dict, "## 1. Core Domains\n**Synonyms to avoid**: `client`\n")
	writeFile(t, filepath.Join(tmp, "code.go"), "package x\n// client\n")

	tool := NewCheckUbiquitousLanguageTool(silentLogger(), analyzers.NewUbiquitousLanguageAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{
		"projectPath":    tmp,
		"dictionaryPath": dict,
	}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out UbiquitousLanguageResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.ViolationsCount != 0 {
		t.Errorf("ViolationsCount: got %d, want 0 (numbered heading should be ignored)", out.ViolationsCount)
	}
}

func TestCheckUbiquitousLanguageTool_Execute_MissingArguments(t *testing.T) {
	tool := NewCheckUbiquitousLanguageTool(silentLogger(), analyzers.NewUbiquitousLanguageAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if !res.IsError {
		t.Error("expected IsError=true when both arguments are missing")
	}
}

func TestCheckUbiquitousLanguageTool_Execute_MissingDictionaryFile(t *testing.T) {
	tmp := t.TempDir()
	writeFile(t, filepath.Join(tmp, "code.go"), "package x\n")

	tool := NewCheckUbiquitousLanguageTool(silentLogger(), analyzers.NewUbiquitousLanguageAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{
		"projectPath":    tmp,
		"dictionaryPath": filepath.Join(tmp, "does-not-exist.md"),
	}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if !res.IsError {
		t.Error("expected IsError=true when the dictionary file doesn't exist")
	}
}

func TestCheckUbiquitousLanguageTool_Execute_MissingProjectPath(t *testing.T) {
	tmp := t.TempDir()
	dict := filepath.Join(tmp, "DOMAIN_DICTIONARY.md")
	writeFile(t, dict, "## User\n**Synonyms to avoid**: `client`\n")

	tool := NewCheckUbiquitousLanguageTool(silentLogger(), analyzers.NewUbiquitousLanguageAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{
		"projectPath":    filepath.Join(tmp, "does-not-exist"),
		"dictionaryPath": dict,
	}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if !res.IsError {
		t.Error("expected IsError=true when the project path doesn't exist")
	}
}
