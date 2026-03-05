package generators

import (
	"strings"
	"testing"
	"time"

	"github.com/orieken/saturday-mcp/internal/models"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/validators"
)

func setupStepGenerator(t *testing.T) *StepGenerator {
	t.Helper()

	registry := templates.NewRegistry()
	loader := templates.NewLoader(registry)
	cache := templates.NewCache(5 * time.Minute)
	processor := templates.NewProcessor(registry, cache)

	if err := loader.LoadAll(); err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	validator := validators.NewValidator()
	return NewStepGenerator(processor, validator)
}

func TestStepGenerator_Generate(t *testing.T) {
	generator := setupStepGenerator(t)

	t.Run("Valid step generation - TypeScript", func(t *testing.T) {
		req := models.StepDefinitionRequest{
			Steps: []models.StepDefinition{
				{Type: "Given", Pattern: "I am on the {string} page", Parameters: "pageName: string"},
				{Type: "When", Pattern: "I click the {string} button", Parameters: "buttonName: string"},
				{Type: "Then", Pattern: "I should see {string}", Parameters: "text: string"},
			},
			Language: "typescript",
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
			"Given('I am on the {string} page'",
			"When('I click the {string} button'",
			"Then('I should see {string}'",
			"pageName: string",
			"buttonName: string",
			"text: string",
			"@cucumber/cucumber",
		}

		for _, expected := range expectedStrings {
			if !strings.Contains(resp.Code, expected) {
				t.Errorf("Expected code to contain %q", expected)
			}
		}

		// Verify filename
		if resp.FileName != "steps.ts" {
			t.Errorf("Expected fileName 'steps.ts', got %q", resp.FileName)
		}

		// Verify metadata
		if resp.Metadata["type"] != "steps" {
			t.Errorf("Expected type 'steps', got %q", resp.Metadata["type"])
		}

		if resp.Metadata["stepCount"] != "3" {
			t.Errorf("Expected stepCount '3', got %q", resp.Metadata["stepCount"])
		}

		if resp.Metadata["language"] != "typescript" {
			t.Errorf("Expected language 'typescript', got %q", resp.Metadata["language"])
		}
	})

	t.Run("Valid step generation - JavaScript", func(t *testing.T) {
		req := models.StepDefinitionRequest{
			Steps: []models.StepDefinition{
				{Type: "Given", Pattern: "I have a user", Parameters: ""},
			},
			Language: "javascript",
		}

		resp, err := generator.Generate(req)
		if err != nil {
			t.Fatalf("Expected successful generation, got error: %v", err)
		}

		if resp.FileName != "steps.js" {
			t.Errorf("Expected fileName 'steps.js', got %q", resp.FileName)
		}

		if resp.Metadata["language"] != "javascript" {
			t.Errorf("Expected language 'javascript', got %q", resp.Metadata["language"])
		}
	})

	t.Run("Default language to TypeScript", func(t *testing.T) {
		req := models.StepDefinitionRequest{
			Steps: []models.StepDefinition{
				{Type: "Given", Pattern: "I am logged in", Parameters: ""},
			},
			// No language specified
		}

		resp, err := generator.Generate(req)
		if err != nil {
			t.Fatalf("Expected successful generation, got error: %v", err)
		}

		if resp.FileName != "steps.ts" {
			t.Errorf("Expected default to TypeScript, got %q", resp.FileName)
		}

		if resp.Metadata["language"] != "typescript" {
			t.Errorf("Expected language 'typescript', got %q", resp.Metadata["language"])
		}
	})

	t.Run("Invalid request - missing steps", func(t *testing.T) {
		req := models.StepDefinitionRequest{
			Steps: []models.StepDefinition{},
		}

		_, err := generator.Generate(req)
		if err == nil {
			t.Error("Expected validation error for empty steps")
		}
	})

	t.Run("Invalid request - invalid step type", func(t *testing.T) {
		req := models.StepDefinitionRequest{
			Steps: []models.StepDefinition{
				{Type: "Invalid", Pattern: "test", Parameters: ""},
			},
		}

		_, err := generator.Generate(req)
		if err == nil {
			t.Error("Expected validation error for invalid step type")
		}
	})

	t.Run("Steps with no parameters", func(t *testing.T) {
		req := models.StepDefinitionRequest{
			Steps: []models.StepDefinition{
				{Type: "Given", Pattern: "I am logged in", Parameters: ""},
				{Type: "When", Pattern: "I navigate to dashboard", Parameters: ""},
				{Type: "Then", Pattern: "I should see the dashboard", Parameters: ""},
			},
		}

		resp, err := generator.Generate(req)
		if err != nil {
			t.Fatalf("Expected successful generation, got error: %v", err)
		}

		if len(resp.Code) == 0 {
			t.Error("Expected generated code to be non-empty")
		}
	})

	t.Run("Mixed step types", func(t *testing.T) {
		req := models.StepDefinitionRequest{
			Steps: []models.StepDefinition{
				{Type: "Given", Pattern: "precondition", Parameters: ""},
				{Type: "Given", Pattern: "another precondition", Parameters: ""},
				{Type: "When", Pattern: "action", Parameters: ""},
				{Type: "Then", Pattern: "assertion", Parameters: ""},
				{Type: "Then", Pattern: "another assertion", Parameters: ""},
			},
		}

		resp, err := generator.Generate(req)
		if err != nil {
			t.Fatalf("Expected successful generation, got error: %v", err)
		}

		if resp.Metadata["stepCount"] != "5" {
			t.Errorf("Expected stepCount '5', got %q", resp.Metadata["stepCount"])
		}
	})
}

func TestStepGenerator_ValidateRequest(t *testing.T) {
	generator := setupStepGenerator(t)

	t.Run("Valid request", func(t *testing.T) {
		req := models.StepDefinitionRequest{
			Steps: []models.StepDefinition{
				{Type: "Given", Pattern: "test", Parameters: ""},
			},
		}

		err := generator.ValidateRequest(req)
		if err != nil {
			t.Errorf("Expected valid request, got error: %v", err)
		}
	})

	t.Run("Invalid request", func(t *testing.T) {
		req := models.StepDefinitionRequest{
			Steps: []models.StepDefinition{},
		}

		err := generator.ValidateRequest(req)
		if err == nil {
			t.Error("Expected validation error")
		}
	})
}
