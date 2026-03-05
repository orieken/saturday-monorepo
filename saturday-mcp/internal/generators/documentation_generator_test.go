package generators

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/orieken/saturday-mcp/internal/models"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/validators"
)

func TestDocumentationGenerator_Generate(t *testing.T) {
	// Setup
	registry := templates.NewRegistry()
	
	templateContent := `
# Project Documentation
{{range .Pages}}
{{.Name}}
Elements: {{len .Elements}}
Methods: {{len .Methods}}
{{end}}
`
	tmpl, err := template.New("documentation").Parse(templateContent)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	if err := registry.Register("documentation", tmpl); err != nil {
		t.Fatalf("Failed to register template: %v", err)
	}

	cache := templates.NewCache(1 * time.Minute)
	processor := templates.NewProcessor(registry, cache)
	validator := validators.NewValidator()
	generator := NewDocumentationGenerator(processor, validator)

	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "doc_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create sample page
	pageContent := `
import { BasePage } from 'core';

export class LoginPage extends BasePage {
  constructor() {
    super();
    this.registerElement('username', '#user');
    this.registerElement('password', '#pass');
  }

  async login() {
    await this.click('username');
  }
  
  public async helper() {
  }
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "LoginPage.ts"), []byte(pageContent), 0644); err != nil {
		t.Fatalf("Failed to write page file: %v", err)
	}

	// Request
	req := models.DocumentationRequest{
		ProjectPath: tmpDir,
		OutputPath:  filepath.Join(tmpDir, "docs.md"),
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

	if !strings.Contains(resp.Code, "LoginPage") {
		t.Error("Expected LoginPage in docs")
	}
	if !strings.Contains(resp.Code, "Elements: 2") {
		t.Errorf("Expected 2 elements, got doc content:\n%s", resp.Code)
	}
	if !strings.Contains(resp.Code, "Methods: 2") { // login and helper
		t.Errorf("Expected 2 methods, got doc content:\n%s", resp.Code)
	}
}
