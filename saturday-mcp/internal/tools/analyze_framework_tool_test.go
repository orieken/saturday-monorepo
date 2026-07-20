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

// newAnalyzeFrameworkTool wires the tool with a real FrameworkAnalyzer
// backed by the silent logger. Analyzers walk a project tree — we let
// them run for real against a t.TempDir fixture rather than fake them,
// matching the "moist tests" guidance in shared/rules/testing-
// conventions.md so the code the tool actually returns stays visible on
// the assertion path.
func newAnalyzeFrameworkTool(t *testing.T) *AnalyzeFrameworkTool {
	t.Helper()
	analyzer := analyzers.NewFrameworkAnalyzer(silentLogger())
	return NewAnalyzeFrameworkTool(silentLogger(), analyzer)
}

func TestAnalyzeFrameworkTool_Metadata(t *testing.T) {
	tool := newAnalyzeFrameworkTool(t)

	if tool.Name() != "analyze_framework" {
		t.Errorf("Name: got %q, want %q", tool.Name(), "analyze_framework")
	}
	if tool.Description() == "" {
		t.Error("Description should not be empty")
	}

	schema := tool.InputSchema()
	if schema.Type != "object" {
		t.Errorf("InputSchema.Type: got %q", schema.Type)
	}
	// Only projectPath is required; patterns is an optional list.
	wantRequired := map[string]bool{"projectPath": true}
	if len(schema.Required) != len(wantRequired) {
		t.Errorf("Required len: got %d, want %d", len(schema.Required), len(wantRequired))
	}
	for _, r := range schema.Required {
		if !wantRequired[r] {
			t.Errorf("unexpected required field: %s", r)
		}
	}
	if _, ok := schema.Properties["patterns"]; !ok {
		t.Error("expected optional 'patterns' property in schema")
	}
	if tool.OutputSchema() == nil {
		t.Error("expected non-nil OutputSchema")
	}
}

func TestAnalyzeFrameworkTool_Execute_Success(t *testing.T) {
	tool := newAnalyzeFrameworkTool(t)

	// Build a minimal Saturday-shaped project on disk so the analyzer
	// finds structure (pagesPath / flowsPath), classes (LoginPage /
	// CheckoutFlow), and step definitions all in one walk. The class
	// regexes require the opening "{" on the same line as the class
	// declaration and rely on a "Page"/"Flow" suffix — the recorded
	// name is the captured prefix + the suffix (see the class-name-
	// suffix trap called out in Batch 2 notes).
	projectDir := t.TempDir()
	writeTSFile(t, projectDir, "lib/pages/login-page.ts", `
import { BasePage } from '@orieken/saturday-core';
export class LoginPage extends BasePage {}
`)
	writeTSFile(t, projectDir, "lib/flows/checkout-flow.ts", `
import { BaseFlow } from '@orieken/saturday-core';
export class CheckoutFlow extends BaseFlow {}
`)
	writeTSFile(t, projectDir, "tests/steps/login-steps.ts", `
import { Given, When } from '@cucumber/cucumber';
Given('I am on the login page', async function() {});
When('I submit the form', async function() {});
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

	var got models.AnalysisResult
	if err := json.Unmarshal([]byte(extractText(t, result)), &got); err != nil {
		t.Fatalf("result not valid JSON: %v", err)
	}
	if got.Stats.PageCount != 1 {
		t.Errorf("PageCount: got %d, want 1", got.Stats.PageCount)
	}
	if got.Stats.FlowCount != 1 {
		t.Errorf("FlowCount: got %d, want 1", got.Stats.FlowCount)
	}
	if got.Stats.StepDefCount != 2 {
		t.Errorf("StepDefCount: got %d, want 2", got.Stats.StepDefCount)
	}
	if got.Structure.RootPath != projectDir {
		t.Errorf("RootPath: got %q, want %q", got.Structure.RootPath, projectDir)
	}
	if len(got.Inventory.Pages) != 1 || got.Inventory.Pages[0].Name != "LoginPage" {
		t.Errorf("Inventory.Pages: got %+v, want [LoginPage]", got.Inventory.Pages)
	}
	if len(got.Inventory.Flows) != 1 || got.Inventory.Flows[0].Name != "CheckoutFlow" {
		t.Errorf("Inventory.Flows: got %+v, want [CheckoutFlow]", got.Inventory.Flows)
	}
}

func TestAnalyzeFrameworkTool_Execute_MissingProjectPathErrors(t *testing.T) {
	tool := newAnalyzeFrameworkTool(t)

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

func TestAnalyzeFrameworkTool_Execute_InvalidProjectPathTypeErrors(t *testing.T) {
	tool := newAnalyzeFrameworkTool(t)

	// projectPath is not a string — the tool's type assertion falls
	// through to the empty-string zero value and the same
	// "projectPath is required" guard fires. Kept as a separate case
	// so a future refactor that adds a distinct type-mismatch message
	// is easy to notice.
	req := buildRequest(map[string]any{
		"projectPath": 42,
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

func TestAnalyzeFrameworkTool_Execute_AnalysisFailureSurfacesAsToolError(t *testing.T) {
	tool := newAnalyzeFrameworkTool(t)

	// Non-existent projectPath — filepath.Walk returns an error the
	// analyzer surfaces verbatim, and the tool wraps it with an
	// "Analysis failed:" prefix. Using a child of TempDir avoids
	// stepping on real filesystem permissions.
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
