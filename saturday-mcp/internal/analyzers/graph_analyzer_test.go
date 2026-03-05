package analyzers

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/orieken/saturday-mcp/internal/logging"
)

func TestGraphAnalyzer_AnalyzeImpact(t *testing.T) {
	// Setup temporary project
	root := t.TempDir()
	
	// mock file structure
	files := map[string]string{
		"lib/pages/LoginPage.ts": `
			export class LoginPage {}
		`,
		"lib/flows/LoginFlow.ts": `
			import { LoginPage } from '../pages/LoginPage';
			export class LoginFlow {}
		`,
		"tests/steps/login_steps.ts": `
			import { LoginFlow } from '../../lib/flows/LoginFlow';
			// step definitions...
		`,
		"lib/pages/OtherPage.ts": `
			export class OtherPage {}
		`,
	}

	for path, content := range files {
		fullPath := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Run Analysis
	logger := logging.NewLogger(os.Stderr)
	analyzer := NewGraphAnalyzer(logger)
	
	graph, err := analyzer.Build(root)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Case 1: Change LoginPage
	// Expected Impact: LoginFlow, login_steps
	impacted, err := analyzer.AnalyzeImpact(graph, "lib/pages/LoginPage.ts")
	if err != nil {
		t.Fatalf("AnalyzeImpact failed: %v", err)
	}

	sort.Strings(impacted)
	expected := []string{"lib/flows/LoginFlow.ts", "tests/steps/login_steps.ts"}
	sort.Strings(expected)

	if len(impacted) != len(expected) {
		t.Errorf("Expected %d impacted files, got %d: %v", len(expected), len(impacted), impacted)
	}

	for i := range expected {
		if impacted[i] != expected[i] {
			t.Errorf("Index %d: Expected %s, got %s", i, expected[i], impacted[i])
		}
	}

	// Case 2: Change OtherPage
	// Expected: None
	impacted2, err := analyzer.AnalyzeImpact(graph, "lib/pages/OtherPage.ts")
	if err != nil {
		t.Fatalf("AnalyzeImpact failed: %v", err)
	}
	if len(impacted2) != 0 {
		t.Errorf("Expected 0 impacted files, got %d", len(impacted2))
	}
}
