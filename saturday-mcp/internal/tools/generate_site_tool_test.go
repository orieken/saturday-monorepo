package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/orieken/saturday-mcp/internal/generators"
)

// newGenerateSiteTool builds the tool with real generator + template
// system on top of a shared logger. Same pattern for every generator
// tool test — the generators themselves are pure functions over their
// request struct, so we exercise them for real rather than faking.
func newGenerateSiteTool(t *testing.T) *GenerateSiteTool {
	t.Helper()
	gen := generators.NewSiteGenerator(newProcessor(t), newValidator())
	return NewGenerateSiteTool(silentLogger(), gen)
}

func TestGenerateSiteTool_Metadata(t *testing.T) {
	tool := newGenerateSiteTool(t)

	if tool.Name() != "generate_site" {
		t.Errorf("Name: got %q, want %q", tool.Name(), "generate_site")
	}
	if tool.Description() == "" {
		t.Error("Description should not be empty")
	}

	schema := tool.InputSchema()
	if schema.Type != "object" {
		t.Errorf("InputSchema.Type: got %q", schema.Type)
	}
	wantRequired := map[string]bool{"name": true, "baseUrl": true, "pages": true}
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

func TestGenerateSiteTool_Execute_Success_NoWrite(t *testing.T) {
	tool := newGenerateSiteTool(t)

	req := buildRequest(map[string]any{
		"name":    "ecommerce",
		"baseUrl": "https://example.com",
		"pages":   []any{"home", "cart"},
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if result.IsError {
		t.Errorf("expected success, got IsError=true; content=%q", extractText(t, result))
	}

	var got GenerationResult
	if err := json.Unmarshal([]byte(extractText(t, result)), &got); err != nil {
		t.Fatalf("result not valid JSON: %v", err)
	}
	if !got.Success {
		t.Error("expected success=true")
	}
	if !strings.Contains(got.Code, "EcommerceSite") {
		t.Errorf("expected code to contain EcommerceSite, got:\n%s", got.Code)
	}
	if got.FileName != "ecommerce-site.ts" {
		t.Errorf("FileName: got %q", got.FileName)
	}
	if got.Written {
		t.Error("expected Written=false when writeToFile omitted")
	}
	if got.FilePath != "" {
		t.Errorf("expected empty FilePath, got %q", got.FilePath)
	}
}

func TestGenerateSiteTool_Execute_Success_WithWrite(t *testing.T) {
	tool := newGenerateSiteTool(t)
	outDir := t.TempDir()

	req := buildRequest(map[string]any{
		"name":        "storefront",
		"baseUrl":     "https://example.com",
		"pages":       []any{"home"},
		"writeToFile": true,
		"outputPath":  outDir,
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if result.IsError {
		t.Fatalf("expected success, got IsError=true; content=%q", extractText(t, result))
	}

	var got GenerationResult
	if err := json.Unmarshal([]byte(extractText(t, result)), &got); err != nil {
		t.Fatalf("result not valid JSON: %v", err)
	}
	if !got.Written {
		t.Error("expected Written=true")
	}
	expectedPath := filepath.Join(outDir, "lib", "sites", got.FileName)
	if got.FilePath != expectedPath {
		t.Errorf("FilePath: got %q, want %q", got.FilePath, expectedPath)
	}
	// The tool must actually write the file — that's the point of writeToFile.
	if _, err := os.Stat(expectedPath); err != nil {
		t.Errorf("expected file at %s, stat err: %v", expectedPath, err)
	}
}

func TestGenerateSiteTool_Execute_WriteWithoutOutputPathErrors(t *testing.T) {
	tool := newGenerateSiteTool(t)

	req := buildRequest(map[string]any{
		"name":        "storefront",
		"baseUrl":     "https://example.com",
		"pages":       []any{"home"},
		"writeToFile": true,
		// outputPath deliberately omitted
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true when writeToFile=true but outputPath missing")
	}
	if !strings.Contains(extractText(t, result), "outputPath is required") {
		t.Errorf("expected outputPath error, got %q", extractText(t, result))
	}
}

func TestGenerateSiteTool_Execute_WriteFailureSurfacesAsToolError(t *testing.T) {
	tool := newGenerateSiteTool(t)
	outDir := t.TempDir()

	// Put a file where the writer needs a directory ("lib" under
	// outputPath) — MkdirAll will fail because the path component is not
	// a directory. This is the closest we can get to a real filesystem
	// failure without depending on OS-specific permission semantics.
	if err := os.WriteFile(filepath.Join(outDir, "lib"), []byte("blocking"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	req := buildRequest(map[string]any{
		"name":        "storefront",
		"baseUrl":     "https://example.com",
		"pages":       []any{"home"},
		"writeToFile": true,
		"outputPath":  outDir,
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected IsError=true when write fails")
	}
	if !strings.Contains(extractText(t, result), "Failed to write file") {
		t.Errorf("expected write-failure prefix, got %q", extractText(t, result))
	}
}

func TestGenerateSiteTool_Execute_InvalidRequestSurfacesAsToolError(t *testing.T) {
	tool := newGenerateSiteTool(t)

	// baseUrl fails the URL validator; the generator returns a validation
	// error which the tool must surface as an IsError result rather than a
	// transport error.
	req := buildRequest(map[string]any{
		"name":    "bad",
		"baseUrl": "not-a-url",
		"pages":   []any{"home"},
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute should not surface transport err, got: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true for invalid request")
	}
	if !strings.Contains(extractText(t, result), "Generation failed") {
		t.Errorf("expected 'Generation failed' prefix, got %q", extractText(t, result))
	}
}
