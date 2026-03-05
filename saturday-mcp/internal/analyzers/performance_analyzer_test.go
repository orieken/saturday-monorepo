package analyzers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/orieken/saturday-mcp/internal/logging"
)

func TestPerformanceAnalyzer_Analyze(t *testing.T) {
	// Setup
	logger := logging.NewLogger(os.Stderr)
	analyzer := NewPerformanceAnalyzer(logger)
	
	// Create temp directory for test
	tmpDir, err := os.MkdirTemp("", "perf_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test file with issues
	badCode := `
export async function test() {
  // Bad: long complex selector
  await page.click("div.container > div.row > div.col-md-12 > form > div.form-group > label[for='username'] + input.form-control[type='text'][name='username'][placeholder='Enter Username']");
  
  // Bad: high timeout
  await page.waitForSelector('.foo', { timeout: 60000 });
}
`
	// Make selector definitely > 100 chars
	selector := "div.container > " + strings.Repeat("div.nested > span.item > ", 10) + "button[type='submit']"
	badCode = strings.Replace(badCode, `page.click("div.container`, `page.click("`+selector+`"); //`, 1)

	if err := os.WriteFile(filepath.Join(tmpDir, "test.ts"), []byte(badCode), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create a large file
	largeCode := strings.Repeat("// line\n", 600)
	if err := os.WriteFile(filepath.Join(tmpDir, "large.ts"), []byte(largeCode), 0644); err != nil {
		t.Fatalf("Failed to write large file: %v", err)
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

	// We expect:
	// 1. simplify-selector (from test.ts)
	// 2. avoid-high-timeouts (from test.ts)
	// 3. large-file (from large.ts)
	// 4. optionally large-file-size (from large.ts if > 50KB, probably not yet with just \n, but let's check issues count)
	
	if len(result.Suggestions) < 3 {
		t.Errorf("Expected at least 3 suggestions, got %d", len(result.Suggestions))
	}

	foundSelector := false
	foundTimeout := false
	foundLarge := false

	for _, s := range result.Suggestions {
		switch s.Rule {
		case "simplify-selector":
			foundSelector = true
		case "avoid-high-timeouts":
			foundTimeout = true
		case "large-file":
			foundLarge = true
		}
	}

	if !foundSelector {
		t.Error("Expected to find 'simplify-selector' suggestion")
	}
	if !foundTimeout {
		t.Error("Expected to find 'avoid-high-timeouts' suggestion")
	}
	if !foundLarge {
		t.Error("Expected to find 'large-file' suggestion")
	}

	// Check summary
	if _, ok := result.Summary["scanDurationMs"]; !ok {
		t.Error("Expected scanDurationMs in summary")
	}
}
