package generators

import (
	"strings"
	"testing"
	"time"

	"github.com/orieken/saturday-mcp/internal/models"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/validators"
)

func setupPageGenerator(t *testing.T) *PageGenerator {
	t.Helper()

	registry := templates.NewRegistry()
	loader := templates.NewLoader(registry)
	cache := templates.NewCache(5 * time.Minute)
	processor := templates.NewProcessor(registry, cache)

	if err := loader.LoadAll(); err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	validator := validators.NewValidator()
	return NewPageGenerator(processor, validator)
}

func TestPageGenerator_Generate(t *testing.T) {
	generator := setupPageGenerator(t)

	t.Run("Valid page generation", func(t *testing.T) {
		req := models.PageGenerationRequest{
			Name: "login",
			Path: "/login",
			Elements: []models.ElementDefinition{
				{Name: "usernameInput", Selector: "#username"},
				{Name: "passwordInput", Selector: "#password"},
				{Name: "submitButton", Selector: "button[type='submit']"},
			},
			Description: "Login page",
		}

		resp, err := generator.Generate(req)
		if err != nil {
			t.Fatalf("Expected successful generation, got error: %v", err)
		}

		if resp == nil {
			t.Fatal("Expected response, got nil")
		}

		// Verify code content
		expectedStrings := []string{
			"LoginPage",
			"extends BasePage",
			"'/login'",
			"usernameInput",
			"passwordInput",
			"submitButton",
			"registerElement",
		}

		for _, expected := range expectedStrings {
			if !strings.Contains(resp.Code, expected) {
				t.Errorf("Expected code to contain %q", expected)
			}
		}

		// Verify filename
		if resp.FileName != "login-page.ts" {
			t.Errorf("Expected fileName 'login-page.ts', got %q", resp.FileName)
		}

		// Verify metadata
		if resp.Metadata["type"] != "page" {
			t.Errorf("Expected type 'page', got %q", resp.Metadata["type"])
		}

		if resp.Metadata["elementCount"] != "3" {
			t.Errorf("Expected elementCount '3', got %q", resp.Metadata["elementCount"])
		}
	})

	t.Run("Invalid request - missing name", func(t *testing.T) {
		req := models.PageGenerationRequest{
			Path: "/test",
			Elements: []models.ElementDefinition{
				{Name: "element1", Selector: "#el1"},
			},
		}

		_, err := generator.Generate(req)
		if err == nil {
			t.Error("Expected validation error for missing name")
		}
	})

	t.Run("Invalid request - empty elements", func(t *testing.T) {
		req := models.PageGenerationRequest{
			Name:     "test",
			Path:     "/test",
			Elements: []models.ElementDefinition{},
		}

		_, err := generator.Generate(req)
		if err == nil {
			t.Error("Expected validation error for empty elements")
		}
	})

	t.Run("Page with special characters in name", func(t *testing.T) {
		req := models.PageGenerationRequest{
			Name: "product_details",
			Path: "/product/:id",
			Elements: []models.ElementDefinition{
				{Name: "title", Selector: "h1"},
			},
		}

		resp, err := generator.Generate(req)
		if err != nil {
			t.Fatalf("Expected successful generation, got error: %v", err)
		}

		if !strings.Contains(resp.Code, "ProductDetailsPage") {
			t.Error("Expected name to be converted to PascalCase")
		}

		if resp.FileName != "product-details-page.ts" {
			t.Errorf("Expected kebab-case filename, got %q", resp.FileName)
		}
	})
}

func TestPageGenerator_ValidateRequest(t *testing.T) {
	generator := setupPageGenerator(t)

	t.Run("Valid request", func(t *testing.T) {
		req := models.PageGenerationRequest{
			Name: "test",
			Path: "/test",
			Elements: []models.ElementDefinition{
				{Name: "element", Selector: "#el"},
			},
		}

		err := generator.ValidateRequest(req)
		if err != nil {
			t.Errorf("Expected valid request, got error: %v", err)
		}
	})

	t.Run("Invalid request", func(t *testing.T) {
		req := models.PageGenerationRequest{
			Name: "test",
			// Missing Path and Elements
		}

		err := generator.ValidateRequest(req)
		if err == nil {
			t.Error("Expected validation error")
		}
	})
}
