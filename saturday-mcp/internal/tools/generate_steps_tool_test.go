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

// newGenerateStepsTool wires the tool with the real step generator.
// Same shape as generate_site's test helper.
func newGenerateStepsTool(t *testing.T) *GenerateStepsTool {
	t.Helper()
	gen := generators.NewStepGenerator(newProcessor(t), newValidator())
	return NewGenerateStepsTool(silentLogger(), gen)
}

func TestGenerateStepsTool_Metadata(t *testing.T) {
	tool := newGenerateStepsTool(t)

	if tool.Name() != "generate_steps" {
		t.Errorf("Name: got %q, want %q", tool.Name(), "generate_steps")
	}
	if tool.Description() == "" {
		t.Error("Description should not be empty")
	}

	schema := tool.InputSchema()
	if schema.Type != "object" {
		t.Errorf("InputSchema.Type: got %q", schema.Type)
	}
	wantRequired := map[string]bool{"steps": true}
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

func TestGenerateStepsTool_Execute_Success_NoWrite(t *testing.T) {
	tool := newGenerateStepsTool(t)

	req := buildRequest(map[string]any{
		"steps": []any{
			map[string]any{"type": "Given", "pattern": "I am on the {string} page", "parameters": "pageName: string"},
			map[string]any{"type": "When", "pattern": "I click {string}", "parameters": "buttonName: string"},
			map[string]any{"type": "Then", "pattern": "I should see {string}", "parameters": "text: string"},
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
	if !strings.Contains(got.Code, "@cucumber/cucumber") {
		t.Errorf("expected code to import @cucumber/cucumber, got:\n%s", got.Code)
	}
	if !strings.Contains(got.Code, "Given('I am on the {string} page'") {
		t.Errorf("expected code to bind Given step, got:\n%s", got.Code)
	}
	// Default language is TypeScript.
	if got.FileName != "steps.ts" {
		t.Errorf("FileName: got %q, want %q", got.FileName, "steps.ts")
	}
	if got.Written {
		t.Error("expected Written=false when writeToFile omitted")
	}
}

func TestGenerateStepsTool_Execute_Success_JavaScriptFilename(t *testing.T) {
	tool := newGenerateStepsTool(t)

	req := buildRequest(map[string]any{
		"language": "javascript",
		"steps": []any{
			map[string]any{"type": "Given", "pattern": "I have a user"},
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
	if got.FileName != "steps.js" {
		t.Errorf("FileName: got %q, want %q", got.FileName, "steps.js")
	}
}

func TestGenerateStepsTool_Execute_Success_WithWrite(t *testing.T) {
	tool := newGenerateStepsTool(t)
	outDir := t.TempDir()

	req := buildRequest(map[string]any{
		"steps": []any{
			map[string]any{"type": "Given", "pattern": "I am logged in"},
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
	// Steps write into tests/steps/, not lib/*.
	expectedPath := filepath.Join(outDir, "tests", "steps", got.FileName)
	if got.FilePath != expectedPath {
		t.Errorf("FilePath: got %q, want %q", got.FilePath, expectedPath)
	}
	if _, err := os.Stat(expectedPath); err != nil {
		t.Errorf("expected file at %s, stat err: %v", expectedPath, err)
	}
}

func TestGenerateStepsTool_Execute_WriteWithoutOutputPathErrors(t *testing.T) {
	tool := newGenerateStepsTool(t)

	req := buildRequest(map[string]any{
		"steps": []any{
			map[string]any{"type": "Given", "pattern": "I am logged in"},
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

func TestGenerateStepsTool_Execute_WriteFailureSurfacesAsToolError(t *testing.T) {
	tool := newGenerateStepsTool(t)
	outDir := t.TempDir()

	// Block MkdirAll by placing a file where a "tests" directory is
	// expected — steps write to tests/steps/, not lib/, so we block the
	// tests component.
	if err := os.WriteFile(filepath.Join(outDir, "tests"), []byte("blocking"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	req := buildRequest(map[string]any{
		"steps": []any{
			map[string]any{"type": "Given", "pattern": "I am logged in"},
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

func TestGenerateStepsTool_Execute_InvalidRequestSurfacesAsToolError(t *testing.T) {
	tool := newGenerateStepsTool(t)

	// Empty steps slice fails min=1; the generator returns a validation
	// error which the tool must surface as an IsError result rather than
	// a transport error.
	req := buildRequest(map[string]any{
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
