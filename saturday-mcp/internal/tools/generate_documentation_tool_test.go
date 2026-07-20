package tools

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/orieken/saturday-mcp/internal/generators"
)

// newGenerateDocumentationTool wires the tool with the real documentation
// generator and a captureFileWriter (see testfixtures_test.go). Unlike the
// generator tools in Batch 1, generate_documentation takes the low-level
// domain.FileWriter interface — not the filewriter.FileWriter concrete —
// so we inject a fake here rather than fighting with real MkdirAll trick.
func newGenerateDocumentationTool(t *testing.T, fw *captureFileWriter) *GenerateDocumentationTool {
	t.Helper()
	gen := generators.NewDocumentationGenerator(newProcessor(t), newValidator())
	return NewGenerateDocumentationTool(silentLogger(), gen, fw)
}

func TestGenerateDocumentationTool_Metadata(t *testing.T) {
	tool := newGenerateDocumentationTool(t, &captureFileWriter{})

	if tool.Name() != "generate_documentation" {
		t.Errorf("Name: got %q, want %q", tool.Name(), "generate_documentation")
	}
	if tool.Description() == "" {
		t.Error("Description should not be empty")
	}

	schema := tool.InputSchema()
	if schema.Type != "object" {
		t.Errorf("InputSchema.Type: got %q", schema.Type)
	}
	// Both projectPath and outputPath are required at the schema layer —
	// the tool's Execute also guards on both being non-empty.
	wantRequired := map[string]bool{"projectPath": true, "outputPath": true}
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

func TestGenerateDocumentationTool_Execute_Success(t *testing.T) {
	fw := &captureFileWriter{}
	tool := newGenerateDocumentationTool(t, fw)

	// Build a minimal Saturday-shaped project on disk that the
	// documentation generator can scan. Two page files (Login, Cart)
	// plus a non-page file that must be ignored.
	projectDir := t.TempDir()
	writeTSFile(t, projectDir, "pages/login-page.ts", `
import { BasePage } from '@orieken/saturday-core';

export class LoginPage extends BasePage {
  protected initializeElements(): void {
    this.registerElement('username', '#username');
    this.registerElement('password', '#password');
  }

  public async login(): Promise<void> { /* ... */ }
}
`)
	writeTSFile(t, projectDir, "pages/cart-page.ts", `
export class CartPage extends BasePage {
  protected initializeElements(): void {
    this.registerElement('checkout', '#checkout');
  }

  public async checkout(): Promise<void> { /* ... */ }
}
`)
	// Non-page TS file — no "Page" or "page" in path/name, so scanPages
	// must skip it entirely. Also has no export class.
	writeTSFile(t, projectDir, "utils/helpers.ts", `export function noop() {}`)

	outPath := filepath.Join(t.TempDir(), "PROJECT_DOCS.md")
	req := buildRequest(map[string]any{
		"projectPath": projectDir,
		"outputPath":  outPath,
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if result.IsError {
		t.Fatalf("expected success, got IsError=true; content=%q", extractText(t, result))
	}

	var got DocumentationResult
	if err := json.Unmarshal([]byte(extractText(t, result)), &got); err != nil {
		t.Fatalf("result not valid JSON: %v", err)
	}
	if !got.Success {
		t.Error("expected success=true")
	}
	if got.Path != outPath {
		t.Errorf("Path: got %q, want %q", got.Path, outPath)
	}
	// Two page files scanned — Login + Cart.
	if got.Pages != "2" {
		t.Errorf("Pages: got %q, want %q", got.Pages, "2")
	}

	// The tool must have delegated the actual write to the injected
	// FileWriter with the caller's outputPath and the rendered markdown.
	if !fw.called {
		t.Fatal("expected FileWriter.WriteFile to be called")
	}
	if fw.lastPath != outPath {
		t.Errorf("WriteFile path: got %q, want %q", fw.lastPath, outPath)
	}
	body := string(fw.lastBody)
	if !strings.Contains(body, "# Project Documentation") {
		t.Errorf("expected doc header in body, got:\n%s", body)
	}
	if !strings.Contains(body, "LoginPage") {
		t.Errorf("expected LoginPage in body, got:\n%s", body)
	}
	if !strings.Contains(body, "CartPage") {
		t.Errorf("expected CartPage in body, got:\n%s", body)
	}
}

func TestGenerateDocumentationTool_Execute_MissingProjectPathErrors(t *testing.T) {
	fw := &captureFileWriter{}
	tool := newGenerateDocumentationTool(t, fw)

	req := buildRequest(map[string]any{
		// projectPath deliberately omitted
		"outputPath": "/tmp/whatever.md",
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true when projectPath missing")
	}
	if !strings.Contains(extractText(t, result), "projectPath and outputPath are required") {
		t.Errorf("expected required-field error, got %q", extractText(t, result))
	}
	if fw.called {
		t.Error("FileWriter should not be called when validation fails up front")
	}
}

func TestGenerateDocumentationTool_Execute_MissingOutputPathErrors(t *testing.T) {
	fw := &captureFileWriter{}
	tool := newGenerateDocumentationTool(t, fw)

	req := buildRequest(map[string]any{
		"projectPath": t.TempDir(),
		// outputPath deliberately omitted
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true when outputPath missing")
	}
	if !strings.Contains(extractText(t, result), "projectPath and outputPath are required") {
		t.Errorf("expected required-field error, got %q", extractText(t, result))
	}
	if fw.called {
		t.Error("FileWriter should not be called when validation fails up front")
	}
}

func TestGenerateDocumentationTool_Execute_GenerationFailureSurfacesAsToolError(t *testing.T) {
	fw := &captureFileWriter{}
	tool := newGenerateDocumentationTool(t, fw)

	// Non-existent projectPath — filepath.Walk returns an error the
	// generator wraps as "failed to scan pages", which the tool surfaces
	// with a "Generation failed" prefix.
	req := buildRequest(map[string]any{
		"projectPath": filepath.Join(t.TempDir(), "does-not-exist"),
		"outputPath":  filepath.Join(t.TempDir(), "docs.md"),
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute should not surface transport err, got: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true when generator fails")
	}
	msg := extractText(t, result)
	if !strings.Contains(msg, "Generation failed") {
		t.Errorf("expected 'Generation failed' prefix, got %q", msg)
	}
	if fw.called {
		t.Error("FileWriter should not be called when generation fails")
	}
}

func TestGenerateDocumentationTool_Execute_WriteFailureSurfacesAsToolError(t *testing.T) {
	// Inject an fs error via captureFileWriter — the tool must catch it
	// and surface as a tool result error rather than a transport error.
	fw := &captureFileWriter{err: errors.New("disk full")}
	tool := newGenerateDocumentationTool(t, fw)

	projectDir := t.TempDir()
	writeTSFile(t, projectDir, "pages/dummy-page.ts", `export class DummyPage {}`)

	req := buildRequest(map[string]any{
		"projectPath": projectDir,
		"outputPath":  filepath.Join(t.TempDir(), "docs.md"),
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected IsError=true when write fails")
	}
	msg := extractText(t, result)
	if !strings.Contains(msg, "Failed to write file") {
		t.Errorf("expected 'Failed to write file' prefix, got %q", msg)
	}
	if !strings.Contains(msg, "disk full") {
		t.Errorf("expected underlying error surfaced, got %q", msg)
	}
	if !fw.called {
		t.Error("expected FileWriter.WriteFile to have been called even though it errored")
	}
}
