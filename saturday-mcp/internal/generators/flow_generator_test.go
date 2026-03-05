package generators

import (
	"strings"
	"testing"
	"time"

	"github.com/orieken/saturday-mcp/internal/models"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/validators"
)

func setupFlowGenerator(t *testing.T) *FlowGenerator {
	t.Helper()

	registry := templates.NewRegistry()
	loader := templates.NewLoader(registry)
	cache := templates.NewCache(5 * time.Minute)
	processor := templates.NewProcessor(registry, cache)

	if err := loader.LoadAll(); err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	validator := validators.NewValidator()
	return NewFlowGenerator(processor, validator)
}

func TestFlowGenerator_Generate(t *testing.T) {
	generator := setupFlowGenerator(t)

	t.Run("Valid flow generation", func(t *testing.T) {
		req := models.FlowGenerationRequest{
			Name:        "checkout",
			Steps:       []string{"addToCart", "enterShipping", "enterPayment", "confirmOrder"},
			Description: "Checkout flow",
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
			"CheckoutFlow",
			"extends BaseFlow",
			"async execute()",
			"addToCart",
			"enterShipping",
			"enterPayment",
			"confirmOrder",
		}

		for _, expected := range expectedStrings {
			if !strings.Contains(resp.Code, expected) {
				t.Errorf("Expected code to contain %q", expected)
			}
		}

		// Verify filename
		if resp.FileName != "checkout-flow.ts" {
			t.Errorf("Expected fileName 'checkout-flow.ts', got %q", resp.FileName)
		}

		// Verify metadata
		if resp.Metadata["type"] != "flow" {
			t.Errorf("Expected type 'flow', got %q", resp.Metadata["type"])
		}

		if resp.Metadata["stepCount"] != "4" {
			t.Errorf("Expected stepCount '4', got %q", resp.Metadata["stepCount"])
		}
	})

	t.Run("Invalid request - missing name", func(t *testing.T) {
		req := models.FlowGenerationRequest{
			Steps: []string{"step1"},
		}

		_, err := generator.Generate(req)
		if err == nil {
			t.Error("Expected validation error for missing name")
		}
	})

	t.Run("Invalid request - empty steps", func(t *testing.T) {
		req := models.FlowGenerationRequest{
			Name:  "test",
			Steps: []string{},
		}

		_, err := generator.Generate(req)
		if err == nil {
			t.Error("Expected validation error for empty steps")
		}
	})

	t.Run("Flow with special characters in name", func(t *testing.T) {
		req := models.FlowGenerationRequest{
			Name:  "user_registration",
			Steps: []string{"fillForm", "submitForm"},
		}

		resp, err := generator.Generate(req)
		if err != nil {
			t.Fatalf("Expected successful generation, got error: %v", err)
		}

		if !strings.Contains(resp.Code, "UserRegistrationFlow") {
			t.Error("Expected name to be converted to PascalCase")
		}

		if resp.FileName != "user-registration-flow.ts" {
			t.Errorf("Expected kebab-case filename, got %q", resp.FileName)
		}
	})

	t.Run("Flow with single step", func(t *testing.T) {
		req := models.FlowGenerationRequest{
			Name:  "simple",
			Steps: []string{"doSomething"},
		}

		resp, err := generator.Generate(req)
		if err != nil {
			t.Fatalf("Expected successful generation, got error: %v", err)
		}

		if resp.Metadata["stepCount"] != "1" {
			t.Errorf("Expected stepCount '1', got %q", resp.Metadata["stepCount"])
		}
	})
}

func TestFlowGenerator_ValidateRequest(t *testing.T) {
	generator := setupFlowGenerator(t)

	t.Run("Valid request", func(t *testing.T) {
		req := models.FlowGenerationRequest{
			Name:  "test",
			Steps: []string{"step1"},
		}

		err := generator.ValidateRequest(req)
		if err != nil {
			t.Errorf("Expected valid request, got error: %v", err)
		}
	})

	t.Run("Invalid request", func(t *testing.T) {
		req := models.FlowGenerationRequest{
			Name: "test",
			// Missing Steps
		}

		err := generator.ValidateRequest(req)
		if err == nil {
			t.Error("Expected validation error")
		}
	})
}
