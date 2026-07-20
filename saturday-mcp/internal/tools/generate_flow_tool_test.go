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

// newGenerateFlowTool wires the tool with the real flow generator.
// Same shape as generate_site's test helper.
func newGenerateFlowTool(t *testing.T) *GenerateFlowTool {
	t.Helper()
	gen := generators.NewFlowGenerator(newProcessor(t), newValidator())
	return NewGenerateFlowTool(silentLogger(), gen)
}

func TestGenerateFlowTool_Metadata(t *testing.T) {
	tool := newGenerateFlowTool(t)

	if tool.Name() != "generate_flow" {
		t.Errorf("Name: got %q, want %q", tool.Name(), "generate_flow")
	}
	if tool.Description() == "" {
		t.Error("Description should not be empty")
	}

	schema := tool.InputSchema()
	if schema.Type != "object" {
		t.Errorf("InputSchema.Type: got %q", schema.Type)
	}
	wantRequired := map[string]bool{"name": true, "steps": true}
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

func TestGenerateFlowTool_Execute_Success_NoWrite(t *testing.T) {
	tool := newGenerateFlowTool(t)

	req := buildRequest(map[string]any{
		"name":  "checkout",
		"steps": []any{"addToCart", "enterPayment", "confirmOrder"},
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
	if !strings.Contains(got.Code, "CheckoutFlow") {
		t.Errorf("expected code to contain CheckoutFlow, got:\n%s", got.Code)
	}
	if !strings.Contains(got.Code, "extends BaseFlow") {
		t.Errorf("expected code to extend BaseFlow, got:\n%s", got.Code)
	}
	if got.FileName != "checkout-flow.ts" {
		t.Errorf("FileName: got %q, want %q", got.FileName, "checkout-flow.ts")
	}
	if got.Written {
		t.Error("expected Written=false when writeToFile omitted")
	}
	if got.FilePath != "" {
		t.Errorf("expected empty FilePath, got %q", got.FilePath)
	}
}

func TestGenerateFlowTool_Execute_Success_WithWrite(t *testing.T) {
	tool := newGenerateFlowTool(t)
	outDir := t.TempDir()

	req := buildRequest(map[string]any{
		"name":        "signup",
		"steps":       []any{"fillForm", "submit"},
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
	expectedPath := filepath.Join(outDir, "lib", "flows", got.FileName)
	if got.FilePath != expectedPath {
		t.Errorf("FilePath: got %q, want %q", got.FilePath, expectedPath)
	}
	if _, err := os.Stat(expectedPath); err != nil {
		t.Errorf("expected file at %s, stat err: %v", expectedPath, err)
	}
}

func TestGenerateFlowTool_Execute_WriteWithoutOutputPathErrors(t *testing.T) {
	tool := newGenerateFlowTool(t)

	req := buildRequest(map[string]any{
		"name":        "signup",
		"steps":       []any{"fillForm"},
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

func TestGenerateFlowTool_Execute_WriteFailureSurfacesAsToolError(t *testing.T) {
	tool := newGenerateFlowTool(t)
	outDir := t.TempDir()

	// Block MkdirAll by placing a file where a "lib" directory is
	// expected — same technique as the site test.
	if err := os.WriteFile(filepath.Join(outDir, "lib"), []byte("blocking"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	req := buildRequest(map[string]any{
		"name":        "signup",
		"steps":       []any{"fillForm"},
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

func TestGenerateFlowTool_Execute_InvalidRequestSurfacesAsToolError(t *testing.T) {
	tool := newGenerateFlowTool(t)

	// Empty steps slice fails min=1 — surfaced as a tool result error
	// rather than a transport error.
	req := buildRequest(map[string]any{
		"name":  "empty",
		"steps": []any{},
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
