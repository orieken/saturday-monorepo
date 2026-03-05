package templates

import (
	"strings"
	"testing"
	"time"
)

// TestTemplateSystemIntegration tests the complete template system workflow
func TestTemplateSystemIntegration(t *testing.T) {
	// Create components
	registry := NewRegistry()
	loader := NewLoader(registry)
	cache := NewCache(5 * time.Minute)
	processor := NewProcessor(registry, cache)

	// Load all templates
	err := loader.LoadAll()
	if err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	t.Run("Generate Page from template", func(t *testing.T) {
		data := map[string]interface{}{
			"Name": "login",
			"Path": "/login",
			"Elements": []map[string]string{
				{"Name": "usernameInput", "Selector": "#username"},
				{"Name": "passwordInput", "Selector": "#password"},
				{"Name": "submitButton", "Selector": "button[type='submit']"},
			},
		}

		result, err := processor.Process("page", data)
		if err != nil {
			t.Fatalf("Failed to process page template: %v", err)
		}

		// Verify output contains expected content
		expectedStrings := []string{
			"LoginPage",
			"extends BasePage",
			"'/login'",
			"usernameInput",
			"passwordInput",
			"submitButton",
		}

		for _, expected := range expectedStrings {
			if !strings.Contains(result, expected) {
				t.Errorf("Expected output to contain %q", expected)
			}
		}
	})

	t.Run("Generate Site from template", func(t *testing.T) {
		data := map[string]interface{}{
			"Name":    "ecommerce",
			"Imports": []string{"home", "product", "cart"},
			"Pages":   []string{"home", "product", "cart"},
		}

		result, err := processor.Process("site", data)
		if err != nil {
			t.Fatalf("Failed to process site template: %v", err)
		}

		// Verify output contains expected content
		expectedStrings := []string{
			"EcommerceSite",
			"extends BaseSite",
			"HomePage",
			"ProductPage",
			"CartPage",
			"registerPage",
		}

		for _, expected := range expectedStrings {
			if !strings.Contains(result, expected) {
				t.Errorf("Expected output to contain %q", expected)
			}
		}
	})

	t.Run("Cache performance", func(t *testing.T) {
		data := map[string]interface{}{
			"Name": "test",
			"Path": "/test",
			"Elements": []map[string]string{
				{"Name": "element1", "Selector": "#el1"},
			},
		}

		// First call - cache miss
		_, err := processor.Process("page", data)
		if err != nil {
			t.Fatalf("Failed to process template: %v", err)
		}

		// Second call - should hit cache
		_, err = processor.Process("page", data)
		if err != nil {
			t.Fatalf("Failed to process template: %v", err)
		}

		if cache.Size() == 0 {
			t.Error("Expected cache to have entries")
		}
	})

	t.Run("Helper functions in templates", func(t *testing.T) {
		data := map[string]interface{}{
			"Name": "my_page",
			"Path": "/my-page",
			"Elements": []map[string]string{},
		}

		result, err := processor.Process("page", data)
		if err != nil {
			t.Fatalf("Failed to process template: %v", err)
		}

		// Should convert my_page to MyPage using pascalCase
		if !strings.Contains(result, "MyPagePage") {
			t.Error("Expected helper function pascalCase to be applied")
		}
	})
}
