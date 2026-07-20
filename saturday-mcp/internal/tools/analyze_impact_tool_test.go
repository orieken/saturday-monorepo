package tools

import (
	"context"
	"encoding/json"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/orieken/saturday-mcp/internal/analyzers"
)

// newAnalyzeImpactTool wires the tool with a real GraphAnalyzer backed
// by the silent logger. Same "moist tests" rationale as the sibling
// analyzer tools — the graph is built from real files on disk against a
// t.TempDir fixture so the impact set the tool serializes back is what
// the analyzer actually produced, not a mock's stand-in.
func newAnalyzeImpactTool(t *testing.T) *AnalyzeImpactTool {
	t.Helper()
	analyzer := analyzers.NewGraphAnalyzer(silentLogger())
	return NewAnalyzeImpactTool(silentLogger(), analyzer)
}

func TestAnalyzeImpactTool_Metadata(t *testing.T) {
	tool := newAnalyzeImpactTool(t)

	if tool.Name() != "analyze_impact" {
		t.Errorf("Name: got %q, want %q", tool.Name(), "analyze_impact")
	}
	if tool.Description() == "" {
		t.Error("Description should not be empty")
	}

	schema := tool.InputSchema()
	if schema.Type != "object" {
		t.Errorf("InputSchema.Type: got %q", schema.Type)
	}
	// Both projectPath and targetFile are required — this tool differs
	// from the single-arg analyzers by requiring a specific file to
	// pivot the impact search around.
	wantRequired := map[string]bool{"projectPath": true, "targetFile": true}
	if len(schema.Required) != len(wantRequired) {
		t.Errorf("Required len: got %d, want %d", len(schema.Required), len(wantRequired))
	}
	for _, r := range schema.Required {
		if !wantRequired[r] {
			t.Errorf("unexpected required field: %s", r)
		}
	}
	if _, ok := schema.Properties["targetFile"]; !ok {
		t.Error("expected 'targetFile' property in schema")
	}
	if tool.OutputSchema() == nil {
		t.Error("expected non-nil OutputSchema")
	}
}

func TestAnalyzeImpactTool_Execute_Success(t *testing.T) {
	tool := newAnalyzeImpactTool(t)

	// Chain of TS files so the dependency graph has edges to walk:
	// login-page  <-imports  login-flow  <-imports  login-steps
	// The graph analyzer follows Incoming edges from the changed file,
	// so changing login-page should surface both the flow and the
	// steps as impacted. The import path regex requires the classic
	// `import ... from '...'` shape — a plain `import '...';` won't
	// match, so keep the from-clause explicit.
	projectDir := t.TempDir()
	writeTSFile(t, projectDir, "lib/pages/login-page.ts", `
export class LoginPage {}
`)
	writeTSFile(t, projectDir, "lib/flows/login-flow.ts", `
import { LoginPage } from '../pages/login-page';
export class LoginFlow {}
`)
	writeTSFile(t, projectDir, "tests/steps/login-steps.ts", `
import { LoginFlow } from '../../lib/flows/login-flow';
`)

	req := buildRequest(map[string]any{
		"projectPath": projectDir,
		"targetFile":  "lib/pages/login-page.ts",
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if result.IsError {
		t.Fatalf("expected success, got IsError=true; content=%q", extractText(t, result))
	}

	var got ImpactResult
	if err := json.Unmarshal([]byte(extractText(t, result)), &got); err != nil {
		t.Fatalf("result not valid JSON: %v", err)
	}
	if got.Target != "lib/pages/login-page.ts" {
		t.Errorf("Target: got %q, want %q", got.Target, "lib/pages/login-page.ts")
	}
	if got.Count != 2 {
		t.Errorf("Count: got %d, want 2 (impacted=%v)", got.Count, got.Impacted)
	}
	if got.Count != len(got.Impacted) {
		t.Errorf("Count/Impacted mismatch: Count=%d, len(Impacted)=%d", got.Count, len(got.Impacted))
	}

	sort.Strings(got.Impacted)
	want := []string{"lib/flows/login-flow.ts", "tests/steps/login-steps.ts"}
	sort.Strings(want)
	for i, w := range want {
		if got.Impacted[i] != w {
			t.Errorf("Impacted[%d]: got %q, want %q", i, got.Impacted[i], w)
		}
	}
}

func TestAnalyzeImpactTool_Execute_MissingProjectPathErrors(t *testing.T) {
	tool := newAnalyzeImpactTool(t)

	// Only targetFile supplied — the tool's single guard fires with a
	// combined message covering both required args (see line 62 of
	// analyze_impact_tool.go).
	req := buildRequest(map[string]any{
		"targetFile": "lib/pages/login-page.ts",
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true when projectPath missing")
	}
	if !strings.Contains(extractText(t, result), "projectPath and targetFile are required") {
		t.Errorf("expected required-fields error, got %q", extractText(t, result))
	}
}

func TestAnalyzeImpactTool_Execute_MissingTargetFileErrors(t *testing.T) {
	tool := newAnalyzeImpactTool(t)

	// Only projectPath supplied — same combined guard fires. Kept as a
	// distinct case so a future refactor that splits the error message
	// per-field will surface both here.
	req := buildRequest(map[string]any{
		"projectPath": t.TempDir(),
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true when targetFile missing")
	}
	if !strings.Contains(extractText(t, result), "projectPath and targetFile are required") {
		t.Errorf("expected required-fields error, got %q", extractText(t, result))
	}
}

func TestAnalyzeImpactTool_Execute_GraphBuildFailureSurfacesAsToolError(t *testing.T) {
	tool := newAnalyzeImpactTool(t)

	// Non-existent projectPath — filepath.Walk fails, GraphAnalyzer.Build
	// wraps that as "failed to walk project", and the tool prefixes with
	// "Failed to build dependency graph:" (distinct from the AnalyzeImpact
	// failure prefix — hitting both paths keeps the two error branches
	// covered).
	req := buildRequest(map[string]any{
		"projectPath": filepath.Join(t.TempDir(), "does-not-exist"),
		"targetFile":  "lib/pages/login-page.ts",
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute should not surface transport err, got: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true when graph build fails")
	}
	msg := extractText(t, result)
	if !strings.Contains(msg, "Failed to build dependency graph") {
		t.Errorf("expected 'Failed to build dependency graph' prefix, got %q", msg)
	}
}

func TestAnalyzeImpactTool_Execute_AnalysisFailureSurfacesAsToolError(t *testing.T) {
	tool := newAnalyzeImpactTool(t)

	// Valid, real projectPath (empty temp dir) — Build succeeds with
	// zero nodes, then AnalyzeImpact can't find the target (no exact
	// match, no fuzzy HasSuffix match) and returns "node not found".
	// The tool prefixes it with "Analysis failed:".
	req := buildRequest(map[string]any{
		"projectPath": t.TempDir(),
		"targetFile":  "lib/pages/does-not-exist.ts",
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
