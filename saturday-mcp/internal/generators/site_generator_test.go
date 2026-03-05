package generators

import (
	"strings"
	"testing"
	"time"

	"github.com/orieken/saturday-mcp/internal/models"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/validators"
)

func setupSiteGenerator(t *testing.T) *SiteGenerator {
	t.Helper()

	// Setup template system
	registry := templates.NewRegistry()
	loader := templates.NewLoader(registry)
	cache := templates.NewCache(5 * time.Minute)
	processor := templates.NewProcessor(registry, cache)

	// Load templates
	if err := loader.LoadAll(); err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	// Setup validator
	validator := validators.NewValidator()

	return NewSiteGenerator(processor, validator)
}

func TestSiteGenerator_Generate(t *testing.T) {
	generator := setupSiteGenerator(t)

	t.Run("Valid site generation", func(t *testing.T) {
		req := models.SiteGenerationRequest{
			Name:        "ecommerce",
			BaseURL:     "https://example.com",
			Pages:       []string{"home", "product", "cart"},
			Flows:       []string{"checkout"},
			Description: "E-commerce site",
		}

		resp, err := generator.Generate(req)
		if err != nil {
			t.Fatalf("Expected successful generation, got error: %v", err)
		}

		if resp == nil {
			t.Fatal("Expected response, got nil")
		}

		// Verify code content
		if len(resp.Code) == 0 {
			t.Error("Expected generated code to be non-empty")
		}

		// Check for expected content in generated code
		expectedStrings := []string{
			"EcommerceSite",
			"extends BaseSite",
			"HomePage",
			"ProductPage",
			"CartPage",
			"registerPage",
		}

		for _, expected := range expectedStrings {
			if !strings.Contains(resp.Code, expected) {
				t.Errorf("Expected code to contain %q", expected)
			}
		}

		// Verify filename
		if resp.FileName != "ecommerce-site.ts" {
			t.Errorf("Expected filename 'ecommerce-site.ts', got %q", resp.FileName)
		}

		// Verify metadata
		if resp.Metadata["type"] != "site" {
			t.Errorf("Expected type 'site', got %q", resp.Metadata["type"])
		}

		if resp.Metadata["name"] != "ecommerce" {
			t.Errorf("Expected name 'ecommerce', got %q", resp.Metadata["name"])
		}

		if resp.Metadata["pageCount"] != "3" {
			t.Errorf("Expected pageCount '3', got %q", resp.Metadata["pageCount"])
		}
	})

	t.Run("Invalid request - missing name", func(t *testing.T) {
		req := models.SiteGenerationRequest{
			BaseURL: "https://example.com",
			Pages:   []string{"home"},
		}

		_, err := generator.Generate(req)
		if err == nil {
			t.Error("Expected validation error for missing name")
		}

		if !strings.Contains(err.Error(), "validation failed") {
			t.Errorf("Expected validation error, got: %v", err)
		}
	})

	t.Run("Invalid request - invalid URL", func(t *testing.T) {
		req := models.SiteGenerationRequest{
			Name:    "test",
			BaseURL: "not-a-url",
			Pages:   []string{"home"},
		}

		_, err := generator.Generate(req)
		if err == nil {
			t.Error("Expected validation error for invalid URL")
		}
	})

	t.Run("Invalid request - empty pages", func(t *testing.T) {
		req := models.SiteGenerationRequest{
			Name:    "test",
			BaseURL: "https://example.com",
			Pages:   []string{},
		}

		_, err := generator.Generate(req)
		if err == nil {
			t.Error("Expected validation error for empty pages")
		}
	})

	t.Run("Site with special characters in name", func(t *testing.T) {
		req := models.SiteGenerationRequest{
			Name:    "my_awesome_site",
			BaseURL: "https://example.com",
			Pages:   []string{"home"},
		}

		resp, err := generator.Generate(req)
		if err != nil {
			t.Fatalf("Expected successful generation, got error: %v", err)
		}

		// Should convert to PascalCase
		if !strings.Contains(resp.Code, "MyAwesomeSite") {
			t.Error("Expected name to be converted to PascalCase")
		}

		// Filename should be kebab-case
		if resp.FileName != "my-awesome-site-site.ts" {
			t.Errorf("Expected kebab-case filename, got %q", resp.FileName)
		}
	})

	t.Run("Site without flows", func(t *testing.T) {
		req := models.SiteGenerationRequest{
			Name:    "simple",
			BaseURL: "https://example.com",
			Pages:   []string{"home"},
		}

		resp, err := generator.Generate(req)
		if err != nil {
			t.Fatalf("Expected successful generation, got error: %v", err)
		}

		if len(resp.Code) == 0 {
			t.Error("Expected generated code to be non-empty")
		}
	})

	t.Run("Site with many pages", func(t *testing.T) {
		req := models.SiteGenerationRequest{
			Name:    "large",
			BaseURL: "https://example.com",
			Pages: []string{
				"home", "about", "contact", "products",
				"cart", "checkout", "account", "orders",
			},
		}

		resp, err := generator.Generate(req)
		if err != nil {
			t.Fatalf("Expected successful generation, got error: %v", err)
		}

		if resp.Metadata["pageCount"] != "8" {
			t.Errorf("Expected pageCount '8', got %q", resp.Metadata["pageCount"])
		}

		// Verify all pages are registered
		for _, page := range req.Pages {
			pascalPage := templates.PascalCase(page) + "Page"
			if !strings.Contains(resp.Code, pascalPage) {
				t.Errorf("Expected code to contain %q", pascalPage)
			}
		}
	})
}

func TestSiteGenerator_ValidateRequest(t *testing.T) {
	generator := setupSiteGenerator(t)

	t.Run("Valid request", func(t *testing.T) {
		req := models.SiteGenerationRequest{
			Name:    "test",
			BaseURL: "https://example.com",
			Pages:   []string{"home"},
		}

		err := generator.ValidateRequest(req)
		if err != nil {
			t.Errorf("Expected valid request, got error: %v", err)
		}
	})

	t.Run("Invalid request", func(t *testing.T) {
		req := models.SiteGenerationRequest{
			Name:  "test",
			Pages: []string{"home"},
			// Missing BaseURL
		}

		err := generator.ValidateRequest(req)
		if err == nil {
			t.Error("Expected validation error")
		}
	})
}
