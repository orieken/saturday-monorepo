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

func TestElementGenerator_Generate(t *testing.T) {
	// Setup
	registry := templates.NewRegistry()
	
	// Define mock template content
	templateContent := `
import { BaseComponent } from '@orieken/saturday-core';

export class {{pascalCase .Name}} extends BaseComponent {
  constructor(page: Page, rootLocator: Locator) {
    super(page, rootLocator);
  }

  {{range .Methods}}
  async {{camelCase .}}(): Promise<void> {}
  {{end}}
}
`
	// Parse template with function map
	tmpl, err := template.New("element").Funcs(templates.DefaultFuncMap()).Parse(templateContent)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	// Register parsed template
	if err := registry.Register("element", tmpl); err != nil {
		t.Fatalf("Failed to register template: %v", err)
	}

	cache := templates.NewCache(1 * time.Minute)
	processor := templates.NewProcessor(registry, cache)
	validator := validators.NewValidator()
	generator := NewElementGenerator(processor, validator)

	// Test Case
	req := models.ElementGenerationRequest{
		Name:         "NavBar",
		RootSelector: "nav.main-nav",
		Methods:      []string{"clickHome", "Search"},
		Description:  "Main navigation bar",
	}

	// Execute
	resp, err := generator.Generate(req)
	if err != nil {
		t.Fatalf("Expected successful generation, got error: %v", err)
	}

	// Verify
	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.FileName != "nav-bar.ts" {
		t.Errorf("Expected filename 'nav-bar.ts', got %s", resp.FileName)
	}

	expectedCodeSnippets := []string{
		"export class NavBar extends BaseComponent",
		"async clickHome(): Promise<void> {}",
		"async search(): Promise<void> {}",
	}

	for _, snippet := range expectedCodeSnippets {
		if !strings.Contains(resp.Code, snippet) {
			t.Errorf("Expected code to contain '%s', but it didn't.\nCode:\n%s", snippet, resp.Code)
		}
	}
}
