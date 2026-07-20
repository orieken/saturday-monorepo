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

// newGenerateElementTool wires the tool with the real element generator.
// Same shape as generate_site's test helper.
func newGenerateElementTool(t *testing.T) *GenerateElementTool {
	t.Helper()
	gen := generators.NewElementGenerator(newProcessor(t), newValidator())
	return NewGenerateElementTool(silentLogger(), gen)
}

func TestGenerateElementTool_Metadata(t *testing.T) {
	tool := newGenerateElementTool(t)

	if tool.Name() != "generate_element" {
		t.Errorf("Name: got %q, want %q", tool.Name(), "generate_element")
	}
	if tool.Description() == "" {
		t.Error("Description should not be empty")
	}

	schema := tool.InputSchema()
	if schema.Type != "object" {
		t.Errorf("InputSchema.Type: got %q", schema.Type)
	}
	wantRequired := map[string]bool{"name": true, "rootSelector": true}
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

func TestGenerateElementTool_Execute_Success_NoWrite(t *testing.T) {
	tool := newGenerateElementTool(t)

	req := buildRequest(map[string]any{
		"name":         "navBar",
		"rootSelector": "nav.main-nav",
		"methods":      []any{"clickHome", "search"},
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
	if !got.Success {
		t.Error("expected success=true")
	}
	if !strings.Contains(got.Code, "NavBar") {
		t.Errorf("expected code to contain NavBar, got:\n%s", got.Code)
	}
	if !strings.Contains(got.Code, "extends BaseComponent") {
		t.Errorf("expected code to extend BaseComponent, got:\n%s", got.Code)
	}
	// Element filename has no "-element" suffix — just kebab-cased name.
	if got.FileName != "nav-bar.ts" {
		t.Errorf("FileName: got %q, want %q", got.FileName, "nav-bar.ts")
	}
	if got.Written {
		t.Error("expected Written=false when writeToFile omitted")
	}
	if got.FilePath != "" {
		t.Errorf("expected empty FilePath, got %q", got.FilePath)
	}
}

func TestGenerateElementTool_Execute_Success_WithWrite(t *testing.T) {
	tool := newGenerateElementTool(t)
	outDir := t.TempDir()

	req := buildRequest(map[string]any{
		"name":         "modal",
		"rootSelector": ".modal",
		"writeToFile":  true,
		"outputPath":   outDir,
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
	expectedPath := filepath.Join(outDir, "lib", "components", got.FileName)
	if got.FilePath != expectedPath {
		t.Errorf("FilePath: got %q, want %q", got.FilePath, expectedPath)
	}
	if _, err := os.Stat(expectedPath); err != nil {
		t.Errorf("expected file at %s, stat err: %v", expectedPath, err)
	}
}

func TestGenerateElementTool_Execute_WriteWithoutOutputPathErrors(t *testing.T) {
	tool := newGenerateElementTool(t)

	req := buildRequest(map[string]any{
		"name":         "modal",
		"rootSelector": ".modal",
		"writeToFile":  true,
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

func TestGenerateElementTool_Execute_WriteFailureSurfacesAsToolError(t *testing.T) {
	tool := newGenerateElementTool(t)
	outDir := t.TempDir()

	// Block MkdirAll by placing a file where a "lib" directory is
	// expected — same technique as the site test.
	if err := os.WriteFile(filepath.Join(outDir, "lib"), []byte("blocking"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	req := buildRequest(map[string]any{
		"name":         "modal",
		"rootSelector": ".modal",
		"writeToFile":  true,
		"outputPath":   outDir,
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

func TestGenerateElementTool_Execute_InvalidRequestSurfacesAsToolError(t *testing.T) {
	tool := newGenerateElementTool(t)

	// Missing rootSelector — required by the validator; the generator
	// returns a validation error which the tool must surface as an
	// IsError result rather than a transport error.
	req := buildRequest(map[string]any{
		"name": "bad",
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
