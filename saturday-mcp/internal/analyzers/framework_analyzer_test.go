package analyzers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/orieken/saturday-mcp/internal/logging"
)

func TestFrameworkAnalyzer_Analyze(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	
	commonSetup(t, tmpDir)

	// Create Analyzer
	logger := logging.NewLogger(os.Stderr)
	analyzer := NewFrameworkAnalyzer(logger)

	// Run Analysis
	result, err := analyzer.Analyze(tmpDir)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Verify Structure
	if result.Structure.PagesPath != filepath.Join(tmpDir, "lib", "pages") {
		t.Errorf("Expected PagesPath %s, got %s", filepath.Join(tmpDir, "lib", "pages"), result.Structure.PagesPath)
	}

	// Verify Stats
	if result.Stats.PageCount != 1 {
		t.Errorf("Expected 1 page, got %d", result.Stats.PageCount)
	}
	if result.Stats.FlowCount != 1 {
		t.Errorf("Expected 1 flow, got %d", result.Stats.FlowCount)
	}
	if result.Stats.StepDefCount != 2 {
		t.Errorf("Expected 2 steps, got %d", result.Stats.StepDefCount)
	}

	// Verify Inventory
	if len(result.Inventory.Pages) > 0 {
		if result.Inventory.Pages[0].Name != "LoginPage" {
			t.Errorf("Expected page name LoginPage, got %s", result.Inventory.Pages[0].Name)
		}
	} else {
		t.Error("No pages found in inventory")
	}

	if len(result.Inventory.Flows) > 0 {
		if result.Inventory.Flows[0].Name != "LoginFlow" {
			t.Errorf("Expected flow name LoginFlow, got %s", result.Inventory.Flows[0].Name)
		}
	} else {
		t.Error("No flows found in inventory")
	}
}

func commonSetup(t *testing.T, root string) {
	// Create directories
	dirs := []string{
		"lib/pages",
		"lib/flows",
		"tests/features",
		"tests/steps",
	}

	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(root, d), 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Create Page file
	pageContent := `
import { BasePage } from './BasePage';

export class LoginPage extends BasePage {
    constructor(page) {
        super(page);
    }
}
`
	if err := os.WriteFile(filepath.Join(root, "lib/pages/login-page.ts"), []byte(pageContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create Flow file
	flowContent := `
import { BaseFlow } from './BaseFlow';

export class LoginFlow extends BaseFlow {
    constructor(page) {
        super(page);
    }
}
`
	if err := os.WriteFile(filepath.Join(root, "lib/flows/login-flow.ts"), []byte(flowContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create Step Definitions file
	stepContent := `
import { Given, When, Then } from '@cucumber/cucumber';

Given('I am on the login page', async function() {
    await page.goto('/login');
});

When('I submit the form', async function() {
    await page.click('#submit');
});
`
	if err := os.WriteFile(filepath.Join(root, "tests/steps/login-steps.ts"), []byte(stepContent), 0644); err != nil {
		t.Fatal(err)
	}
}
