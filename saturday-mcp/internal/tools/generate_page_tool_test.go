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

// newGeneratePageTool wires the tool with the real page generator and
// template registry, matching the shape used by generate_site's test.
// See generate_site_tool_test.go for the template pattern this mirrors.
func newGeneratePageTool(t *testing.T) *GeneratePageTool {
	t.Helper()
	gen := generators.NewPageGenerator(newProcessor(t), newValidator())
	return NewGeneratePageTool(silentLogger(), gen)
}

func TestGeneratePageTool_Metadata(t *testing.T) {
	tool := newGeneratePageTool(t)

	if tool.Name() != "generate_page" {
		t.Errorf("Name: got %q, want %q", tool.Name(), "generate_page")
	}
	if tool.Description() == "" {
		t.Error("Description should not be empty")
	}

	schema := tool.InputSchema()
	if schema.Type != "object" {
		t.Errorf("InputSchema.Type: got %q", schema.Type)
	}
	wantRequired := map[string]bool{"name": true, "path": true, "elements": true}
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

func TestGeneratePageTool_Execute_Success_NoWrite(t *testing.T) {
	tool := newGeneratePageTool(t)

	req := buildRequest(map[string]any{
		"name": "login",
		"path": "/login",
		"elements": []any{
			map[string]any{"name": "usernameInput", "selector": "#username"},
			map[string]any{"name": "passwordInput", "selector": "#password"},
		},
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
	if !strings.Contains(got.Code, "LoginPage") {
		t.Errorf("expected code to contain LoginPage, got:\n%s", got.Code)
	}
	if !strings.Contains(got.Code, "extends BasePage") {
		t.Errorf("expected code to extend BasePage, got:\n%s", got.Code)
	}
	if got.FileName != "login-page.ts" {
		t.Errorf("FileName: got %q, want %q", got.FileName, "login-page.ts")
	}
	if got.Written {
		t.Error("expected Written=false when writeToFile omitted")
	}
	if got.FilePath != "" {
		t.Errorf("expected empty FilePath, got %q", got.FilePath)
	}
}

func TestGeneratePageTool_Execute_Success_WithWrite(t *testing.T) {
	tool := newGeneratePageTool(t)
	outDir := t.TempDir()

	req := buildRequest(map[string]any{
		"name": "cart",
		"path": "/cart",
		"elements": []any{
			map[string]any{"name": "checkoutButton", "selector": "#checkout"},
		},
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
	expectedPath := filepath.Join(outDir, "lib", "pages", got.FileName)
	if got.FilePath != expectedPath {
		t.Errorf("FilePath: got %q, want %q", got.FilePath, expectedPath)
	}
	if _, err := os.Stat(expectedPath); err != nil {
		t.Errorf("expected file at %s, stat err: %v", expectedPath, err)
	}
}

func TestGeneratePageTool_Execute_WriteWithoutOutputPathErrors(t *testing.T) {
	tool := newGeneratePageTool(t)

	req := buildRequest(map[string]any{
		"name": "cart",
		"path": "/cart",
		"elements": []any{
			map[string]any{"name": "checkoutButton", "selector": "#checkout"},
		},
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

func TestGeneratePageTool_Execute_WriteFailureSurfacesAsToolError(t *testing.T) {
	tool := newGeneratePageTool(t)
	outDir := t.TempDir()

	// Block the writer by placing a regular file where it needs a "lib"
	// directory. MkdirAll fails because the path component is not a
	// directory — same trick as the generate_site write-failure test.
	if err := os.WriteFile(filepath.Join(outDir, "lib"), []byte("blocking"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	req := buildRequest(map[string]any{
		"name": "cart",
		"path": "/cart",
		"elements": []any{
			map[string]any{"name": "checkoutButton", "selector": "#checkout"},
		},
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

func TestGeneratePageTool_Execute_InvalidRequestSurfacesAsToolError(t *testing.T) {
	tool := newGeneratePageTool(t)

	// Empty elements slice fails the min=1 validator; the generator
	// returns a validation error which the tool must surface as an
	// IsError result rather than a transport error.
	req := buildRequest(map[string]any{
		"name":     "empty",
		"path":     "/empty",
		"elements": []any{},
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
