package generators

import (
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/orieken/saturday-mcp/internal/models"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/validators"
)

func TestMigrationGenerator_Generate(t *testing.T) {
	// Setup
	registry := templates.NewRegistry()
	
	// Reuse 'page' template content for test (simplified)
	templateContent := `
import { BasePage } from '@orieken/saturday-core';

export class {{pascalCase .Name}} extends BasePage {
  constructor(page: Page) {
    super(page);
    {{range .Elements}}
    this.add{{if eq .Type "button"}}Button{{else if eq .Type "input"}}Input{{else}}Element{{end}}("{{.Name}}", "{{.Selector}}");
    {{end}}
  }
}
`
	tmpl, err := template.New("page").Funcs(templates.DefaultFuncMap()).Parse(templateContent)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	if err := registry.Register("page", tmpl); err != nil {
		t.Fatalf("Failed to register template: %v", err)
	}

	cache := templates.NewCache(1 * time.Minute)
	processor := templates.NewProcessor(registry, cache)
	validator := validators.NewValidator()
	generator := NewMigrationGenerator(processor, validator)

	// Test Case
	sourceCode := `
test('basic test', async ({ page }) => {
  await page.click('#login-btn');
  await page.fill('input[name="user"]', 'admin');
  await page.click('text=Submit');
});
`
	req := models.MigrationRequest{
		SourceCode: sourceCode,
		Type:       "page",
	}

	// Execute
	resp, err := generator.Generate(req)
	if err != nil {
		t.Fatalf("Expected successful migration, got error: %v", err)
	}

	// Verify
	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	expectedSnippets := []string{
		"export class DraftMigrated extends BasePage",
		`this.addButton("button1", "#login-btn")`,
		`this.addInput("input2", "input[name=\"user\"]")`,
		// The 3rd one doesn't match standard heuristics so default element
		`this.addElement("element3", "text=Submit")`, 
	}

	for _, snippet := range expectedSnippets {
		if !strings.Contains(resp.Code, snippet) {
			t.Errorf("Expected code to contain '%s', but it didn't.\nCode:\n%s", snippet, resp.Code)
		}
	}
}
