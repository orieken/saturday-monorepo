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

// newValidatePatternsTool wires the tool with a real PatternValidator
// (NamingConvention + Inheritance rules). Same "moist tests" reasoning
// as the other analyzer tools — a fake validator would hide the actual
// rule output the tool has to serialize.
func newValidatePatternsTool(t *testing.T) *ValidatePatternsTool {
	t.Helper()
	validator := analyzers.NewPatternValidator(silentLogger())
	return NewValidatePatternsTool(silentLogger(), validator)
}

func TestValidatePatternsTool_Metadata(t *testing.T) {
	tool := newValidatePatternsTool(t)

	if tool.Name() != "validate_patterns" {
		t.Errorf("Name: got %q, want %q", tool.Name(), "validate_patterns")
	}
	if tool.Description() == "" {
		t.Error("Description should not be empty")
	}

	schema := tool.InputSchema()
	if schema.Type != "object" {
		t.Errorf("InputSchema.Type: got %q", schema.Type)
	}
	// Only projectPath is required; checkTypes is an optional list.
	wantRequired := map[string]bool{"projectPath": true}
	if len(schema.Required) != len(wantRequired) {
		t.Errorf("Required len: got %d, want %d", len(schema.Required), len(wantRequired))
	}
	for _, r := range schema.Required {
		if !wantRequired[r] {
			t.Errorf("unexpected required field: %s", r)
		}
	}
	if _, ok := schema.Properties["checkTypes"]; !ok {
		t.Error("expected optional 'checkTypes' property in schema")
	}
	if tool.OutputSchema() == nil {
		t.Error("expected non-nil OutputSchema")
	}
}

func TestValidatePatternsTool_Execute_Success(t *testing.T) {
	tool := newValidatePatternsTool(t)

	// Two page files under lib/pages/: one violates naming (class
	// WrongName in invalid-naming-page.ts), one violates inheritance
	// (LoginPage without extends BasePage). Both rules should fire so
	// Valid=false with 2 issues.
	projectDir := t.TempDir()
	writeTSFile(t, projectDir, "lib/pages/invalid-naming-page.ts", `
import { BasePage } from '@orieken/saturday-core';
export class WrongName extends BasePage {}
`)
	writeTSFile(t, projectDir, "lib/pages/login-page.ts", `
export class LoginPage {}
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

	var got models.ValidationResponse
	if err := json.Unmarshal([]byte(extractText(t, result)), &got); err != nil {
		t.Fatalf("result not valid JSON: %v", err)
	}
	if got.Valid {
		t.Error("expected Valid=false when rule violations are present")
	}
	if len(got.Issues) < 2 {
		t.Fatalf("expected >=2 issues, got %d: %+v", len(got.Issues), got.Issues)
	}

	// One issue must be a NamingConvention hit on invalid-naming-page.ts;
	// another must be an Inheritance hit on login-page.ts.
	namingHit := false
	inheritanceHit := false
	for _, issue := range got.Issues {
		if issue.Rule == "NamingConvention" && strings.HasSuffix(issue.File, "invalid-naming-page.ts") {
			namingHit = true
		}
		if issue.Rule == "Inheritance" && strings.HasSuffix(issue.File, "login-page.ts") {
			inheritanceHit = true
		}
	}
	if !namingHit {
		t.Error("expected NamingConvention issue on invalid-naming-page.ts")
	}
	if !inheritanceHit {
		t.Error("expected Inheritance issue on login-page.ts")
	}
}

func TestValidatePatternsTool_Execute_MissingProjectPathErrors(t *testing.T) {
	tool := newValidatePatternsTool(t)

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

func TestValidatePatternsTool_Execute_InvalidProjectPathTypeErrors(t *testing.T) {
	tool := newValidatePatternsTool(t)

	// Non-string projectPath — same "projectPath is required" guard
	// fires via zero-value fallthrough. See notes in
	// analyze_framework_tool_test.go for rationale.
	req := buildRequest(map[string]any{
		"projectPath": true,
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

func TestValidatePatternsTool_Execute_ValidationFailureSurfacesAsToolError(t *testing.T) {
	tool := newValidatePatternsTool(t)

	// Non-existent projectPath — filepath.Walk fails, the validator
	// surfaces it verbatim, and the tool wraps it with a
	// "Validation failed:" prefix (distinct from the "Analysis failed:"
	// prefix on the other analyzer tools — see Batch 2 note about
	// per-tool error prefixes).
	req := buildRequest(map[string]any{
		"projectPath": filepath.Join(t.TempDir(), "does-not-exist"),
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute should not surface transport err, got: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true when validation fails")
	}
	msg := extractText(t, result)
	if !strings.Contains(msg, "Validation failed") {
		t.Errorf("expected 'Validation failed' prefix, got %q", msg)
	}
}
