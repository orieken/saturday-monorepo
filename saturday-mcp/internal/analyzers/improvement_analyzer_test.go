package analyzers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/orieken/saturday-mcp/internal/logging"
)

func TestImprovementAnalyzer_Analyze(t *testing.T) {
	// Setup
	logger := logging.NewLogger(os.Stderr)
	analyzer := NewImprovementAnalyzer(logger)
	
	// Create temp directory for test
	tmpDir, err := os.MkdirTemp("", "improvement_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test file with issues
	badCode := `
import { Page } from 'playwright';

export async function test(page: Page) {
  // Bad practice: hard wait
  await page.waitForTimeout(5000);
  
  // Bad practice: console log
  console.log("Debugging something");

  // Bad practice: explicit any
  const data: any = { foo: 'bar' };
}
`
	filePath := filepath.Join(tmpDir, "test.ts")
	if err := os.WriteFile(filePath, []byte(badCode), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Execute
	result, err := analyzer.Analyze(tmpDir)
	if err != nil {
		t.Fatalf("Expected successful analysis, got error: %v", err)
	}

	// Verify
	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if len(result.Suggestions) != 3 {
		t.Errorf("Expected 3 suggestions, got %d", len(result.Suggestions))
	}

	// Check specific rules
	foundHardWait := false
	foundConsoleLog := false
	foundAny := false

	for _, s := range result.Suggestions {
		switch s.Rule {
		case "no-hard-waits":
			foundHardWait = true
		case "no-console-log":
			foundConsoleLog = true
		case "no-any":
			foundAny = true
		}
	}

	if !foundHardWait {
		t.Error("Expected to find 'no-hard-waits' suggestion")
	}
	if !foundConsoleLog {
		t.Error("Expected to find 'no-console-log' suggestion")
	}
	if !foundAny {
		t.Error("Expected to find 'no-any' suggestion")
	}
}
