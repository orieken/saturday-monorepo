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

func TestServiceGenerator_Generate(t *testing.T) {
	// Setup
	registry := templates.NewRegistry()
	
	// Define mock template content
	templateContent := `
import { BaseService, APIRequestContext } from '@orieken/saturday-core';

export class {{pascalCase .Name}}Service extends BaseService {
  constructor(request: APIRequestContext, baseURL: string = '{{.BaseURL}}') {
    super(request, baseURL);
  }

  {{range .Endpoints}}
  async {{camelCase .Name}}(): Promise<any> {
    return this.{{lower .Method}}('{{.Path}}');
  }
  {{end}}
}
`
	// Parse template with function map
	tmpl, err := template.New("service").Funcs(templates.DefaultFuncMap()).Parse(templateContent)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	// Register parsed template
	if err := registry.Register("service", tmpl); err != nil {
		t.Fatalf("Failed to register template: %v", err)
	}

	cache := templates.NewCache(1 * time.Minute)
	processor := templates.NewProcessor(registry, cache)
	validator := validators.NewValidator()
	generator := NewServiceGenerator(processor, validator)

	// Test Case
	req := models.ServiceGenerationRequest{
		Name:    "User",
		BaseURL: "https://api.example.com",
		Endpoints: []models.EndpointDefinition{
			{Name: "GetUser", Method: "GET", Path: "/users/1"},
			{Name: "CreateUser", Method: "POST", Path: "/users"},
		},
		Description: "User management service",
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

	if resp.FileName != "user-service.ts" {
		t.Errorf("Expected filename 'user-service.ts', got %s", resp.FileName)
	}

	expectedCodeSnippets := []string{
		"export class UserService extends BaseService",
		"constructor(request: APIRequestContext, baseURL: string = 'https://api.example.com')",
		"async getUser(): Promise<any>",
		"return this.get('/users/1')",
		"async createUser(): Promise<any>",
		"return this.post('/users')",
	}

	for _, snippet := range expectedCodeSnippets {
		if !strings.Contains(resp.Code, snippet) {
			t.Errorf("Expected code to contain '%s', but it didn't.\nCode:\n%s", snippet, resp.Code)
		}
	}
}
